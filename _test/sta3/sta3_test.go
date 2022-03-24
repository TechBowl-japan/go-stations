package sta3_test

import (
	"reflect"
	"testing"

	"github.com/TechBowl-japan/go-stations/model"
)

func TestStation3(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target    model.HealthzResponse
		FieldName string
		WantKind  reflect.Kind
	}{
		"HealthzResponse has Message field": {
			Target:    model.HealthzResponse{},
			FieldName: "Message",
			WantKind:  reflect.String,
		},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tp := reflect.TypeOf(tc.Target)
			f, ok := tp.FieldByName(tc.FieldName)
			if !ok {
				t.Error(tc.FieldName + " field が見つかりません")
				return
			}

			if f.Type.Kind() == tc.WantKind {
				return
			}
			t.Errorf(tc.FieldName+" が期待している kind ではありません, got = %s, want = %s", f.Type.Kind(), tc.WantKind)
		})
	}
}
