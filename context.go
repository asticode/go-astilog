package astilog

import "context"

const contextKeyFields = "astilog.fields"

func fieldsFromContext(ctx context.Context) map[string]interface{} {
	v, ok := ctx.Value(contextKeyFields).(map[string]interface{})
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
		cfs = make(map[string]interface{})
	}
	for k, v := range fs {
		cfs[k] = v
	}
	return context.WithValue(ctx, contextKeyFields, cfs)
}
