// Copyright (c) 2017-2018 Laurent Moussault. All rights reserved.
// Licensed under a simplified BSD license (see LICENSE file).

package atlas

import (
	"fmt"
	"image"
)

////////////////////////////////////////////////////////////////////////////////

type region struct {
	first, second *region
	bounds        image.Rectangle
	content       int
}

////////////////////////////////////////////////////////////////////////////////

const (
	something = -1
	nothing   = -2
)

////////////////////////////////////////////////////////////////////////////////

func (n *region) String() string {
	if n.first == nil {
		switch n.content {
		case nothing:
			return "nothing"
		case something:
			return "something"
		default:
			return fmt.Sprintf("%d", n.content)
		}
	}
	return "{ " + n.first.String() + ", " + n.second.String() + " }"
}

////////////////////////////////////////////////////////////////////////////////

func (n *region) insert(items []image.Image, index int) *region {
	// If already split, recurse

	if n.first != nil {
		f := n.first.insert(items, index)
		if f != nil {
			return f
		}

		return n.second.insert(items, index)
	}

	// It's a leaf

	if n.content != nothing {
		// Already filled
		return nil
	}

	s := items[index].Bounds().Size()

	if n.bounds.Dx() < s.X || n.bounds.Dy() < s.Y {
		// Too small
		return nil
	}

	if n.bounds.Size() == s {
		// It's a match!
		n.content = index
		return n
	}

	// Split the leaf

	n.first = &region{content: nothing}
	n.second = &region{content: nothing}

	if n.bounds.Dx()-s.X > n.bounds.Dy()-s.Y {
		n.first.bounds = image.Rectangle {
			Min: n.bounds.Min,
			Max: image.Pt(n.bounds.Min.X+s.X, n.bounds.Max.Y),
		}
		// n.first.x, n.first.y = n.x, n.y
		// n.first.w, n.first.h = w, n.h

		n.second.bounds = image.Rectangle {
			Min: image.Pt(n.bounds.Min.X+s.X, n.bounds.Min.Y),
			Max: image.Pt(n.bounds.Max.X, n.bounds.Max.Y),
		}
		// n.second.x, n.second.y = n.x+w, n.y
		// n.second.w, n.second.h = n.w-w, n.h

	} else {
		n.first.bounds = image.Rectangle {
			Min: n.bounds.Min,
			Max: image.Pt(n.bounds.Max.X, n.bounds.Min.Y+s.Y),
		}
		// n.first.x, n.first.y = n.x, n.y
		// n.first.w, n.first.h = n.w, h

		n.second.bounds = image.Rectangle {
			Min: image.Pt(n.bounds.Min.X, n.bounds.Min.Y+s.Y),
			Max: image.Pt(n.bounds.Max.X, n.bounds.Max.Y),
		}
		// n.second.x, n.second.y = n.x, n.y+h
		// n.second.w, n.second.h = n.w, n.h-h

	}

	return n.first.insert(items, index)
}

////////////////////////////////////////////////////////////////////////////////
