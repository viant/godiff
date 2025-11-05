package godiff

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BoolOnly struct {
	Flag bool
}

type OuterBool struct {
	Inner *BoolOnly
}

func TestBoolField_UpdateFalseToTrue(t *testing.T) {
	from := &BoolOnly{Flag: false}
	to := &BoolOnly{Flag: true}

	differ, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
	if !assert.NoError(t, err) {
		return
	}
	got := differ.Diff(from, to)

	expect := &ChangeLog{Changes: []*Change{
		{Type: ChangeTypeUpdate, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Flag"}, From: false, To: true},
	}}
	assert.EqualValues(t, expect, got)
}

func TestBoolField_UpdateTrueToFalse(t *testing.T) {
	from := &BoolOnly{Flag: true}
	to := &BoolOnly{Flag: false}

	differ, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
	if !assert.NoError(t, err) {
		return
	}
	got := differ.Diff(from, to)

	expect := &ChangeLog{Changes: []*Change{
		{Type: ChangeTypeUpdate, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Flag"}, From: true, To: false},
	}}
	assert.EqualValues(t, expect, got)
}

func TestBoolField_CreateWithFalse(t *testing.T) {
	var from *BoolOnly = nil
	to := &BoolOnly{Flag: false}

	differ, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
	if !assert.NoError(t, err) {
		return
	}
	got := differ.Diff(from, to)

	expect := &ChangeLog{Changes: []*Change{
		// Library represents field-level as update with From=nil when parent was created
		{Type: ChangeTypeUpdate, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Flag"}, From: nil, To: false},
	}}
	assert.EqualValues(t, expect, got)
}

func TestBoolField_DeleteWithFalse(t *testing.T) {
	from := &BoolOnly{Flag: false}
	var to *BoolOnly = nil

	differ, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
	if !assert.NoError(t, err) {
		return
	}
	got := differ.Diff(from, to)

	expect := &ChangeLog{Changes: []*Change{
		// Library represents field-level as update with To=nil when parent was deleted
		{Type: ChangeTypeUpdate, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Flag"}, From: false, To: nil},
	}}
	assert.EqualValues(t, expect, got)
}

func TestNestedBool_NoPanicAndCreateFalse(t *testing.T) {
	from := &OuterBool{Inner: nil}
	to := &OuterBool{Inner: &BoolOnly{Flag: false}}

	differ, err := New(reflect.TypeOf(from), reflect.TypeOf(to))
	if !assert.NoError(t, err) {
		return
	}

	// Ensure no panic and correct create on nested bool field
	got := differ.Diff(from, to)
	expect := &ChangeLog{Changes: []*Change{
		// When parent Inner is created, nested field is a create as well
		{Type: ChangeTypeCreate, Path: &Path{Kind: PathKinField, Path: &Path{Kind: PathKinField, Path: &Path{}, Name: "Inner"}, Name: "Flag"}, To: false},
	}}
	assert.EqualValues(t, expect, got)
}
