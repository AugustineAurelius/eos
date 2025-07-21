package generics_test

import (
	"reflect"
	"testing"

	"github.com/AugustineAurelius/eos/pkg/generics"
)

func TestSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	genericsSlice := generics.FromSlice(slice)
	if !reflect.DeepEqual(genericsSlice, generics.SliceOps[int]{1, 2, 3, 4, 5}) {
		t.Fatal("FromSlice failed")
	}

	genericsSlice = genericsSlice.FilterFunc(func(i int) bool {
		return i > 3
	})

	if !reflect.DeepEqual(genericsSlice, generics.SliceOps[int]{4, 5}) {
		t.Fatal("FilterFunc failed")
	}

	founded, ok := genericsSlice.FindFunc(func(i int) bool {
		return i == 4
	})

	if !ok {
		t.Fatal("FindFunc failed")
	}
	if founded != 4 {
		t.Fatal("FindFunc failed")
	}

}

func TestExtract(t *testing.T) {
	slice := []struct {
		ID   int
		Name string
	}{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
		{ID: 3, Name: "Jim"},
		{ID: 4, Name: "Jill"},
		{ID: 5, Name: "Jack"},
	}

	genericsSlice := generics.FromSlice(slice)
	extracted := generics.Extract(genericsSlice, func(i struct {
		ID   int
		Name string
	}) string {
		return i.Name
	})

	if !reflect.DeepEqual(extracted, generics.SliceOps[string]{"John", "Jane", "Jim", "Jill", "Jack"}) {
		t.Fatal("Extract failed")
	}
}
