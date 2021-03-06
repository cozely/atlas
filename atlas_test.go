package atlas_test

import (
	"image"
	"reflect"
	"testing"

	"github.com/cozely/atlas"
)

func Test1(t *testing.T) {
	items := []image.Image{
		image.Rect(0, 0, 1, 1),
		image.Rect(0, 0, 2, 1),
		image.Rect(0, 0, 1, 3),
		image.Rect(0, 0, 4, 1),
	}
	a := atlas.New(image.Pt(4, 4))
	mappings, err := a.Pack(items)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	out := []atlas.Mapping{
		{Bin:0, Bounds:image.Rect(3, 1, 4, 2)},
		{Bin:0, Bounds:image.Rect(1, 1, 3, 2)},
		{Bin:0, Bounds:image.Rect(0, 1, 1, 4)},
		{Bin:0, Bounds:image.Rect(0, 0, 4, 1)},
	}

	for i := range items {
		if items[i].Bounds().Size() != mappings[i].Bounds.Size() {
			t.Errorf("wrong mapping size for item %d: %d, %d", i, mappings[i].Bounds.Dx(), mappings[i].Bounds.Dy())
		}
	}

	if !reflect.DeepEqual(mappings, out) {
		t.Errorf("wrong mappings: %#v", mappings)
	}
}

func Test2(t *testing.T) {
	items := []image.Image{
		image.Rect(0, 0, 1, 1),
		image.Rect(0, 0, 2, 1),
		image.Rect(0, 0, 1, 3),
		image.Rect(0, 0, 5, 1),
		image.Rect(0, 0, 4, 1),
	}
	a := atlas.New(image.Pt(4, 4))
	mappings, err := a.Pack(items)
	if err == nil {
		t.Errorf("should not have succeded!")
	}

	out := []atlas.Mapping{
		{Bin:0, Bounds:image.Rect(0, 0, 0, 0)},
		{Bin:0, Bounds:image.Rect(0, 0, 0, 0)},
		{Bin:0, Bounds:image.Rect(0, 0, 0, 0)},
		{Bin:0, Bounds:image.Rect(0, 0, 0, 0)},
		{Bin:0, Bounds:image.Rect(0, 0, 0, 0)},
	}

	if !reflect.DeepEqual(mappings, out) {
		t.Errorf("wrong mappings: %#v", mappings)
	}
}

func Test3(t *testing.T) {
	items := []image.Image{
		image.Rect(0, 0, 1, 1),
		image.Rect(0, 0, 2, 1),
		image.Rect(0, 0, 1, 3),
		image.Rect(0, 0, 4, 1),
	}
	a := atlas.New(image.Pt(5, 5))
	mappings, err := a.Pack(items)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	out := []atlas.Mapping{
		atlas.Mapping{Bin:0, Bounds:image.Rectangle{Min:image.Point{X:4, Y:0}, Max:image.Point{X:5, Y:1}}},
		atlas.Mapping{Bin:0, Bounds:image.Rectangle{Min:image.Point{X:1, Y:1}, Max:image.Point{X:3, Y:2}}},
		atlas.Mapping{Bin:0, Bounds:image.Rectangle{Min:image.Point{X:0, Y:1}, Max:image.Point{X:1, Y:4}}},
		atlas.Mapping{Bin:0, Bounds:image.Rectangle{Min:image.Point{X:0, Y:0}, Max:image.Point{X:4, Y:1}}},
	}

	for i := range items {
		if items[i].Bounds().Size() != mappings[i].Bounds.Size() {
			t.Errorf("wrong mapping size for item %d: %d, %d", i, mappings[i].Bounds.Dx(), mappings[i].Bounds.Dy())
		}
	}

	if !reflect.DeepEqual(mappings, out) {
		t.Errorf("wrong mappings: %#v", mappings)
	}
}
