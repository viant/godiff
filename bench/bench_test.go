package bench_test

import (
	//	diff "github.com/r3labs/diff"
	"github.com/stretchr/testify/assert"
	"github.com/viant/godiff"

	"reflect"
	"testing"
)

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

var benchDiff, _ = godiff.New(reflect.TypeOf(&Record{}), reflect.TypeOf(&Record{}))

func Benchmark_GoDiff(b *testing.B) {

	record1 := &Record{
		Id:   1,
		Name: "Rec1",
	}
	record2 := &Record{
		Id:    2,
		Name:  "Rec1",
		Dep:   &Record{Id: 10},
		Flags: []Flag{{Value: 12}},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		changeLog := benchDiff.Diff(record1, record2)
		assert.True(b, len(changeLog.Changes) == 3)
	}
}

/*
var s3LabDif, _ = diff.NewDiffer()

func Benchmark_S3LabDiff(b *testing.B) {

	record1 := &Record{
		ID:   1,
		Name: "Rec1",
	}
	record2 := &Record{
		ID:    2,
		Name:  "Rec1",
		Dep:   &Record{ID: 10},
		Flags: []Flag{{Value: 12}},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		changeLog, _ := s3LabDif.Diff(record1, record2)
		assert.True(b, changeLog != nil)
	}
}


*/
