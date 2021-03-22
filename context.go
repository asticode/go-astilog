package astilog

import (
	"context"
	"sync"
)

type contextKey string

const contextKeyFields contextKey = "astilog.fields"

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
	if ctx == nil {
		return nil
	}
	v, ok := ctx.Value(contextKeyFields).(*contextFields)
	if !ok {
		return nil
	}
	return v
}

func FieldsFromContext(ctx context.Context) (fs map[string]interface{}) {
	if cfs := fieldsFromContext(ctx); cfs != nil {
		cfs.m.Lock()
		fs = make(map[string]interface{})
		for k, v := range cfs.fs {
			fs[k] = v
		}
		cfs.m.Unlock()
		return
	}
	return
}

func ContextWithField(ctx context.Context, k string, v interface{}) context.Context {
	return ContextWithFields(ctx, map[string]interface{}{k: v})
}

func ContextWithFields(ctx context.Context, fs map[string]interface{}) context.Context {
	cfs := newContextFields()
	if ccfs := fieldsFromContext(ctx); ccfs != nil {
		ccfs.m.Lock()
		for k, v := range ccfs.fs {
			cfs.fs[k] = v
		}
		ccfs.m.Unlock()
	}
	for k, v := range fs {
		cfs.fs[k] = v
	}
	return context.WithValue(ctx, contextKeyFields, cfs)
}
