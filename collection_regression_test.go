package godiff

import (
	"reflect"
	"testing"
)

func TestCollectionNilPointersAndDeletion(t *testing.T) {
	var nilInt *int
	value := 2
	for _, tc := range []struct {
		name     string
		from, to interface{}
		want     map[string]ChangeType
	}{
		{"typed nil map deletion", map[string]interface{}{"gone": nilInt}, map[string]interface{}{}, map[string]ChangeType{"[gone]": ChangeTypeDelete}},
		{"pointer map update", map[string]interface{}{"a": nilInt}, map[string]interface{}{"a": &value}, map[string]ChangeType{"[a]": ChangeTypeUpdate}},
		{"typed slice shrink", []int{1, 2}, []int{1}, map[string]ChangeType{"[1]": ChangeTypeDelete}},
		{"typed slice empty", []int{1}, []int{}, map[string]ChangeType{"[0]": ChangeTypeDelete}},
		{"interface slice shrink", []interface{}{1, 2}, []interface{}{1}, map[string]ChangeType{"[1]": ChangeTypeDelete}},
		{"interface slice grow nil", []interface{}{1}, []interface{}{1, nilInt}, map[string]ChangeType{"[1]": ChangeTypeCreate}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if v := recover(); v != nil {
					t.Fatalf("panic: %v", v)
				}
			}()
			d, err := New(reflect.TypeOf(tc.from), reflect.TypeOf(tc.to))
			if err != nil {
				t.Fatal(err)
			}
			got := map[string]ChangeType{}
			for _, c := range d.Diff(tc.from, tc.to).Changes {
				if c.Error != "" {
					t.Fatal(c.Error)
				}
				got[c.Path.String()] = c.Type
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestInterfaceConcreteTypeChanges(t *testing.T) {
	for _, pair := range [][2]interface{}{{1, "one"}, {"one", 1}, {[]int{1}, "one"}, {nil, "one"}} {
		from := map[string]interface{}{"value": pair[0]}
		to := map[string]interface{}{"value": pair[1]}
		d, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
		if err != nil {
			t.Fatal(err)
		}
		changes := d.Diff(from, to).Changes
		if len(changes) != 1 || changes[0].Error != "" || changes[0].Type != ChangeTypeUpdate || !reflect.DeepEqual(changes[0].From, pair[0]) || !reflect.DeepEqual(changes[0].To, pair[1]) {
			t.Fatalf("pair=%v changes=%+v", pair, changes)
		}
	}
}
