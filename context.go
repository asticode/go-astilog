package astilog

import (
	"context"
	"sync"
)

const contextKeyFields = "astilog.fields"

type contextFields struct {
	fs map[string]interface{}
	m  *sync.Mutex
}

func newContextFields() *contextFields {
	return &contextFields{
		fs: make(map[string]interface{}),
		m:  &sync.Mutex{},
	}
}

func fieldsFromContext(ctx context.Context) *contextFields {
	v, ok := ctx.Value(contextKeyFields).(*contextFields)
	if !ok {
		return nil
	}
	return v
}

func ContextWithField(ctx context.Context, k string, v interface{}) context.Context {
	return ContextWithFields(ctx, map[string]interface{}{k: v})
}

func ContextWithFields(ctx context.Context, fs map[string]interface{}) context.Context {
	cfs := fieldsFromContext(ctx)
	if cfs == nil {
		cfs = newContextFields()
	}
	cfs.m.Lock()
	for k, v := range fs {
		cfs.fs[k] = v
	}
	cfs.m.Unlock()
	return context.WithValue(ctx, contextKeyFields, cfs)
}
