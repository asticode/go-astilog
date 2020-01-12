package astilog

import (
	"context"
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {
	fs := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	ctx := ContextWithFields(context.Background(), fs)
	if g := fieldsFromContext(ctx); !reflect.DeepEqual(fs, g) {
		t.Errorf("expected %+v, got %+v", fs, g)
	}
}
