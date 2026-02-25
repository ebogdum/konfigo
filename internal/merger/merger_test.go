package merger

import (
	"reflect"
	"testing"
)

func TestMergeSlices_BasicUnionWithDedup(t *testing.T) {
	dst := []interface{}{"a", "b"}
	src := []interface{}{"b", "c"}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_NoDuplicatesInSrc(t *testing.T) {
	dst := []interface{}{"a"}
	src := []interface{}{"b", "c"}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_AllDuplicates(t *testing.T) {
	dst := []interface{}{"a", "b"}
	src := []interface{}{"a", "b"}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_EmptyDst(t *testing.T) {
	dst := []interface{}{}
	src := []interface{}{"a", "b"}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_EmptySrc(t *testing.T) {
	dst := []interface{}{"a", "b"}
	src := []interface{}{}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_BothEmpty(t *testing.T) {
	dst := []interface{}{}
	src := []interface{}{}
	got := mergeSlices(dst, src)
	if len(got) != 0 {
		t.Errorf("mergeSlices([], []) = %v, want []", got)
	}
}

func TestMergeSlices_MixedTypes(t *testing.T) {
	dst := []interface{}{"a", 1, true}
	src := []interface{}{1, "b", false}
	got := mergeSlices(dst, src)
	want := []interface{}{"a", 1, true, "b", false}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices(%v, %v) = %v, want %v", dst, src, got, want)
	}
}

func TestMergeSlices_NestedMaps(t *testing.T) {
	m1 := map[string]interface{}{"host": "a"}
	m2 := map[string]interface{}{"host": "b"}
	m1dup := map[string]interface{}{"host": "a"}
	dst := []interface{}{m1}
	src := []interface{}{m1dup, m2}
	got := mergeSlices(dst, src)
	want := []interface{}{m1, m2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mergeSlices with nested maps: got %v, want %v", got, want)
	}
}

func TestMerge_ArraysReplacedWhenFlagFalse(t *testing.T) {
	dst := map[string]interface{}{
		"tags": []interface{}{"a", "b"},
	}
	src := map[string]interface{}{
		"tags": []interface{}{"c", "d"},
	}
	Merge(dst, src, true, nil, false)
	got := dst["tags"]
	want := []interface{}{"c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge with mergeArrays=false: tags = %v, want %v", got, want)
	}
}

func TestMerge_ArraysMergedWhenFlagTrue(t *testing.T) {
	dst := map[string]interface{}{
		"tags": []interface{}{"a", "b"},
	}
	src := map[string]interface{}{
		"tags": []interface{}{"b", "c"},
	}
	Merge(dst, src, true, nil, true)
	got := dst["tags"]
	want := []interface{}{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge with mergeArrays=true: tags = %v, want %v", got, want)
	}
}

func TestMerge_ArraysMergedCaseInsensitive(t *testing.T) {
	dst := map[string]interface{}{
		"Tags": []interface{}{"a", "b"},
	}
	src := map[string]interface{}{
		"tags": []interface{}{"b", "c"},
	}
	Merge(dst, src, false, nil, true)
	got := dst["tags"]
	want := []interface{}{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge case-insensitive with mergeArrays=true: tags = %v, want %v", got, want)
	}
}

func TestMerge_MapsStillDeepMergeRegardless(t *testing.T) {
	dst := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	src := map[string]interface{}{
		"db": map[string]interface{}{"host": "prod.example.com"},
	}
	Merge(dst, src, true, nil, true)
	gotDB := dst["db"].(map[string]interface{})
	if gotDB["host"] != "prod.example.com" {
		t.Errorf("expected host=prod.example.com, got %v", gotDB["host"])
	}
	if gotDB["port"] != 5432 {
		t.Errorf("expected port=5432, got %v", gotDB["port"])
	}
}

func TestMerge_NestedArraysMerged(t *testing.T) {
	dst := map[string]interface{}{
		"config": map[string]interface{}{
			"plugins": []interface{}{"auth", "cache"},
		},
	}
	src := map[string]interface{}{
		"config": map[string]interface{}{
			"plugins": []interface{}{"cache", "logging"},
		},
	}
	Merge(dst, src, true, nil, true)
	got := dst["config"].(map[string]interface{})["plugins"]
	want := []interface{}{"auth", "cache", "logging"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nested array merge: got %v, want %v", got, want)
	}
}
