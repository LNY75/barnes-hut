package main

import (
	"fmt"
	"gifhelper"
	"testing"
)

func TestBuildQuadTree(t *testing.T) {
	// u := CreateCustomUniverse()
	// qt := BuildQuadTree(u)
	// qt.Print()

	// j := CreateJupiterSystem()
	// qt2 := BuildQuadTree(j)
	// qt2.Print()
}

func TestAssignClusterMassAndPos(t *testing.T) {
	// u := CreateCustomUniverse()
	// qt := BuildQuadTree(u)
	// AssignClusterPos(qt.root)
	// qt.Print()
}

func TestJupiterSimulation(t *testing.T) {
	// jupiter := CreateJupiterSystem()

	// var numGen int = 100000
	// var time float64 = 1
	// var imgWidth int = 500
	// var outputFilename string = "out.png"
	// var animOutputFile string = "animated-jupiter"
	// var frameRate int = 1000

	// evolution := BarnesHut(jupiter, numGen, time, 0.5)

	// fmt.Println("evolution is complete")

	// // write out an animation of the universe
	// frames := AnimateSystem(evolution, imgWidth, frameRate, 2)
	// gifhelper.ImagesToGIF(frames, animOutputFile)

	// // export the final frame as an png
	// img := evolution[len(evolution)-1].DrawToCanvas(imgWidth, 2)
	// var c Canvas = CreateNewCanvas(imgWidth, imgWidth)
	// c.img = img
	// c.SaveToPNG(outputFilename)
}

func TestGalaxySimulation(t *testing.T) {
	// g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
	// width := 1.0e23
	// galaxies := []Galaxy{g0}
	// initialUniverse := InitializeUniverse(galaxies, width)
	// // now evolve the universe: feel free to adjust the following parameters.
	// numGens := 50000
	// time := 2e14
	// theta := 0.5

	// timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	// fmt.Println("Simulation run. Now drawing images.")
	// canvasWidth := 900
	// frequency := 1000
	// scalingFactor := 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
	// imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	// fmt.Println("Images drawn. Now generating GIF.")
	// gifhelper.ImagesToGIF(imageList, "galaxy")
	// fmt.Println("GIF drawn.")
}

func TestCollisionSimulation(t *testing.T) {
	// the following sample parameters may be helpful for the "collide" command
	// all units are in SI (meters, kg, etc.)
	// but feel free to change the positions of the galaxies.

	g0 := InitializeGalaxy(500, 4e21, 4e22, 3e22)
	g1 := InitializeGalaxy(500, 4e21, 3e22, 3e22)

	// you probably want to apply a "push" function at this point to these galaxies to move
	// them toward each other to collide.
	// be careful: if you push them too fast, they'll just fly through each other.
	// too slow and the black holes at the center collide and hilarity ensues.

	push(&g0, OrderedPair{-100, 200})
	push(&g1, OrderedPair{200, -100})

	width := 1.0e23

	galaxies := []Galaxy{g0, g1}

	initialUniverse := InitializeUniverse(galaxies, width)

	// now evolve the universe: feel free to adjust the following parameters.
	// There is a maximum number of steps we can simulate the collision. If we reach a point where two black holes have collided with eacher, their position will be the same, and the quadtree will be stuck in an infinite loop trying to allocate them into the right places. The tree will grow indefinitely because no matter how we divide the quadrants, these two black holes will always end up in the same quadrant.
	numGens := 10000
	time := 3e15
	theta := 0.5

	timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	fmt.Println("Simulation run. Now drawing images.")
	canvasWidth := 1000
	frequency := 1000
	scalingFactor := 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, "collision")
	fmt.Println("GIF drawn.")
}

func CreateCustomUniverse() *Universe {
	var A, B, C, D, E, F, G Star
	A.position.x, A.position.y = 1, 13
	B.position.x, B.position.y = 2, 12
	C.position.x, C.position.y = 6, 8
	D.position.x, D.position.y = 9, 9
	E.position.x, E.position.y = 12, 6
	F.position.x, F.position.y = 2, 5
	G.position.x, G.position.y = 8, 2

	A.mass = 1
	B.mass = 2
	C.mass = 3
	D.mass = 4
	E.mass = 5
	F.mass = 6
	G.mass = 7

	var customUniverse Universe
	customUniverse.width = 14
	customUniverse.AddStar(A)
	customUniverse.AddStar(B)
	customUniverse.AddStar(C)
	customUniverse.AddStar(D)
	customUniverse.AddStar(E)
	customUniverse.AddStar(F)
	customUniverse.AddStar(G)

	return &customUniverse
}
