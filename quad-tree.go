/*
	stores all quad-tree related operations
*/

package main

import (
	"fmt"
)

//BuildQuadTree builds a quad tree for a universe
func BuildQuadTree(u *Universe) *QuadTree {
	var qt QuadTree = QuadTree{root: InitRoot(u)}

	// start building the quad tree:
	for i := 0; i < len((*u).stars); i++ {
		parent := qt.root
		star := u.stars[i]
		q := star.whichSubQuad(parent.sector)
		next := parent.children[q]

		// find	where to insert the new star into the quad tree
		for next != nil {
			// is the occupied child a dummy or actual star?
			if next.children != nil {
				// the occupied child is a dummy
				parent = next
			} else {
				// the occupied child is an actual star
				// we need another dummy star as the new parent
				var newDummy *Node = parent.NewDummy(q)
				parent.children[q] = newDummy
				// put the previous star back under the dummy
				nq := next.star.whichSubQuad(newDummy.sector)
				newDummy.children[nq] = next
				parent = newDummy
			}
			q = star.whichSubQuad(parent.sector)
			next = parent.children[q]
		}

		//we've found the parent of the star; insert star under that parent:
		parent.children[q] = &Node{
			children: nil,
			star:     star,
			sector:   parent.sector.findNewQuad(q),
		}

	}
	return &qt
}

// prints the content of the quad tree in a BFS manner. For debugging purpose only
func (qt *QuadTree) Print() {
	var queue []*Node = make([]*Node, 1)
	queue[0] = qt.root

	for len(queue) != 0 {
		current := queue[0]
		fmt.Println("(", current.star.position.x, ", ", current.star.position.y, ") ", "mass: ", current.star.mass)
		queue = queue[1:] // omg this works, and it doesn't throw outofbound error!

		for _, c := range current.children {
			if c != nil {
				queue = append(queue, c)
			}
		}
	}
}

// // After constructing the quad tree, we need to assign mass to each dummy star
// // assignment of masses is performed in a recursive manner
// // **warning: always get set all positions first before calculating their masses
// func AssignClutserMass(n *Node) float64 {
// 	if n.children != nil {
// 		n.star.mass = 0
// 		for _, c := range n.children {
// 			if c != nil && c.children == nil {
// 				n.star.mass += c.star.mass
// 			} else if c != nil {
// 				n.star.mass += AssignClutserMass(c)
// 			}
// 		}
// 	}
// 	return n.star.mass
// }

// this data structure is used for assigning positions of internal nodes of a quatree only; because I need to pass on the x,y coordinates AND the combined mass for the recursion
type PseudoStar struct {
	x, y, mass float64
}

// After constructing the quad tree, we need to update the mass and position of each dummy based on the center of mass of its children
func AssignClusterPos(n *Node) PseudoStar {
	if n.children != nil {
		// get a list of non-nil children:
		children := n.GetRealChildren()
		n.star.position.x = 0
		n.star.position.y = 0
		n.star.mass = 0

		for i := 0; i < len(children); i++ {
			if children[i].children == nil {
				n.star.position = CenterOfMass(n.star, children[i].star)
				n.star.mass += children[i].star.mass
			} else {
				ps := AssignClusterPos(children[i])
				n.star.position = CenterOfMass2(n.star, &ps)
				n.star.mass += ps.mass
			}
		}
	}
	return PseudoStar{
		x:    n.star.position.x,
		y:    n.star.position.y,
		mass: n.star.mass,
	}
}

// CenterOfMass computes the position of the center of mass of s1 and s2
func CenterOfMass(n *Star, m *Star) OrderedPair {
	nx := n.position.x
	ny := n.position.y
	mx := m.position.x
	my := m.position.y
	nMass := n.mass
	mMass := m.mass

	return OrderedPair{
		(nx*nMass + mx*mMass) / (nMass + mMass),
		(ny*nMass + my*mMass) / (nMass + mMass),
	}
}

// Calculate center of mass between a star and a psedostar
func CenterOfMass2(n *Star, m *PseudoStar) OrderedPair {
	nx := n.position.x
	ny := n.position.y
	mx := m.x
	my := m.y
	nMass := n.mass
	mMass := m.mass

	return OrderedPair{
		(nx*nMass + mx*mMass) / (nMass + mMass),
		(ny*nMass + my*mMass) / (nMass + mMass),
	}

}

// returns a list of all non-nil children of node n
func (n *Node) GetRealChildren() []*Node {
	var children []*Node
	for _, c := range n.children {
		if c != nil {
			children = append(children, c)
		}
	}
	return children
}

// // check whether nodes n exists in nodes. helper funciton for debugging only
// func (n *Node) in(nodes []*Node) bool {
// 	for _, e := range nodes {
// 		if e == n {
// 			return true
// 		}
// 	}
// 	return false
// }

//InitRoot initializes a root node for a quad tree, given the size of the universe:
func InitRoot(u *Universe) *Node {
	// define the outermost quadrant:
	var q Quadrant = Quadrant{
		x:     0,
		y:     0,
		width: u.width,
	}

	var s Star = Star{
		position: OrderedPair{
			x: u.width / 2,
			y: u.width / 2,
		},
	}

	var root Node = Node{
		children: make([]*Node, 4),
		star:     &s,
		sector:   q,
	}
	return &root
}

//whichSubQuad determines which quadrant(NW, NE, SW, SE) the star belongs to
//returns a number correponding to each sub-quadrant:
/*
	0 -> NW
	1 -> NE
	2 -> SW
	3 -> SE
*/
func (s *Star) whichSubQuad(q Quadrant) int {
	x := s.position.x
	y := s.position.y
	center_x := q.x + q.width/2
	center_y := q.y + q.width/2
	if y > center_y {
		if x < center_x {
			return 0
		} else {
			return 1
		}
	} else {
		if x < center_x {
			return 2
		} else {
			return 3
		}
	}
}

//isInUniverse determines whether the star is still in the universe
func (s *Star) isInUniverse(u *Universe) bool {
	x := s.position.x
	y := s.position.y
	if x >= 0 && x <= u.width && y >= 0 && y <= u.width {
		return true
	}
	return false
}

//creates a new quadrant (NW, NE, SW, SE) from an existing quadrant
func (q *Quadrant) findNewQuad(i int) Quadrant {
	var x float64
	var y float64
	switch i {
	case 0:
		// create quadrant from NW of the original quadrant
		x = q.x
		y = q.y + q.width/2
	case 1:
		// create quadrant from NE of the original quadrant
		x = q.x + q.width/2
		y = q.y + q.width/2
	case 2:
		// create quadrant from SW of the original quadrant
		x = q.x
		y = q.y
	case 3:
		// create quadrant from SE of the original quadrant
		x = q.x + q.width/2
		y = q.y
	}
	return Quadrant{
		x:     x,
		y:     y,
		width: q.width / 2,
	}
}

// creates a new dummy node from an existing node from one of its quadrants (NW/NE/SW/SE, represnted by an int)
func (n *Node) NewDummy(i int) *Node {
	x := n.star.position.x
	y := n.star.position.y
	w := n.sector.width / 4

	switch i {
	case 0:
		x -= w
		y += w
	case 1:
		x += w
		y += w
	case 2:
		x -= w
		y -= w
	case 3:
		x += w
		y -= w
	}
	star := Star{
		position: OrderedPair{x: x, y: y},
	}
	sect := Quadrant{
		x:     star.position.x - w,
		y:     star.position.y - w,
		width: 2 * w,
	}
	return &Node{
		children: make([]*Node, 4),
		star:     &star,
		sector:   sect,
	}
}
