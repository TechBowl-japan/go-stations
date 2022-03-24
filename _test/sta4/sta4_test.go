package sta4_test

import (
	"reflect"
	"testing"

	"github.com/TechBowl-japan/go-stations/model"
)

func TestStation4(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target       model.HealthzResponse
		FieldName    string
		WantKind     reflect.Kind
		JSONTagValue string
	}{
		"HealthzResponse has Message field": {
			Target:       model.HealthzResponse{},
			FieldName:    "Message",
			WantKind:     reflect.String,
			JSONTagValue: "message",
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
			if f.Type.Kind() != tc.WantKind {
				t.Errorf(tc.FieldName+" が期待している kind ではありませ, got = %s, want = %s", f.Type.Kind(), tc.WantKind)
			}
			v, ok := f.Tag.Lookup("json")
			if !ok {
				t.Error("json tag が見つかりません")
				return
			}
			if v != tc.JSONTagValue {
				t.Errorf("json tag の内容が期待している内容ではありません, got = %s, want = %s", v, tc.JSONTagValue)
			}
		})
	}
}
