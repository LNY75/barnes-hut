package main

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse

	for i := 1; i <= numGens; i++ {
		// update universe and add the new universe to timePoints
		timePoints[i] = UpdateUniverse(timePoints[i-1], time, theta)

	}
	return timePoints
}

// UpdateUniverse updates the current universe after a time t; it returns a pointer to the updated universe
func UpdateUniverse(currentUniverse *Universe, t, theta float64) *Universe {
	newUniverse := CopyUniverse(currentUniverse)

	// construct the quad tree HERE. Doing this anywhere else is going to cause problems. SERIOUS problems. Problems that took hours to resolve.
	qt := BuildQuadTree(newUniverse)
	AssignClusterPos(qt.root)

	for s := range currentUniverse.stars {
		// update pos, vel and accel
		newUniverse.stars[s].Update(qt, currentUniverse, t, theta)

	}
	return newUniverse
}

// need to change this to use current acc/vel/pos
func (s *Star) Update(qt *QuadTree, univ *Universe, t, theta float64) {
	acc := s.NewAccel(qt, univ, theta)
	vel := s.NewVelocity(t)
	pos := s.NewPosition(t)
	s.acceleration, s.velocity, s.position = acc, vel, pos
}

/*
we want to sum all forces acting on Star s
initilize netForce
initialize a queue ([]*Node) of length 1
set the root of the quad tree to be queue[0]

while queue.length != 0:
	current <- queue[0]
	if current is an actual star (has no children) (and very importantly, current is not s itself): compute gravity between current and s and add it to netFroce
	else if current has children:
		calculate s/d for current, compare it with theta (usually 0.5)
		if s/d > theta: we add all of current's children to the queue to be explored later
		otherwise: compute gravitational force between current and s, add that to netForce
	remove the first element from queue (queue  = queue[1:])
*/
// ComputeNetForce sums the forces of all bodies in the universe acting on b.
func ComputeNetForce(qt *QuadTree, star *Star, theta float64) OrderedPair {
	var netForce OrderedPair

	// loop through all nodes of the quad tree in a BFS manner
	var queue []*Node = make([]*Node, 1)
	queue[0] = qt.root

	for len(queue) != 0 {
		current := queue[0]

		// if this is an actual star, not a dummy, we add the force it exerts on star directly to netForce
		if current.children == nil && current.star != star {
			F := ComputeGravityForce(star, current.star)
			netForce.Add(F)
		} else {
			sd := Theta(current, star)
			if sd > theta {
				// if sd > theta, we add its children to the queue to be explored later
				for _, c := range current.children {
					if c != nil {
						queue = append(queue, c)
					}
				}
			} else {
				// treat the collection of stars in this subtree as a single object
				netForce.Add(ComputeGravityForce(star, current.star))
			}
		}
		queue = queue[1:]
	}
	return netForce
}

// Computes theta = s/d = (sector width)/(distance) that determines whether individual gravitatioal effect from a cluster should be considered
// big theta means we should probably consider the individual gravitational effect of stars in the cluster
func Theta(n *Node, star *Star) float64 {
	s := n.sector.width
	d := Dist(n.star, star)
	return s / d
}
