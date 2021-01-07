package astilog

import (
	"context"
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {
	fs1 := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	ctx1 := ContextWithFields(context.Background(), fs1)
	fs2 := map[string]interface{}{
		"k3": "v3",
	}
	fs := make(map[string]interface{})
	for k, v := range fs1 {
		fs[k] = v
	}
	for k, v := range fs2 {
		fs[k] = v
	}
	ctx2 := ContextWithFields(ctx1, fs2)
	if g := FieldsFromContext(ctx1); !reflect.DeepEqual(fs1, g) {
		t.Errorf("expected %+v, got %+v", fs1, g)
	}
	if g := FieldsFromContext(ctx2); !reflect.DeepEqual(fs, g) {
		t.Errorf("expected %+v, got %+v", fs, g)
	}
}
