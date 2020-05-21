// Code generated by counterfeiter. DO NOT EDIT.
package korefakes

import (
	"context"
	"sync"

	v1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
)

type FakeServicePlans struct {
	DeleteStub        func(context.Context, string) (*v1.ServicePlan, error)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	deleteReturns struct {
		result1 *v1.ServicePlan
		result2 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 *v1.ServicePlan
		result2 error
	}
	GetStub        func(context.Context, string) (*v1.ServicePlan, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	getReturns struct {
		result1 *v1.ServicePlan
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 *v1.ServicePlan
		result2 error
	}
	GetCredentialSchemaStub        func(context.Context, string) (string, error)
	getCredentialSchemaMutex       sync.RWMutex
	getCredentialSchemaArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	getCredentialSchemaReturns struct {
		result1 string
		result2 error
	}
	getCredentialSchemaReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetEditablePlanParamsStub        func(context.Context, string, string) (map[string]bool, error)
	getEditablePlanParamsMutex       sync.RWMutex
	getEditablePlanParamsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getEditablePlanParamsReturns struct {
		result1 map[string]bool
		result2 error
	}
	getEditablePlanParamsReturnsOnCall map[int]struct {
		result1 map[string]bool
		result2 error
	}
	GetSchemaStub        func(context.Context, string) (string, error)
	getSchemaMutex       sync.RWMutex
	getSchemaArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	getSchemaReturns struct {
		result1 string
		result2 error
	}
	getSchemaReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	HasStub        func(context.Context, string) (bool, error)
	hasMutex       sync.RWMutex
	hasArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	hasReturns struct {
		result1 bool
		result2 error
	}
	hasReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	ListStub        func(context.Context) (*v1.ServicePlanList, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 context.Context
	}
	listReturns struct {
		result1 *v1.ServicePlanList
		result2 error
	}
	listReturnsOnCall map[int]struct {
		result1 *v1.ServicePlanList
		result2 error
	}
	ListFilteredStub        func(context.Context, func(v1.ServicePlan) bool) (*v1.ServicePlanList, error)
	listFilteredMutex       sync.RWMutex
	listFilteredArgsForCall []struct {
		arg1 context.Context
		arg2 func(v1.ServicePlan) bool
	}
	listFilteredReturns struct {
		result1 *v1.ServicePlanList
		result2 error
	}
	listFilteredReturnsOnCall map[int]struct {
		result1 *v1.ServicePlanList
		result2 error
	}
	UpdateStub        func(context.Context, *v1.ServicePlan) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 context.Context
		arg2 *v1.ServicePlan
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeServicePlans) Delete(arg1 context.Context, arg2 string) (*v1.ServicePlan, error) {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Delete", []interface{}{arg1, arg2})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeServicePlans) DeleteCalls(stub func(context.Context, string) (*v1.ServicePlan, error)) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeServicePlans) DeleteArgsForCall(i int) (context.Context, string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) DeleteReturns(result1 *v1.ServicePlan, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 *v1.ServicePlan
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) DeleteReturnsOnCall(i int, result1 *v1.ServicePlan, result2 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 *v1.ServicePlan
			result2 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 *v1.ServicePlan
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) Get(arg1 context.Context, arg2 string) (*v1.ServicePlan, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Get", []interface{}{arg1, arg2})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeServicePlans) GetCalls(stub func(context.Context, string) (*v1.ServicePlan, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeServicePlans) GetArgsForCall(i int) (context.Context, string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) GetReturns(result1 *v1.ServicePlan, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *v1.ServicePlan
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetReturnsOnCall(i int, result1 *v1.ServicePlan, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *v1.ServicePlan
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *v1.ServicePlan
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetCredentialSchema(arg1 context.Context, arg2 string) (string, error) {
	fake.getCredentialSchemaMutex.Lock()
	ret, specificReturn := fake.getCredentialSchemaReturnsOnCall[len(fake.getCredentialSchemaArgsForCall)]
	fake.getCredentialSchemaArgsForCall = append(fake.getCredentialSchemaArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetCredentialSchema", []interface{}{arg1, arg2})
	fake.getCredentialSchemaMutex.Unlock()
	if fake.GetCredentialSchemaStub != nil {
		return fake.GetCredentialSchemaStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getCredentialSchemaReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) GetCredentialSchemaCallCount() int {
	fake.getCredentialSchemaMutex.RLock()
	defer fake.getCredentialSchemaMutex.RUnlock()
	return len(fake.getCredentialSchemaArgsForCall)
}

func (fake *FakeServicePlans) GetCredentialSchemaCalls(stub func(context.Context, string) (string, error)) {
	fake.getCredentialSchemaMutex.Lock()
	defer fake.getCredentialSchemaMutex.Unlock()
	fake.GetCredentialSchemaStub = stub
}

func (fake *FakeServicePlans) GetCredentialSchemaArgsForCall(i int) (context.Context, string) {
	fake.getCredentialSchemaMutex.RLock()
	defer fake.getCredentialSchemaMutex.RUnlock()
	argsForCall := fake.getCredentialSchemaArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) GetCredentialSchemaReturns(result1 string, result2 error) {
	fake.getCredentialSchemaMutex.Lock()
	defer fake.getCredentialSchemaMutex.Unlock()
	fake.GetCredentialSchemaStub = nil
	fake.getCredentialSchemaReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetCredentialSchemaReturnsOnCall(i int, result1 string, result2 error) {
	fake.getCredentialSchemaMutex.Lock()
	defer fake.getCredentialSchemaMutex.Unlock()
	fake.GetCredentialSchemaStub = nil
	if fake.getCredentialSchemaReturnsOnCall == nil {
		fake.getCredentialSchemaReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getCredentialSchemaReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetEditablePlanParams(arg1 context.Context, arg2 string, arg3 string) (map[string]bool, error) {
	fake.getEditablePlanParamsMutex.Lock()
	ret, specificReturn := fake.getEditablePlanParamsReturnsOnCall[len(fake.getEditablePlanParamsArgsForCall)]
	fake.getEditablePlanParamsArgsForCall = append(fake.getEditablePlanParamsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetEditablePlanParams", []interface{}{arg1, arg2, arg3})
	fake.getEditablePlanParamsMutex.Unlock()
	if fake.GetEditablePlanParamsStub != nil {
		return fake.GetEditablePlanParamsStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getEditablePlanParamsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) GetEditablePlanParamsCallCount() int {
	fake.getEditablePlanParamsMutex.RLock()
	defer fake.getEditablePlanParamsMutex.RUnlock()
	return len(fake.getEditablePlanParamsArgsForCall)
}

func (fake *FakeServicePlans) GetEditablePlanParamsCalls(stub func(context.Context, string, string) (map[string]bool, error)) {
	fake.getEditablePlanParamsMutex.Lock()
	defer fake.getEditablePlanParamsMutex.Unlock()
	fake.GetEditablePlanParamsStub = stub
}

func (fake *FakeServicePlans) GetEditablePlanParamsArgsForCall(i int) (context.Context, string, string) {
	fake.getEditablePlanParamsMutex.RLock()
	defer fake.getEditablePlanParamsMutex.RUnlock()
	argsForCall := fake.getEditablePlanParamsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeServicePlans) GetEditablePlanParamsReturns(result1 map[string]bool, result2 error) {
	fake.getEditablePlanParamsMutex.Lock()
	defer fake.getEditablePlanParamsMutex.Unlock()
	fake.GetEditablePlanParamsStub = nil
	fake.getEditablePlanParamsReturns = struct {
		result1 map[string]bool
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetEditablePlanParamsReturnsOnCall(i int, result1 map[string]bool, result2 error) {
	fake.getEditablePlanParamsMutex.Lock()
	defer fake.getEditablePlanParamsMutex.Unlock()
	fake.GetEditablePlanParamsStub = nil
	if fake.getEditablePlanParamsReturnsOnCall == nil {
		fake.getEditablePlanParamsReturnsOnCall = make(map[int]struct {
			result1 map[string]bool
			result2 error
		})
	}
	fake.getEditablePlanParamsReturnsOnCall[i] = struct {
		result1 map[string]bool
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetSchema(arg1 context.Context, arg2 string) (string, error) {
	fake.getSchemaMutex.Lock()
	ret, specificReturn := fake.getSchemaReturnsOnCall[len(fake.getSchemaArgsForCall)]
	fake.getSchemaArgsForCall = append(fake.getSchemaArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetSchema", []interface{}{arg1, arg2})
	fake.getSchemaMutex.Unlock()
	if fake.GetSchemaStub != nil {
		return fake.GetSchemaStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getSchemaReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) GetSchemaCallCount() int {
	fake.getSchemaMutex.RLock()
	defer fake.getSchemaMutex.RUnlock()
	return len(fake.getSchemaArgsForCall)
}

func (fake *FakeServicePlans) GetSchemaCalls(stub func(context.Context, string) (string, error)) {
	fake.getSchemaMutex.Lock()
	defer fake.getSchemaMutex.Unlock()
	fake.GetSchemaStub = stub
}

func (fake *FakeServicePlans) GetSchemaArgsForCall(i int) (context.Context, string) {
	fake.getSchemaMutex.RLock()
	defer fake.getSchemaMutex.RUnlock()
	argsForCall := fake.getSchemaArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) GetSchemaReturns(result1 string, result2 error) {
	fake.getSchemaMutex.Lock()
	defer fake.getSchemaMutex.Unlock()
	fake.GetSchemaStub = nil
	fake.getSchemaReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) GetSchemaReturnsOnCall(i int, result1 string, result2 error) {
	fake.getSchemaMutex.Lock()
	defer fake.getSchemaMutex.Unlock()
	fake.GetSchemaStub = nil
	if fake.getSchemaReturnsOnCall == nil {
		fake.getSchemaReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getSchemaReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) Has(arg1 context.Context, arg2 string) (bool, error) {
	fake.hasMutex.Lock()
	ret, specificReturn := fake.hasReturnsOnCall[len(fake.hasArgsForCall)]
	fake.hasArgsForCall = append(fake.hasArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Has", []interface{}{arg1, arg2})
	fake.hasMutex.Unlock()
	if fake.HasStub != nil {
		return fake.HasStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.hasReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) HasCallCount() int {
	fake.hasMutex.RLock()
	defer fake.hasMutex.RUnlock()
	return len(fake.hasArgsForCall)
}

func (fake *FakeServicePlans) HasCalls(stub func(context.Context, string) (bool, error)) {
	fake.hasMutex.Lock()
	defer fake.hasMutex.Unlock()
	fake.HasStub = stub
}

func (fake *FakeServicePlans) HasArgsForCall(i int) (context.Context, string) {
	fake.hasMutex.RLock()
	defer fake.hasMutex.RUnlock()
	argsForCall := fake.hasArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) HasReturns(result1 bool, result2 error) {
	fake.hasMutex.Lock()
	defer fake.hasMutex.Unlock()
	fake.HasStub = nil
	fake.hasReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) HasReturnsOnCall(i int, result1 bool, result2 error) {
	fake.hasMutex.Lock()
	defer fake.hasMutex.Unlock()
	fake.HasStub = nil
	if fake.hasReturnsOnCall == nil {
		fake.hasReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.hasReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) List(arg1 context.Context) (*v1.ServicePlanList, error) {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("List", []interface{}{arg1})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeServicePlans) ListCalls(stub func(context.Context) (*v1.ServicePlanList, error)) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *FakeServicePlans) ListArgsForCall(i int) context.Context {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeServicePlans) ListReturns(result1 *v1.ServicePlanList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 *v1.ServicePlanList
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) ListReturnsOnCall(i int, result1 *v1.ServicePlanList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 *v1.ServicePlanList
			result2 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 *v1.ServicePlanList
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) ListFiltered(arg1 context.Context, arg2 func(v1.ServicePlan) bool) (*v1.ServicePlanList, error) {
	fake.listFilteredMutex.Lock()
	ret, specificReturn := fake.listFilteredReturnsOnCall[len(fake.listFilteredArgsForCall)]
	fake.listFilteredArgsForCall = append(fake.listFilteredArgsForCall, struct {
		arg1 context.Context
		arg2 func(v1.ServicePlan) bool
	}{arg1, arg2})
	fake.recordInvocation("ListFiltered", []interface{}{arg1, arg2})
	fake.listFilteredMutex.Unlock()
	if fake.ListFilteredStub != nil {
		return fake.ListFilteredStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listFilteredReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServicePlans) ListFilteredCallCount() int {
	fake.listFilteredMutex.RLock()
	defer fake.listFilteredMutex.RUnlock()
	return len(fake.listFilteredArgsForCall)
}

func (fake *FakeServicePlans) ListFilteredCalls(stub func(context.Context, func(v1.ServicePlan) bool) (*v1.ServicePlanList, error)) {
	fake.listFilteredMutex.Lock()
	defer fake.listFilteredMutex.Unlock()
	fake.ListFilteredStub = stub
}

func (fake *FakeServicePlans) ListFilteredArgsForCall(i int) (context.Context, func(v1.ServicePlan) bool) {
	fake.listFilteredMutex.RLock()
	defer fake.listFilteredMutex.RUnlock()
	argsForCall := fake.listFilteredArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) ListFilteredReturns(result1 *v1.ServicePlanList, result2 error) {
	fake.listFilteredMutex.Lock()
	defer fake.listFilteredMutex.Unlock()
	fake.ListFilteredStub = nil
	fake.listFilteredReturns = struct {
		result1 *v1.ServicePlanList
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) ListFilteredReturnsOnCall(i int, result1 *v1.ServicePlanList, result2 error) {
	fake.listFilteredMutex.Lock()
	defer fake.listFilteredMutex.Unlock()
	fake.ListFilteredStub = nil
	if fake.listFilteredReturnsOnCall == nil {
		fake.listFilteredReturnsOnCall = make(map[int]struct {
			result1 *v1.ServicePlanList
			result2 error
		})
	}
	fake.listFilteredReturnsOnCall[i] = struct {
		result1 *v1.ServicePlanList
		result2 error
	}{result1, result2}
}

func (fake *FakeServicePlans) Update(arg1 context.Context, arg2 *v1.ServicePlan) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 context.Context
		arg2 *v1.ServicePlan
	}{arg1, arg2})
	fake.recordInvocation("Update", []interface{}{arg1, arg2})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.updateReturns
	return fakeReturns.result1
}

func (fake *FakeServicePlans) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeServicePlans) UpdateCalls(stub func(context.Context, *v1.ServicePlan) error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *FakeServicePlans) UpdateArgsForCall(i int) (context.Context, *v1.ServicePlan) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServicePlans) UpdateReturns(result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeServicePlans) UpdateReturnsOnCall(i int, result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeServicePlans) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getCredentialSchemaMutex.RLock()
	defer fake.getCredentialSchemaMutex.RUnlock()
	fake.getEditablePlanParamsMutex.RLock()
	defer fake.getEditablePlanParamsMutex.RUnlock()
	fake.getSchemaMutex.RLock()
	defer fake.getSchemaMutex.RUnlock()
	fake.hasMutex.RLock()
	defer fake.hasMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	fake.listFilteredMutex.RLock()
	defer fake.listFilteredMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeServicePlans) recordInvocation(key string, args []interface{}) {
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

var _ kore.ServicePlans = new(FakeServicePlans)
