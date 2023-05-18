package godiff

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestNewDiffer(t *testing.T) {

	type PresRecordHas struct {
		ID         bool
		Name       bool
		ScratchPad bool
		Time       bool
	}
	type PresRecord struct {
		ID         int
		Name       string
		ScratchPad string `diff:"-"`
		Time       time.Time
		Bar        int
		Has        *PresRecordHas `diff:"presence=true"`
	}

	type Record struct {
		ID         int
		Name       string
		ScratchPad string `diff:"-"`
		Time       time.Time
	}
	type Nums struct {
		Ints     []int     `diff:"sort=true"`
		Numerics []float32 `diff:"indexBy=."`
	}

	type IRecrod struct {
		ID    int
		Value interface{}
	}

	type Repeated struct {
		ID      int `diff:"name=ID"`
		Nums    []int
		Records []*Record
		Ifs     []interface{}
		Map     map[string]interface{}
	}
	type Nested struct {
		ID     int
		Record *Record
	}
	now := time.Now()
	inFuture := now.Add(time.Hour)
	type RecordPtr struct {
		ID   *int
		Name *string
		Time *time.Time
	}

	type Custom struct {
		ID       int
		Expr     string `diff:"pairdelimiter=AND|OR,pairSeparator=:,whitespace=]|[|\""`
		List     string `diff:"itemSeparator=$coma"`
		ExprList string `diff:"pairdelimiter=AND|OR,pairSeparator=:,whitespace=]|[|\",itemSeparator=$coma,indexBy=."`
	}

	var testCases = []struct {
		description string
		options     []Option
		from        interface{}
		to          interface{}
		expect      *ChangeLog
	}{
		{
			description: "basic struct - delete",
			from:        &Record{ID: 123, Name: "abc"},
			to:          nil,
			expect: &ChangeLog{Changes: []*Change{
				{Type: "delete", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ID"}, From: 123},
				{Type: "delete", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Name"}, From: "abc"},
			}},
		},
		{
			description: "basic struct - update",
			from:        &Record{ID: 123, Name: "abc"},
			to:          &Record{ID: 323, Name: "xtz"},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ID"}, From: 123, To: 323},
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Name"}, From: "abc", To: "xtz"},
			}},
		},
		{
			description: "basic struct - create",
			from:        nil,
			to:          &Record{ID: 323, Name: "az"},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "create", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ID"}, To: 323},
				{Type: "create", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Name"}, To: "az"},
			}},
		},
		{
			description: "basic struct ptr - update",
			from:        &RecordPtr{ID: intPtr(123), Name: stringPtr("abc")},
			to:          &RecordPtr{ID: intPtr(323), Name: stringPtr("xtz")},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ID"}, From: 123, To: 323},
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Name"}, From: "abc", To: "xtz"},
			}},
		},
		{
			description: "basic struct ptr - time",
			from:        &RecordPtr{ID: intPtr(123), Time: &now},
			to:          &RecordPtr{ID: intPtr(123), Time: &inFuture},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Time"}, From: now, To: inFuture},
			}},
		},

		{
			description: "nested - update",
			from:        &Nested{ID: 1, Record: &Record{ID: 123, Name: "abc"}},
			to:          &Nested{ID: 1, Record: &Record{ID: 323, Name: "xtz"}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{Kind: PathKinField, Name: "Record", Path: &Path{}}, Name: "ID"}, From: 123, To: 323},
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{Kind: PathKinField, Name: "Record", Path: &Path{}}, Name: "Name"}, From: "abc", To: "xtz"},
			}},
		},

		{
			description: "repeated - update",
			from:        &Repeated{ID: 1, Records: []*Record{{ID: 12}}, Nums: []int{10, 2}},
			to:          &Repeated{ID: 1, Records: []*Record{{ID: 23}}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "create", Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Nums"}, Index: 0}, To: 10},
				{Type: "create", Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Nums"}, Index: 1}, To: 2},
				{Type: "update", Path: &Path{Kind: 1, Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Records"}, Index: 0}, Name: "ID"}, From: 12, To: 23},
			}},
		},
		{
			description: "repeated - delete",
			from:        &Repeated{ID: 1, Records: []*Record{{ID: 23}}},
			to:          &Repeated{ID: 1},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "delete", Path: &Path{Kind: 1, Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Records"}, Index: 0}, Name: "ID"}, From: 23},
			}},
		},

		{
			description: "repeated - sort primitive",
			from:        &Nums{Ints: []int{4, 6, 1}},
			to:          &Nums{Ints: []int{1, 6, 4, 7}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "create", Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Ints"}, Index: 3}, From: nil, To: 7},
			}},
		},
		{
			description: "repeated - index",
			from:        &Nums{Numerics: []float32{4, 6, 1}},
			to:          &Nums{Numerics: []float32{1, 2, 6, 4}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "create", Path: &Path{Kind: 3, Path: &Path{Kind: 1, Path: &Path{}, Name: "Numerics"}, Index: 1}, From: nil, To: float32(2.0)},
			}},
		},

		{
			description: "interface field",
			from:        &IRecrod{ID: 1, Value: &Record{Name: "test 1"}},
			to:          &IRecrod{ID: 1, Value: Record{Name: "test 2"}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Value"}, Name: "Name"}, From: "test 1", To: "test 2"},
			}},
		},
		{
			description: "[]interface field",
			from:        &Repeated{Ifs: []interface{}{&Record{Name: "test 1"}}},
			to:          &Repeated{Ifs: []interface{}{Record{Name: "test 2"}}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{Kind: PathKindIndex, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Ifs"}}, Name: "Name"}, From: "test 1", To: "test 2"},
			}},
		},
		{
			description: "map field",
			from: &Repeated{Map: map[string]interface{}{
				"k1": "123",
				"k2": "abc",
			}},
			to: &Repeated{Map: map[string]interface{}{
				"k1": "123",
				"k2": "xyz",
			}},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKindKey, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Map"}, Key: "k2"}, From: "abc", To: "xyz"},
			}},
		},

		{
			description: "item separator",
			from:        &Custom{List: "1,2,3"},
			to:          &Custom{List: "1,2,4"},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKindIndex, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "List"}, Index: 2}, From: "3", To: "4"},
			}},
		},

		{
			description: "pair delimiter",
			from:        &Custom{Expr: "k1:v1 AND k2:v2 OR k3:v3"},
			to:          &Custom{Expr: "k1:v1 AND k2:v2 OR k3:v3.1 AND k4:v4"},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKindKey, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Expr"}, Key: "k3"}, From: "v3", To: "v3.1"},
				{Type: "create", Path: &Path{Kind: PathKindKey, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Expr"}, Key: "k4"}, To: "v4"},
			}},
		},
		{
			description: "pair delimiter with list",
			from:        &Custom{ExprList: "k1:[1,2] AND k2:[v2] "},
			to:          &Custom{ExprList: "k1:[0,1] AND k2:[v2]"},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "delete", Path: &Path{Kind: PathKindIndex, Path: &Path{Kind: PathKindKey, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ExprList"}, Key: "k1"}, Index: 1}, From: "2"},
				{Type: "create", Path: &Path{Kind: PathKindIndex, Path: &Path{Kind: PathKindKey, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "ExprList"}, Key: "k1"}, Index: 0}, To: "0"},
			}},
		},
		{description: "diff with field presence check",
			from: &PresRecord{ID: 20, Name: "abc", Bar: 23,
				Has: &PresRecordHas{
					Name: true,
					ID:   true,
				}},
			to: &PresRecord{ID: 21, Name: "xyz", ScratchPad: "xx",
				Has: &PresRecordHas{
					Name:       true,
					ScratchPad: true,
				},
			},
			options: []Option{WithPresence(true)},
			expect: &ChangeLog{Changes: []*Change{
				{Type: "update", Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Name"}, From: "abc", To: "xyz"},
			}},
		},
	}
	for _, testCase := range testCases {
		differ, err := New(reflect.TypeOf(testCase.from), reflect.TypeOf(testCase.to), testCase.options...)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		changeLog := differ.Diff(testCase.from, testCase.to)
		if !assert.EqualValues(t, testCase.expect, changeLog, testCase.description) {
			data, _ := json.Marshal(changeLog)
			fmt.Println(string(data))
		}
	}
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
