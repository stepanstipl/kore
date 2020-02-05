package namespaceclaim

import (
	"context"
	"fmt"

	clusterv1 "github.com/appvia/hub-apis/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/hub-apis/pkg/apis/config/v1"
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"
	orgv1 "github.com/appvia/hub-apis/pkg/apis/org/v1"
	kubev1 "github.com/appvia/kube-operator/pkg/apis/kube/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_namespaceclaim")

// Add creates a new NamespaceClaim Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNamespaceClaim{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("namespaceclaim-controller", mgr, controller.Options{
		MaxConcurrentReconciles: 10,
		Reconciler:              r,
	})
	if err != nil {
		return err
	}
	// @step: we need to start watching teams in order to provision the rbac on the other
	if err := configv1.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("failed to add schema: %s into manager: %s", configv1.GroupName, err)
	}
	if err := orgv1.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("failed to add schema: %s into manager: %s", orgv1.GroupName, err)
	}
	if err := clusterv1.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("failed to add schema: %s into manager: %s", clusterv1.GroupName, err)
	}

	// @step: we need to start watching the kinds - we only care about changes
	// to the spec not the /status
	log.WithValues(
		"apiGroup", kubev1.SchemeGroupVersion.String(),
	).Info("adding the watch api group")

	err = c.Watch(&source.Kind{Type: &kubev1.NamespaceClaim{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	// @clause: whenever a team membership is changed we need queue all namespaceclaims
	// within the team namespace to be reconciled as well
	log.WithValues(
		"apiGroup", orgv1.SchemeGroupVersion.String(),
	).Info("adding the watch api group")

	ctx := context.Background()

	err = c.Watch(&source.Kind{Type: &orgv1.TeamMembership{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(o handler.MapObject) []reconcile.Request {
			requests, err := ReconcileNamespaceClaims(ctx, mgr.GetClient(), o.Meta.GetName(), o.Meta.GetNamespace())
			if err != nil {
				log.Error(err, "failed to force reconcilation of namespaceclaims from trigger")

				return []reconcile.Request{}
			}

			return requests
		}),
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNamespaceClaim implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNamespaceClaim{}

// ReconcileNamespaceClaim reconciles a NamespaceClaim object
type ReconcileNamespaceClaim struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NamespaceClaim object and makes changes based on the state read
// and what is in the NamespaceClaim.Spec
func (r *ReconcileNamespaceClaim) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues(
		"resource.namespace", request.Namespace,
		"resource.name", request.Name)

	reqLogger.Info("attempting to reconcile resource")

	ctx := context.Background()

	// Fetch the NamespaceClaim instance
	resource := &kubev1.NamespaceClaim{}
	err := r.client.Get(context.TODO(), request.NamespacedName, resource)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// Logic:
	// a) check if there is anyone else looking to finalize on the resource
	// b) create a kubernetes client from the cluster reference
	// c) check if the namespace claim is being delete and if so remove the rbac rules
	// d) if not deleting verify the rolebinding is correct in the remote cluster

	// @step: lets be pesimesstic by nature
	resource.Status.Status = metav1.StatusFailure

	err = func() error {
		teamName := HubLabel(resource, "team")

		// @step: check the team the namespace is associated to exists - adding a guard clause
		_, found, err := IsTeam(ctx, r.client, teamName)
		if err != nil {
			resource.Status.Conditions = []core.Condition{{
				Detail:  err.Error(),
				Message: fmt.Sprintf("failed to check for existence of team: %s", teamName),
			}}

			return err
		} else if !found {
			resource.Status.Conditions = []core.Condition{{
				Detail:  err.Error(),
				Message: fmt.Sprintf("team: %s does not exist", teamName),
			}}

			// meaning to don't bother to retry the reconcilation as there's nothing
			// we can do to resolve this
			return nil
		}

		// @step: retrieve the credentials from the associated cluster
		kc, err := MakeClusterKubeClient(ctx, r.client, types.NamespacedName{
			Namespace: resource.Spec.Use.Namespace,
			Name:      resource.Spec.Use.Name,
		})
		if err != nil {
			resource.Status.Conditions = []core.Condition{
				{
					Detail: err.Error(),
					Message: fmt.Sprintf("failed creating kubernetes client from (%s/%s)",
						resource.Spec.Use.Namespace, resource.Spec.Use.Name),
				},
			}

			reqLogger.WithValues(
				"cluster.name", resource.Spec.Use.Name,
				"cluster.namespace", resource.Spec.Use.Namespace,
				"team", teamName,
			).Error(err, "failed to create a kubernetes client from cluster configuration")

			return err
		}

		// @step: check if the resource is being deleted
		if resource.DeletionTimestamp != nil {
			return r.Delete(ctx, r.client, kc, resource)
		}

		return r.Update(ctx, r.client, kc, resource)
	}()
	if err != nil {
		reqLogger.Error(err, "failed to reconcile the namespace claim")

		if err := r.client.Status().Update(ctx, resource); err != nil {
			reqLogger.Error(err, "failed to update the status of resource")

			return reconcile.Result{}, err
		}
	}

	// @thank f$ck it was all ok in the end
	return reconcile.Result{}, nil
}
