package sta6_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func TestStation6(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target    interface{}
		FieldName string
		WantKinds []reflect.Kind
		WantType  reflect.Type
	}{
		"TODO has ID field": {
			Target:    model.TODO{},
			FieldName: "ID",
			WantKinds: []reflect.Kind{reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16,
				reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64},
		},
		"TODO has Subject field": {
			Target:    model.TODO{},
			FieldName: "Subject",
			WantKinds: []reflect.Kind{reflect.String},
		},
		"TODO has Description field": {
			Target:    model.TODO{},
			FieldName: "Description",
			WantKinds: []reflect.Kind{reflect.String},
		},
		"TODO has CreatedAt field": {
			Target:    model.TODO{},
			FieldName: "CreatedAt",
			WantKinds: []reflect.Kind{reflect.Struct},
			WantType:  reflect.TypeOf(time.Time{}),
		},
		"TODO has UpdatedAt field": {
			Target:    model.TODO{},
			FieldName: "UpdatedAt",
			WantKinds: []reflect.Kind{reflect.Struct},
			WantType:  reflect.TypeOf(time.Time{}),
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

			notFound := true
			for _, k := range tc.WantKinds {
				if f.Type.Kind() == k {
					notFound = false
					break
				}
			}
			if notFound {
				t.Errorf(tc.FieldName+" が期待している kind ではありません, got = %s, want = %s", f.Type.Kind(), tc.WantKinds)
				return
			}

			if tc.WantType == nil {
				return
			}

			if f.Type != tc.WantType {
				t.Errorf(tc.FieldName+" が期待している Type ではありません, got = %s, want = %s", f.Type, tc.WantType)
				return
			}
		})
	}
}
