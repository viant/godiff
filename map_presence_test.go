package godiff

import (
	"reflect"
	"testing"
)

func TestMapKeyPresenceChangeKinds(t *testing.T) {
	typ := reflect.TypeOf(map[string]interface{}{})
	for _, tc := range []struct {
		name     string
		from, to interface{}
		want     map[string]ChangeType
	}{
		{"removed", map[string]interface{}{"gone": 2}, map[string]interface{}{}, map[string]ChangeType{"[gone]": ChangeTypeDelete}},
		{"added nil", map[string]interface{}{}, map[string]interface{}{"new": nil}, map[string]ChangeType{"[new]": ChangeTypeCreate}},
		{"removed nil", map[string]interface{}{"gone": nil}, map[string]interface{}{}, map[string]ChangeType{"[gone]": ChangeTypeDelete}},
		{"unchanged nil", map[string]interface{}{"same": nil}, map[string]interface{}{"same": nil}, map[string]ChangeType{}},
		{"value to nil", map[string]interface{}{"same": 2}, map[string]interface{}{"same": nil}, map[string]ChangeType{"[same]": ChangeTypeUpdate}},
		{"nil root create", nil, map[string]interface{}{"new": 2, "empty": nil}, map[string]ChangeType{"[new]": ChangeTypeCreate, "[empty]": ChangeTypeCreate}},
		{"nil root delete", map[string]interface{}{"gone": 2, "empty": nil}, nil, map[string]ChangeType{"[gone]": ChangeTypeDelete, "[empty]": ChangeTypeDelete}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := New(typ, typ)
			if err != nil {
				t.Fatal(err)
			}
			got := map[string]ChangeType{}
			for _, change := range d.Diff(tc.from, tc.to).Changes {
				if change.Error != "" {
					t.Fatal(change.Error)
				}
				got[change.Path.String()] = change.Type
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got=%v want=%v", got, tc.want)
			}
		})
	}
}
