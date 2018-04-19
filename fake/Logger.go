// This file was generated by counterfeiter
package fake

import (
	"sync"

	"github.com/apex/log"
)

type Logger struct {
	WithFieldsStub        func(fields log.Fielder) *log.Entry
	withFieldsMutex       sync.RWMutex
	withFieldsArgsForCall []struct {
		fields log.Fielder
	}
	withFieldsReturns struct {
		result1 *log.Entry
	}
	WithFieldStub        func(key string, value interface{}) *log.Entry
	withFieldMutex       sync.RWMutex
	withFieldArgsForCall []struct {
		key   string
		value interface{}
	}
	withFieldReturns struct {
		result1 *log.Entry
	}
	WithErrorStub        func(err error) *log.Entry
	withErrorMutex       sync.RWMutex
	withErrorArgsForCall []struct {
		err error
	}
	withErrorReturns struct {
		result1 *log.Entry
	}
	DebugStub        func(msg string)
	debugMutex       sync.RWMutex
	debugArgsForCall []struct {
		msg string
	}
	InfoStub        func(msg string)
	infoMutex       sync.RWMutex
	infoArgsForCall []struct {
		msg string
	}
	WarnStub        func(msg string)
	warnMutex       sync.RWMutex
	warnArgsForCall []struct {
		msg string
	}
	ErrorStub        func(msg string)
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
		msg string
	}
	FatalStub        func(msg string)
	fatalMutex       sync.RWMutex
	fatalArgsForCall []struct {
		msg string
	}
	DebugfStub        func(msg string, v ...interface{})
	debugfMutex       sync.RWMutex
	debugfArgsForCall []struct {
		msg string
		v   []interface{}
	}
	InfofStub        func(msg string, v ...interface{})
	infofMutex       sync.RWMutex
	infofArgsForCall []struct {
		msg string
		v   []interface{}
	}
	WarnfStub        func(msg string, v ...interface{})
	warnfMutex       sync.RWMutex
	warnfArgsForCall []struct {
		msg string
		v   []interface{}
	}
	ErrorfStub        func(msg string, v ...interface{})
	errorfMutex       sync.RWMutex
	errorfArgsForCall []struct {
		msg string
		v   []interface{}
	}
	FatalfStub        func(msg string, v ...interface{})
	fatalfMutex       sync.RWMutex
	fatalfArgsForCall []struct {
		msg string
		v   []interface{}
	}
	TraceStub        func(msg string) *log.Entry
	traceMutex       sync.RWMutex
	traceArgsForCall []struct {
		msg string
	}
	traceReturns struct {
		result1 *log.Entry
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Logger) WithFields(fields log.Fielder) *log.Entry {
	fake.withFieldsMutex.Lock()
	fake.withFieldsArgsForCall = append(fake.withFieldsArgsForCall, struct {
		fields log.Fielder
	}{fields})
	fake.recordInvocation("WithFields", []interface{}{fields})
	fake.withFieldsMutex.Unlock()
	if fake.WithFieldsStub != nil {
		return fake.WithFieldsStub(fields)
	}
	return fake.withFieldsReturns.result1
}

func (fake *Logger) WithFieldsCallCount() int {
	fake.withFieldsMutex.RLock()
	defer fake.withFieldsMutex.RUnlock()
	return len(fake.withFieldsArgsForCall)
}

func (fake *Logger) WithFieldsArgsForCall(i int) log.Fielder {
	fake.withFieldsMutex.RLock()
	defer fake.withFieldsMutex.RUnlock()
	return fake.withFieldsArgsForCall[i].fields
}

func (fake *Logger) WithFieldsReturns(result1 *log.Entry) {
	fake.WithFieldsStub = nil
	fake.withFieldsReturns = struct {
		result1 *log.Entry
	}{result1}
}

func (fake *Logger) WithField(key string, value interface{}) *log.Entry {
	fake.withFieldMutex.Lock()
	fake.withFieldArgsForCall = append(fake.withFieldArgsForCall, struct {
		key   string
		value interface{}
	}{key, value})
	fake.recordInvocation("WithField", []interface{}{key, value})
	fake.withFieldMutex.Unlock()
	if fake.WithFieldStub != nil {
		return fake.WithFieldStub(key, value)
	}
	return fake.withFieldReturns.result1
}

func (fake *Logger) WithFieldCallCount() int {
	fake.withFieldMutex.RLock()
	defer fake.withFieldMutex.RUnlock()
	return len(fake.withFieldArgsForCall)
}

func (fake *Logger) WithFieldArgsForCall(i int) (string, interface{}) {
	fake.withFieldMutex.RLock()
	defer fake.withFieldMutex.RUnlock()
	return fake.withFieldArgsForCall[i].key, fake.withFieldArgsForCall[i].value
}

