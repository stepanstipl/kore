// Code generated by counterfeiter. DO NOT EDIT.
package korefakes

import (
	"sync"

	v1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	v1a "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
)

type FakeClusterProvider struct {
	BeforeComponentsUpdateStub        func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error
	beforeComponentsUpdateMutex       sync.RWMutex
	beforeComponentsUpdateArgsForCall []struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}
	beforeComponentsUpdateReturns struct {
		result1 error
	}
	beforeComponentsUpdateReturnsOnCall map[int]struct {
		result1 error
	}
	DefaultPlanPolicyStub        func() *v1a.PlanPolicy
	defaultPlanPolicyMutex       sync.RWMutex
	defaultPlanPolicyArgsForCall []struct {
	}
	defaultPlanPolicyReturns struct {
		result1 *v1a.PlanPolicy
	}
	defaultPlanPolicyReturnsOnCall map[int]struct {
		result1 *v1a.PlanPolicy
	}
	DefaultPlansStub        func() []v1a.Plan
	defaultPlansMutex       sync.RWMutex
	defaultPlansArgsForCall []struct {
	}
	defaultPlansReturns struct {
		result1 []v1a.Plan
	}
	defaultPlansReturnsOnCall map[int]struct {
		result1 []v1a.Plan
	}
	PlanJSONSchemaStub        func() string
	planJSONSchemaMutex       sync.RWMutex
	planJSONSchemaArgsForCall []struct {
	}
	planJSONSchemaReturns struct {
		result1 string
	}
	planJSONSchemaReturnsOnCall map[int]struct {
		result1 string
	}
	SetComponentsStub        func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error
	setComponentsMutex       sync.RWMutex
	setComponentsArgsForCall []struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}
	setComponentsReturns struct {
		result1 error
	}
	setComponentsReturnsOnCall map[int]struct {
		result1 error
	}
	SetProviderDataStub        func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error
	setProviderDataMutex       sync.RWMutex
	setProviderDataArgsForCall []struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}
	setProviderDataReturns struct {
		result1 error
	}
	setProviderDataReturnsOnCall map[int]struct {
		result1 error
	}
	TypeStub        func() string
	typeMutex       sync.RWMutex
	typeArgsForCall []struct {
	}
	typeReturns struct {
		result1 string
	}
	typeReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClusterProvider) BeforeComponentsUpdate(arg1 kore.Context, arg2 *v1.Cluster, arg3 *kore.ClusterComponents) error {
	fake.beforeComponentsUpdateMutex.Lock()
	ret, specificReturn := fake.beforeComponentsUpdateReturnsOnCall[len(fake.beforeComponentsUpdateArgsForCall)]
	fake.beforeComponentsUpdateArgsForCall = append(fake.beforeComponentsUpdateArgsForCall, struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}{arg1, arg2, arg3})
	fake.recordInvocation("BeforeComponentsUpdate", []interface{}{arg1, arg2, arg3})
	fake.beforeComponentsUpdateMutex.Unlock()
	if fake.BeforeComponentsUpdateStub != nil {
		return fake.BeforeComponentsUpdateStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.beforeComponentsUpdateReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) BeforeComponentsUpdateCallCount() int {
	fake.beforeComponentsUpdateMutex.RLock()
	defer fake.beforeComponentsUpdateMutex.RUnlock()
	return len(fake.beforeComponentsUpdateArgsForCall)
}

func (fake *FakeClusterProvider) BeforeComponentsUpdateCalls(stub func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error) {
	fake.beforeComponentsUpdateMutex.Lock()
	defer fake.beforeComponentsUpdateMutex.Unlock()
	fake.BeforeComponentsUpdateStub = stub
}

