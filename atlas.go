// Copyright (c) 2017-2018 Laurent Moussault. All rights reserved.
// Licensed under a simplified BSD license (see LICENSE file).

package atlas

import (
	"fmt"
	"image"
	"sort"
)

/*
The current implementation is rather naive. The objective is to later improve it
with ideas from the PhD Thesis of Andrea Lodi, "Algorithms for Two Dimensional
Bin Packing and Assignment Problems":

  http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.98.3502&rep=rep1&type=pdf

TODO:

- add flipped items, to be able to make all items horizontal
- find a way to compute the "touching perimeter" score
- pre-allocate the bins
*/

////////////////////////////////////////////////////////////////////////////////

// An Atlas contains the mapping information to pack a set of images (called
// items) into an array of bigger images (called bins).
type Atlas struct {
	size  image.Point
	trees []region
	ideal int
}

// A Mapping is the location of an item inside an atlas.
type Mapping struct {
	Bin    int
	Bounds image.Rectangle
}

////////////////////////////////////////////////////////////////////////////////

// New returns a new atlas.
func New(size image.Point) *Atlas {
	return &Atlas{
		size: size,
	}
}

////////////////////////////////////////////////////////////////////////////////

// BinCount returns the number of bins currently in the atlas.
func (a *Atlas) BinCount() int16 {
	return int16(len(a.trees))
}

// Unused returns the number of unused pixels (i.e. not allocated to any image)
// in the atlas.
func (a *Atlas) Unused() int {
	return len(a.trees)*a.size.X*a.size.Y - a.ideal
}

////////////////////////////////////////////////////////////////////////////////

// Pack fits all the items into the atlas.
func (a *Atlas) Pack(items []image.Image) ([]Mapping, error) {
	indices := make([]int, len(items), len(items))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		si := items[indices[i]].Bounds().Size()
		sj := items[indices[j]].Bounds().Size()
		return si.X*2+si.Y*2 > sj.X*2+sj.Y*2
	})

	mappings := make([]Mapping, 0, len(items))
	for i := range indices {
		s := items[indices[i]].Bounds().Size()
		a.ideal += s.X * s.Y
		var reg *region
		for j := range a.trees {
			reg = a.trees[j].insert(items, indices[i])
			if reg != nil {
				mappings = append(mappings, Mapping{
					Bin:    j,
					Bounds: reg.bounds,
				})
				break
			}
		}
		if reg == nil {
			a.trees = append(
				a.trees,
				region{
					bounds:  image.Rectangle{Max: a.size},
					content: nothing,
				},
			)
			reg = a.trees[len(a.trees)-1].insert(items, indices[i])
			if reg != nil {
				mappings = append(mappings, Mapping{
					Bin:    len(a.trees) - 1,
					Bounds: reg.bounds,
				})
			} else {
				return mappings, fmt.Errorf(
					"atlas.Atlas.Pack: unable to fit item %d",
					indices[i],
				)
			}
		}
	}

	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Bin < mappings[j].Bin
	})

	return mappings, nil
}

////////////////////////////////////////////////////////////////////////////////
