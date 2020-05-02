package vpc

import (
	"errors"

	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ kore.ServiceProvider = Provider{}

type Provider struct {
	name string
}

func (p Provider) Name() string {
	return p.name
}

func (p Provider) Kinds() []string {
	return []string{"aws-vpc"}
}

func (p Provider) Plans() []servicesv1.ServicePlan {
	return []servicesv1.ServicePlan{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServicePlan",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "eks-dev",
				Namespace: "kore",
			},
			Spec: servicesv1.ServicePlanSpec{
				Kind:        "aws-vpc",
				Summary:     "AWS VPC for a development EKS cluster",
				Description: "Two private subnets with NAT gateways and two public subnets with direct Internet connection",
				Configuration: v1beta1.JSON{Raw: []byte(`{
					"region": "eu-west-2",
					"privateIPV4Cidr": "10.0.0.0/16",
					"subnetCount": 2
				}`)},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServicePlan",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "eks-prod",
				Namespace: "kore",
			},
			Spec: servicesv1.ServicePlanSpec{
				Kind:        "aws-vpc",
				Summary:     "AWS VPC for a production EKS cluster",
				Description: "Three private subnets with NAT gateways and three public subnets with direct Internet connection",
				Configuration: v1beta1.JSON{Raw: []byte(`{
					"region": "eu-west-2",
					"privateIPV4Cidr": "10.0.0.0/16",
					"subnetCount": 3
				}`)},
			},
		},
	}
}

func (p Provider) PlanJSONSchema(kind string, plan string) (string, error) {
	return `{
		"$id": "https://appvia.io/schemas/services/dummy/dummy.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "Dummy service plan schema",
		"type": "object",
		"additionalProperties": false,
		"required": [
			"region",
			"privateIPV4Cidr",
			"subnetCount"
		],
		"properties": {
			"name": {
				"type": "string",
				"minLength": 1
			},
			"region": {
				"type": "string",
				"minLength": 1
			},
			"privateIPV4Cidr": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"subnetCount": {
				"type": "integer"
			}
		}
	}`, nil
}

func (p Provider) CredentialsJSONSchema(kind string, plan string) (string, error) {
	return "", nil
}

func (p Provider) RequiredCredentialTypes(kind string) ([]schema.GroupVersionKind, error) {
	return []schema.GroupVersionKind{
		eksv1alpha1.EKSCredentialsGVK,
	}, nil
}

func (p Provider) ReconcileCredentials(
	kore.ServiceProviderContext,
	*servicesv1.Service,
	*servicesv1.ServiceCredentials) (reconcile.Result, map[string]string, error) {
	return reconcile.Result{}, nil, errors.New("service credentials are not supported for this service")
}

func (p Provider) DeleteCredentials(
	kore.ServiceProviderContext,
	*servicesv1.Service,
	*servicesv1.ServiceCredentials) (reconcile.Result, error) {
	return reconcile.Result{}, errors.New("service credentials are not supported for this service")
}
