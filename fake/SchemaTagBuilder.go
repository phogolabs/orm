// This file was generated by counterfeiter
package fake

import (
	"sync"

	"github.com/phogolabs/oak/schema"
)

type SchemaTagBuilder struct {
	BuildStub        func(column *schema.Column) string
	buildMutex       sync.RWMutex
	buildArgsForCall []struct {
		column *schema.Column
	}
	buildReturns struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SchemaTagBuilder) Build(column *schema.Column) string {
	fake.buildMutex.Lock()
	fake.buildArgsForCall = append(fake.buildArgsForCall, struct {
		column *schema.Column
	}{column})
	fake.recordInvocation("Build", []interface{}{column})
	fake.buildMutex.Unlock()
	if fake.BuildStub != nil {
		return fake.BuildStub(column)
	}
	return fake.buildReturns.result1
}

func (fake *SchemaTagBuilder) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *SchemaTagBuilder) BuildArgsForCall(i int) *schema.Column {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return fake.buildArgsForCall[i].column
}

func (fake *SchemaTagBuilder) BuildReturns(result1 string) {
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 string
	}{result1}
}

func (fake *SchemaTagBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return fake.invocations
}

func (fake *SchemaTagBuilder) recordInvocation(key string, args []interface{}) {
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

var _ schema.TagBuilder = new(SchemaTagBuilder)