func (fake *FakeClusterProvider) BeforeComponentsUpdateArgsForCall(i int) (kore.Context, *v1.Cluster, *kore.ClusterComponents) {
	fake.beforeComponentsUpdateMutex.RLock()
	defer fake.beforeComponentsUpdateMutex.RUnlock()
	argsForCall := fake.beforeComponentsUpdateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClusterProvider) BeforeComponentsUpdateReturns(result1 error) {
	fake.beforeComponentsUpdateMutex.Lock()
	defer fake.beforeComponentsUpdateMutex.Unlock()
	fake.BeforeComponentsUpdateStub = nil
	fake.beforeComponentsUpdateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) BeforeComponentsUpdateReturnsOnCall(i int, result1 error) {
	fake.beforeComponentsUpdateMutex.Lock()
	defer fake.beforeComponentsUpdateMutex.Unlock()
	fake.BeforeComponentsUpdateStub = nil
	if fake.beforeComponentsUpdateReturnsOnCall == nil {
		fake.beforeComponentsUpdateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.beforeComponentsUpdateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) DefaultPlanPolicy() *v1a.PlanPolicy {
	fake.defaultPlanPolicyMutex.Lock()
	ret, specificReturn := fake.defaultPlanPolicyReturnsOnCall[len(fake.defaultPlanPolicyArgsForCall)]
	fake.defaultPlanPolicyArgsForCall = append(fake.defaultPlanPolicyArgsForCall, struct {
	}{})
	fake.recordInvocation("DefaultPlanPolicy", []interface{}{})
	fake.defaultPlanPolicyMutex.Unlock()
	if fake.DefaultPlanPolicyStub != nil {
		return fake.DefaultPlanPolicyStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.defaultPlanPolicyReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) DefaultPlanPolicyCallCount() int {
	fake.defaultPlanPolicyMutex.RLock()
	defer fake.defaultPlanPolicyMutex.RUnlock()
	return len(fake.defaultPlanPolicyArgsForCall)
}

func (fake *FakeClusterProvider) DefaultPlanPolicyCalls(stub func() *v1a.PlanPolicy) {
	fake.defaultPlanPolicyMutex.Lock()
	defer fake.defaultPlanPolicyMutex.Unlock()
	fake.DefaultPlanPolicyStub = stub
}

func (fake *FakeClusterProvider) DefaultPlanPolicyReturns(result1 *v1a.PlanPolicy) {
	fake.defaultPlanPolicyMutex.Lock()
	defer fake.defaultPlanPolicyMutex.Unlock()
	fake.DefaultPlanPolicyStub = nil
	fake.defaultPlanPolicyReturns = struct {
		result1 *v1a.PlanPolicy
	}{result1}
}

func (fake *FakeClusterProvider) DefaultPlanPolicyReturnsOnCall(i int, result1 *v1a.PlanPolicy) {
	fake.defaultPlanPolicyMutex.Lock()
	defer fake.defaultPlanPolicyMutex.Unlock()
	fake.DefaultPlanPolicyStub = nil
	if fake.defaultPlanPolicyReturnsOnCall == nil {
		fake.defaultPlanPolicyReturnsOnCall = make(map[int]struct {
			result1 *v1a.PlanPolicy
		})
	}
	fake.defaultPlanPolicyReturnsOnCall[i] = struct {
		result1 *v1a.PlanPolicy
	}{result1}
}

func (fake *FakeClusterProvider) DefaultPlans() []v1a.Plan {
	fake.defaultPlansMutex.Lock()
	ret, specificReturn := fake.defaultPlansReturnsOnCall[len(fake.defaultPlansArgsForCall)]
	fake.defaultPlansArgsForCall = append(fake.defaultPlansArgsForCall, struct {
	}{})
	fake.recordInvocation("DefaultPlans", []interface{}{})
	fake.defaultPlansMutex.Unlock()
	if fake.DefaultPlansStub != nil {
		return fake.DefaultPlansStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.defaultPlansReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) DefaultPlansCallCount() int {
	fake.defaultPlansMutex.RLock()
	defer fake.defaultPlansMutex.RUnlock()
	return len(fake.defaultPlansArgsForCall)
}

func (fake *FakeClusterProvider) DefaultPlansCalls(stub func() []v1a.Plan) {
	fake.defaultPlansMutex.Lock()
	defer fake.defaultPlansMutex.Unlock()
	fake.DefaultPlansStub = stub
}

func (fake *FakeClusterProvider) DefaultPlansReturns(result1 []v1a.Plan) {
	fake.defaultPlansMutex.Lock()
	defer fake.defaultPlansMutex.Unlock()
	fake.DefaultPlansStub = nil
	fake.defaultPlansReturns = struct {
		result1 []v1a.Plan
	}{result1}
}

func (fake *FakeClusterProvider) DefaultPlansReturnsOnCall(i int, result1 []v1a.Plan) {
	fake.defaultPlansMutex.Lock()
	defer fake.defaultPlansMutex.Unlock()
	fake.DefaultPlansStub = nil
	if fake.defaultPlansReturnsOnCall == nil {
		fake.defaultPlansReturnsOnCall = make(map[int]struct {
			result1 []v1a.Plan
		})
	}
	fake.defaultPlansReturnsOnCall[i] = struct {
		result1 []v1a.Plan
	}{result1}
}

func (fake *FakeClusterProvider) PlanJSONSchema() string {
	fake.planJSONSchemaMutex.Lock()
	ret, specificReturn := fake.planJSONSchemaReturnsOnCall[len(fake.planJSONSchemaArgsForCall)]
	fake.planJSONSchemaArgsForCall = append(fake.planJSONSchemaArgsForCall, struct {
	}{})
	fake.recordInvocation("PlanJSONSchema", []interface{}{})
	fake.planJSONSchemaMutex.Unlock()
	if fake.PlanJSONSchemaStub != nil {
		return fake.PlanJSONSchemaStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.planJSONSchemaReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) PlanJSONSchemaCallCount() int {
	fake.planJSONSchemaMutex.RLock()
	defer fake.planJSONSchemaMutex.RUnlock()
	return len(fake.planJSONSchemaArgsForCall)
}

func (fake *FakeClusterProvider) PlanJSONSchemaCalls(stub func() string) {
	fake.planJSONSchemaMutex.Lock()
	defer fake.planJSONSchemaMutex.Unlock()
	fake.PlanJSONSchemaStub = stub
}

func (fake *FakeClusterProvider) PlanJSONSchemaReturns(result1 string) {
	fake.planJSONSchemaMutex.Lock()
	defer fake.planJSONSchemaMutex.Unlock()
	fake.PlanJSONSchemaStub = nil
	fake.planJSONSchemaReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeClusterProvider) PlanJSONSchemaReturnsOnCall(i int, result1 string) {
	fake.planJSONSchemaMutex.Lock()
	defer fake.planJSONSchemaMutex.Unlock()
	fake.PlanJSONSchemaStub = nil
	if fake.planJSONSchemaReturnsOnCall == nil {
		fake.planJSONSchemaReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.planJSONSchemaReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeClusterProvider) SetComponents(arg1 kore.Context, arg2 *v1.Cluster, arg3 *kore.ClusterComponents) error {
	fake.setComponentsMutex.Lock()
	ret, specificReturn := fake.setComponentsReturnsOnCall[len(fake.setComponentsArgsForCall)]
	fake.setComponentsArgsForCall = append(fake.setComponentsArgsForCall, struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}{arg1, arg2, arg3})
	fake.recordInvocation("SetComponents", []interface{}{arg1, arg2, arg3})
	fake.setComponentsMutex.Unlock()
	if fake.SetComponentsStub != nil {
		return fake.SetComponentsStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setComponentsReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) SetComponentsCallCount() int {
	fake.setComponentsMutex.RLock()
	defer fake.setComponentsMutex.RUnlock()
	return len(fake.setComponentsArgsForCall)
}

func (fake *FakeClusterProvider) SetComponentsCalls(stub func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error) {
	fake.setComponentsMutex.Lock()
	defer fake.setComponentsMutex.Unlock()
	fake.SetComponentsStub = stub
}

func (fake *FakeClusterProvider) SetComponentsArgsForCall(i int) (kore.Context, *v1.Cluster, *kore.ClusterComponents) {
	fake.setComponentsMutex.RLock()
	defer fake.setComponentsMutex.RUnlock()
	argsForCall := fake.setComponentsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClusterProvider) SetComponentsReturns(result1 error) {
	fake.setComponentsMutex.Lock()
	defer fake.setComponentsMutex.Unlock()
	fake.SetComponentsStub = nil
	fake.setComponentsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) SetComponentsReturnsOnCall(i int, result1 error) {
	fake.setComponentsMutex.Lock()
	defer fake.setComponentsMutex.Unlock()
	fake.SetComponentsStub = nil
	if fake.setComponentsReturnsOnCall == nil {
		fake.setComponentsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setComponentsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) SetProviderData(arg1 kore.Context, arg2 *v1.Cluster, arg3 *kore.ClusterComponents) error {
	fake.setProviderDataMutex.Lock()
	ret, specificReturn := fake.setProviderDataReturnsOnCall[len(fake.setProviderDataArgsForCall)]
	fake.setProviderDataArgsForCall = append(fake.setProviderDataArgsForCall, struct {
		arg1 kore.Context
		arg2 *v1.Cluster
		arg3 *kore.ClusterComponents
	}{arg1, arg2, arg3})
	fake.recordInvocation("SetProviderData", []interface{}{arg1, arg2, arg3})
	fake.setProviderDataMutex.Unlock()
	if fake.SetProviderDataStub != nil {
		return fake.SetProviderDataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setProviderDataReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) SetProviderDataCallCount() int {
	fake.setProviderDataMutex.RLock()
	defer fake.setProviderDataMutex.RUnlock()
	return len(fake.setProviderDataArgsForCall)
}

func (fake *FakeClusterProvider) SetProviderDataCalls(stub func(kore.Context, *v1.Cluster, *kore.ClusterComponents) error) {
	fake.setProviderDataMutex.Lock()
	defer fake.setProviderDataMutex.Unlock()
	fake.SetProviderDataStub = stub
}

func (fake *FakeClusterProvider) SetProviderDataArgsForCall(i int) (kore.Context, *v1.Cluster, *kore.ClusterComponents) {
	fake.setProviderDataMutex.RLock()
	defer fake.setProviderDataMutex.RUnlock()
	argsForCall := fake.setProviderDataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClusterProvider) SetProviderDataReturns(result1 error) {
	fake.setProviderDataMutex.Lock()
	defer fake.setProviderDataMutex.Unlock()
	fake.SetProviderDataStub = nil
	fake.setProviderDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) SetProviderDataReturnsOnCall(i int, result1 error) {
	fake.setProviderDataMutex.Lock()
	defer fake.setProviderDataMutex.Unlock()
	fake.SetProviderDataStub = nil
	if fake.setProviderDataReturnsOnCall == nil {
		fake.setProviderDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setProviderDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClusterProvider) Type() string {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct {
	}{})
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.typeReturns
	return fakeReturns.result1
}

func (fake *FakeClusterProvider) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeClusterProvider) TypeCalls(stub func() string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = stub
}

func (fake *FakeClusterProvider) TypeReturns(result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeClusterProvider) TypeReturnsOnCall(i int, result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeClusterProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.beforeComponentsUpdateMutex.RLock()
	defer fake.beforeComponentsUpdateMutex.RUnlock()
	fake.defaultPlanPolicyMutex.RLock()
	defer fake.defaultPlanPolicyMutex.RUnlock()
	fake.defaultPlansMutex.RLock()
	defer fake.defaultPlansMutex.RUnlock()
	fake.planJSONSchemaMutex.RLock()
	defer fake.planJSONSchemaMutex.RUnlock()
	fake.setComponentsMutex.RLock()
	defer fake.setComponentsMutex.RUnlock()
	fake.setProviderDataMutex.RLock()
	defer fake.setProviderDataMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClusterProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ kore.ClusterProvider = new(FakeClusterProvider)