func (fake *Logger) WithFieldReturns(result1 *log.Entry) {
	fake.WithFieldStub = nil
	fake.withFieldReturns = struct {
		result1 *log.Entry
	}{result1}
}

func (fake *Logger) WithError(err error) *log.Entry {
	fake.withErrorMutex.Lock()
	fake.withErrorArgsForCall = append(fake.withErrorArgsForCall, struct {
		err error
	}{err})
	fake.recordInvocation("WithError", []interface{}{err})
	fake.withErrorMutex.Unlock()
	if fake.WithErrorStub != nil {
		return fake.WithErrorStub(err)
	}
	return fake.withErrorReturns.result1
}

func (fake *Logger) WithErrorCallCount() int {
	fake.withErrorMutex.RLock()
	defer fake.withErrorMutex.RUnlock()
	return len(fake.withErrorArgsForCall)
}

func (fake *Logger) WithErrorArgsForCall(i int) error {
	fake.withErrorMutex.RLock()
	defer fake.withErrorMutex.RUnlock()
	return fake.withErrorArgsForCall[i].err
}

func (fake *Logger) WithErrorReturns(result1 *log.Entry) {
	fake.WithErrorStub = nil
	fake.withErrorReturns = struct {
		result1 *log.Entry
	}{result1}
}

func (fake *Logger) Debug(msg string) {
	fake.debugMutex.Lock()
	fake.debugArgsForCall = append(fake.debugArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Debug", []interface{}{msg})
	fake.debugMutex.Unlock()
	if fake.DebugStub != nil {
		fake.DebugStub(msg)
	}
}

func (fake *Logger) DebugCallCount() int {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return len(fake.debugArgsForCall)
}

func (fake *Logger) DebugArgsForCall(i int) string {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return fake.debugArgsForCall[i].msg
}

func (fake *Logger) Info(msg string) {
	fake.infoMutex.Lock()
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Info", []interface{}{msg})
	fake.infoMutex.Unlock()
	if fake.InfoStub != nil {
		fake.InfoStub(msg)
	}
}

func (fake *Logger) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *Logger) InfoArgsForCall(i int) string {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return fake.infoArgsForCall[i].msg
}

func (fake *Logger) Warn(msg string) {
	fake.warnMutex.Lock()
	fake.warnArgsForCall = append(fake.warnArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Warn", []interface{}{msg})
	fake.warnMutex.Unlock()
	if fake.WarnStub != nil {
		fake.WarnStub(msg)
	}
}

func (fake *Logger) WarnCallCount() int {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	return len(fake.warnArgsForCall)
}

func (fake *Logger) WarnArgsForCall(i int) string {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	return fake.warnArgsForCall[i].msg
}

func (fake *Logger) Error(msg string) {
	fake.errorMutex.Lock()
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Error", []interface{}{msg})
	fake.errorMutex.Unlock()
	if fake.ErrorStub != nil {
		fake.ErrorStub(msg)
	}
}

func (fake *Logger) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *Logger) ErrorArgsForCall(i int) string {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return fake.errorArgsForCall[i].msg
}

