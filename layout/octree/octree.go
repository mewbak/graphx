package octree

import (
	"errors"
	"fmt"
)

// Octree represents Octree data structure.
// See https://en.wikipedia.org/wiki/Octree for details.
type Octree struct {
	Root Octant
	ids  map[string]Octant // for fast lookup in findLeaf
}

// octant represent a node in octree, which is an octant of a cube.
// See: http://en.wikipedia.org/wiki/Octant_(solid_geometry)
type Octant interface {
	Center() Point
	Insert(p Point) Octant
}

// New inits new octree.
func New() *Octree {
	return &Octree{
		ids: make(map[string]Octant),
	}
}

// Insert adds new Point into the Octree data structure.
func (o *Octree) Insert(p Point) {
	if o.Root == nil {
		o.Root = o.NewLeaf(p)
		return
	}

	o.Root = o.Root.Insert(p)
}

// FindLeafs searches for the leaf with the given id.
func (o *Octree) FindLeaf(id string) (Octant, error) {
	oct, ok := o.ids[id]
	if !ok {
		return nil, errors.New("leaf not found")
	}
	return oct, nil
}

// String implements Stringer interface for octree.
func (o *Octree) String() string {
	return fmt.Sprintf("Root: %T, leafs: %v", o.Root, o.Root.(*Node).Leafs)
}
