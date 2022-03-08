package sta17_test

import (
	"reflect"
	"testing"

	"github.com/osamingo/go-todo-app/model"
)

func TestStation17(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target       interface{}
		FieldName    string
		WantKinds    []reflect.Kind
		WantTypes    []reflect.Type
		JSONTagValue string
	}{
		"UpdateTODORequest has ID field": {
			Target:    model.DeleteTODORequest{},
			FieldName: "IDs",
			WantKinds: []reflect.Kind{reflect.Slice},
			WantTypes: []reflect.Type{
				reflect.SliceOf(reflect.TypeOf(int(0))),
				reflect.SliceOf(reflect.TypeOf(uint(0))),
				reflect.SliceOf(reflect.TypeOf(int8(0))),
				reflect.SliceOf(reflect.TypeOf(uint8(0))),
				reflect.SliceOf(reflect.TypeOf(int16(0))),
				reflect.SliceOf(reflect.TypeOf(uint16(0))),
				reflect.SliceOf(reflect.TypeOf(int32(0))),
				reflect.SliceOf(reflect.TypeOf(uint32(0))),
				reflect.SliceOf(reflect.TypeOf(int64(0))),
				reflect.SliceOf(reflect.TypeOf(uint64(0))),
			},
			JSONTagValue: "ids",
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
				t.Error(tc.FieldName + " field がみつかりません")
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

			if tc.WantTypes != nil {
				notFound = true
				for _, et := range tc.WantTypes {
					if f.Type == et {
						notFound = false
						break
					}
				}
				if notFound {
					t.Errorf(tc.FieldName+" が期待している Type ではありません, got = %s, want = %s", f.Type, tc.WantTypes)
					return
				}
			}

			v, ok := f.Tag.Lookup("json")
			if !ok {
				t.Error("json tag が見つかりません")
				return
			}

			if v != tc.JSONTagValue {
				t.Errorf("json tag の内容が期待している内容ではありません, got = %s, want = %s", v, tc.JSONTagValue)
				return
			}
		})
	}
}
