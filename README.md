# godiff (Data structure diff for GoLang)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/godiff)](https://goreportcard.com/report/github.com/viant/godiff)
[![GoDoc](https://godoc.org/github.com/viant/godiff?status.svg)](https://godoc.org/github.com/viant/godiff)

This library is compatible with Go 1.17+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Motivation](#motivation)
- [Usage](#usage)
- [Contribution](#contributing-to-godiff)
- [License](#license)

## Motivation

The goal of this library is to both produce diff and patch
for arbitrary data structure in GoLang in performant manner

Other project for diffing golang structures and values:

- [r3labs/diff](https://github.com/r3labs/diff)

See 3rd party differ performance comparison.

- [Benchmark](#benchmark)

## Usage

```go
package godiff_test

import (
	"fmt"
	"github.com/viant/godiff"
	"log"
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

//ExampleNew shows basic differ usage
func ExampleNew() {

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
```

Supported [tags](tag.go):

- name - optional name in the change log
- indexBy - index elements before comparing
- sort - sort elements before comparing
- whitespace - remove specified whitespace chars when converting string to list or map
- pairSeparator - pair separator to convert string to a map comparison
- pairDelimiter - pair delimiter
- itemSeparator - item separator to convert string to a slice comparison
- '-' (ignore)

Work in progress tag:

- timeLayout
- precision

## Config option
- WithTagName
- WithRegistry
- NullifyEmpty
- WithConfig

## Diff option
- WithPresence
- WithShallow

## Benchmark

GoDiff is around 5x faster than s3lab diff.

```text
pkg: github.com/viant/godiff/bench
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
Benchmark_GoDiff
Benchmark_GoDiff-16       	 1432188	       815.9 ns/op	     768 B/op	      14 allocs/op
Benchmark_S3LabDiff
Benchmark_S3LabDiff-16    	  288633	      4531 ns/op	    2262 B/op	      42 allocs/op
```

## Contributing to godiff

godiff is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.



## Credits and Acknowledgements

**Library Author:** Adrian Witas

