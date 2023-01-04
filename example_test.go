package godiff_test

import (
	"fmt"
	"github.com/viant/godiff"
	"log"
	"reflect"
)

//ExampleNew shows basic differ usage
func ExampleNew() {

	type Flag struct {
		Value int
	}
	type Record struct {
		Id        int
		Name      string
		Dep       *Record
		Transient string `diff:"-"`
		Flags     []Flag
	}
	record1 := &Record{
		Id:    1,
		Name:  "Rec1",
		Flags: []Flag{{Value: 3}, {Value: 15}},
	}

	record2 := &Record{
		Id:    2,
		Name:  "Rec1",
		Dep:   &Record{Id: 10},
		Flags: []Flag{{Value: 12}},
	}

	diff, err := godiff.New(reflect.TypeOf(record1), reflect.TypeOf(record2))
	if err != nil {
		log.Fatal(err)
	}
	changeLog := diff.Diff(record1, record2)
	fmt.Println(changeLog.String())
}