func (fake *Logger) Fatal(msg string) {
	fake.fatalMutex.Lock()
	fake.fatalArgsForCall = append(fake.fatalArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Fatal", []interface{}{msg})
	fake.fatalMutex.Unlock()
	if fake.FatalStub != nil {
		fake.FatalStub(msg)
	}
}

func (fake *Logger) FatalCallCount() int {
	fake.fatalMutex.RLock()
	defer fake.fatalMutex.RUnlock()
	return len(fake.fatalArgsForCall)
}

func (fake *Logger) FatalArgsForCall(i int) string {
	fake.fatalMutex.RLock()
	defer fake.fatalMutex.RUnlock()
	return fake.fatalArgsForCall[i].msg
}

func (fake *Logger) Debugf(msg string, v ...interface{}) {
	fake.debugfMutex.Lock()
	fake.debugfArgsForCall = append(fake.debugfArgsForCall, struct {
		msg string
		v   []interface{}
	}{msg, v})
	fake.recordInvocation("Debugf", []interface{}{msg, v})
	fake.debugfMutex.Unlock()
	if fake.DebugfStub != nil {
		fake.DebugfStub(msg, v...)
	}
}

func (fake *Logger) DebugfCallCount() int {
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	return len(fake.debugfArgsForCall)
}

func (fake *Logger) DebugfArgsForCall(i int) (string, []interface{}) {
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	return fake.debugfArgsForCall[i].msg, fake.debugfArgsForCall[i].v
}

func (fake *Logger) Infof(msg string, v ...interface{}) {
	fake.infofMutex.Lock()
	fake.infofArgsForCall = append(fake.infofArgsForCall, struct {
		msg string
		v   []interface{}
	}{msg, v})
	fake.recordInvocation("Infof", []interface{}{msg, v})
	fake.infofMutex.Unlock()
	if fake.InfofStub != nil {
		fake.InfofStub(msg, v...)
	}
}

func (fake *Logger) InfofCallCount() int {
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	return len(fake.infofArgsForCall)
}

func (fake *Logger) InfofArgsForCall(i int) (string, []interface{}) {
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	return fake.infofArgsForCall[i].msg, fake.infofArgsForCall[i].v
}

func (fake *Logger) Warnf(msg string, v ...interface{}) {
	fake.warnfMutex.Lock()
	fake.warnfArgsForCall = append(fake.warnfArgsForCall, struct {
		msg string
		v   []interface{}
	}{msg, v})
	fake.recordInvocation("Warnf", []interface{}{msg, v})
	fake.warnfMutex.Unlock()
	if fake.WarnfStub != nil {
		fake.WarnfStub(msg, v...)
	}
}

func (fake *Logger) WarnfCallCount() int {
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	return len(fake.warnfArgsForCall)
}

func (fake *Logger) WarnfArgsForCall(i int) (string, []interface{}) {
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	return fake.warnfArgsForCall[i].msg, fake.warnfArgsForCall[i].v
}

func (fake *Logger) Errorf(msg string, v ...interface{}) {
	fake.errorfMutex.Lock()
	fake.errorfArgsForCall = append(fake.errorfArgsForCall, struct {
		msg string
		v   []interface{}
	}{msg, v})
	fake.recordInvocation("Errorf", []interface{}{msg, v})
	fake.errorfMutex.Unlock()
	if fake.ErrorfStub != nil {
		fake.ErrorfStub(msg, v...)
	}
}

func (fake *Logger) ErrorfCallCount() int {
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	return len(fake.errorfArgsForCall)
}

func (fake *Logger) ErrorfArgsForCall(i int) (string, []interface{}) {
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	return fake.errorfArgsForCall[i].msg, fake.errorfArgsForCall[i].v
}

func (fake *Logger) Fatalf(msg string, v ...interface{}) {
	fake.fatalfMutex.Lock()
	fake.fatalfArgsForCall = append(fake.fatalfArgsForCall, struct {
		msg string
		v   []interface{}
	}{msg, v})
	fake.recordInvocation("Fatalf", []interface{}{msg, v})
	fake.fatalfMutex.Unlock()
	if fake.FatalfStub != nil {
		fake.FatalfStub(msg, v...)
	}
}

func (fake *Logger) FatalfCallCount() int {
	fake.fatalfMutex.RLock()
	defer fake.fatalfMutex.RUnlock()
	return len(fake.fatalfArgsForCall)
}

func (fake *Logger) FatalfArgsForCall(i int) (string, []interface{}) {
	fake.fatalfMutex.RLock()
	defer fake.fatalfMutex.RUnlock()
	return fake.fatalfArgsForCall[i].msg, fake.fatalfArgsForCall[i].v
}

func (fake *Logger) Trace(msg string) *log.Entry {
	fake.traceMutex.Lock()
	fake.traceArgsForCall = append(fake.traceArgsForCall, struct {
		msg string
	}{msg})
	fake.recordInvocation("Trace", []interface{}{msg})
	fake.traceMutex.Unlock()
	if fake.TraceStub != nil {
		return fake.TraceStub(msg)
	}
	return fake.traceReturns.result1
}

func (fake *Logger) TraceCallCount() int {
	fake.traceMutex.RLock()
	defer fake.traceMutex.RUnlock()
	return len(fake.traceArgsForCall)
}

func (fake *Logger) TraceArgsForCall(i int) string {
	fake.traceMutex.RLock()
	defer fake.traceMutex.RUnlock()
	return fake.traceArgsForCall[i].msg
}

func (fake *Logger) TraceReturns(result1 *log.Entry) {
	fake.TraceStub = nil
	fake.traceReturns = struct {
		result1 *log.Entry
	}{result1}
}

func (fake *Logger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.withFieldsMutex.RLock()
	defer fake.withFieldsMutex.RUnlock()
	fake.withFieldMutex.RLock()
	defer fake.withFieldMutex.RUnlock()
	fake.withErrorMutex.RLock()
	defer fake.withErrorMutex.RUnlock()
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	fake.fatalMutex.RLock()
	defer fake.fatalMutex.RUnlock()
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	fake.warnfMutex.RLock()
	defer fake.warnfMutex.RUnlock()
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	fake.fatalfMutex.RLock()
	defer fake.fatalfMutex.RUnlock()
	fake.traceMutex.RLock()
	defer fake.traceMutex.RUnlock()
	return fake.invocations
}

func (fake *Logger) recordInvocation(key string, args []interface{}) {
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

var _ log.Interface = new(Logger)