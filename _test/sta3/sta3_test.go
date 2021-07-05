package sta3_test

import (
	"reflect"
	"testing"

	"github.com/TechBowl-japan/go-stations/model"
)

func TestStation3(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target    interface{}
		FieldName string
		WantKinds []reflect.Kind
	}{
		"HealthzResponse has Message field": {
			Target:    model.HealthzResponse{},
			FieldName: "Message",
			WantKinds: []reflect.Kind{reflect.String},
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tp := reflect.TypeOf(tc.Target)
			f, ok := tp.FieldByName(tc.FieldName)
			if !ok {
				t.Error(tc.FieldName + " field が見つかりません")
				return
			}

			for _, k := range tc.WantKinds {
				if f.Type.Kind() == k {
					return
				}
			}
			t.Errorf(tc.FieldName+" が期待している kind ではありません, got = %s, want = %s", f.Type.Kind(), tc.WantKinds)
		})
	}
}
