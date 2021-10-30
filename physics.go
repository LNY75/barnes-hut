package main

import (
	"math"
)

//Performs addition on 2 2-d vectors
func (v *OrderedPair) Add(v2 OrderedPair) {
	v.x += v2.x
	v.y += v2.y
}

// NewVelocity makes the velocity of this object consistent with the acceleration.
func (s *Star) NewVelocity(t float64) OrderedPair {
	return OrderedPair{
		x: s.velocity.x + s.acceleration.x*t,
		y: s.velocity.y + s.acceleration.y*t,
	}
}

// NewPosition computes the new poosition given the updated acc and velocity.
func (s *Star) NewPosition(t float64) OrderedPair {
	return OrderedPair{
		x: s.position.x + s.velocity.x*t + 0.5*s.acceleration.x*t*t,
		y: s.position.y + s.velocity.y*t + 0.5*s.acceleration.y*t*t,
	}
}

// UpdateAccel computes the new accerlation vector for b
func (s *Star) NewAccel(qt *QuadTree, univ *Universe, theta float64) OrderedPair {
	F := ComputeNetForce(qt, s, theta)
	return OrderedPair{
		x: F.x / s.mass,
		y: F.y / s.mass,
	}
}

// ComputeGravityForce computes the gravity force between star 1 and star 2.
func ComputeGravityForce(s1, s2 *Star) OrderedPair {
	d := Dist(s1, s2)
	deltaX := s2.position.x - s1.position.x
	deltaY := s2.position.y - s1.position.y
	F := G * s1.mass * s2.mass / (d * d)

	r := OrderedPair{
		x: F * deltaX / d,
		y: F * deltaY / d,
	}
	return r
}

// Compute the Euclidian Distance between two stars
func Dist(s1, s2 *Star) float64 {
	dx := s1.position.x - s2.position.x
	dy := s1.position.y - s2.position.y
	d := math.Sqrt(dx*dx + dy*dy)
	return d
}
