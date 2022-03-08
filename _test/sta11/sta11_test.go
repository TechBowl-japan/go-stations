package sta11_test

import (
	"reflect"
	"testing"

	"github.com/TechBowl-japan/go-stations/model"
)

func TestStation11(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target       interface{}
		FieldName    string
		WantKinds    []reflect.Kind
		WantTypes    []reflect.Type
		JSONTagValue string
	}{
		"UpdateTODORequest has ID field": {
			Target:    model.UpdateTODORequest{},
			FieldName: "ID",
			WantKinds: []reflect.Kind{reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16,
				reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64},
			JSONTagValue: "id",
		},
		"UpdateTODORequest has Subject field": {
			Target:       model.UpdateTODORequest{},
			FieldName:    "Subject",
			WantKinds:    []reflect.Kind{reflect.String},
			JSONTagValue: "subject",
		},
		"UpdateTODORequest has Description field": {
			Target:       model.UpdateTODORequest{},
			FieldName:    "Description",
			WantKinds:    []reflect.Kind{reflect.String},
			JSONTagValue: "description",
		},
		"UpdateTODOResponse has TODO field": {
			Target:       model.UpdateTODOResponse{},
			FieldName:    "TODO",
			WantKinds:    []reflect.Kind{reflect.Struct, reflect.Ptr},
			WantTypes:    []reflect.Type{reflect.TypeOf(model.TODO{}), reflect.TypeOf(&model.TODO{})},
			JSONTagValue: "todo",
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
