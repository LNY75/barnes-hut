package main

import (
	"fmt"
	"gifhelper"
	"math"
	"os"
)

func CreateJupiterSystem() *Universe {
	// manually generate Jupiter:
	// declaring objects

	var jupiter, io, europa, ganymede, callisto Star

	jupiter.red, jupiter.green, jupiter.blue = 223, 227, 202
	io.red, io.green, io.blue = 249, 249, 165
	europa.red, europa.green, europa.blue = 132, 83, 52
	ganymede.red, ganymede.green, ganymede.blue = 76, 0, 153
	callisto.red, callisto.green, callisto.blue = 0, 153, 76

	jupiter.mass = 1.898 * math.Pow(10, 27)
	io.mass = 8.9319 * math.Pow(10, 22)
	europa.mass = 4.7998 * math.Pow(10, 22)
	ganymede.mass = 1.4819 * math.Pow(10, 23)
	callisto.mass = 1.0759 * math.Pow(10, 23)

	jupiter.radius = 71000000
	io.radius = 1821000
	europa.radius = 1569000
	ganymede.radius = 2631000
	callisto.radius = 2410000

	jupiter.position.x, jupiter.position.y = 2000000000, 2000000000
	io.position.x, io.position.y = 2000000000-421600000, 2000000000
	europa.position.x, europa.position.y = 2000000000, 2000000000+670900000
	ganymede.position.x, ganymede.position.y = 2000000000+1070400000, 2000000000
	callisto.position.x, callisto.position.y = 2000000000, 2000000000-1882700000

	jupiter.velocity.x, jupiter.velocity.y = 0, 0
	io.velocity.x, io.velocity.y = 0, -17320
	europa.velocity.x, europa.velocity.y = -13740, 0
	ganymede.velocity.x, ganymede.velocity.y = 0, 10870
	callisto.velocity.x, callisto.velocity.y = 8200, 0

	// declaring universe and setting its fields.
	var jupiterSystem Universe
	jupiterSystem.width = 4000000000
	jupiterSystem.AddStar(jupiter)
	jupiterSystem.AddStar(io)
	jupiterSystem.AddStar(europa)
	jupiterSystem.AddStar(ganymede)
	jupiterSystem.AddStar(callisto)

	return &jupiterSystem
}

func JupiterSimulation() {
	jupiter := CreateJupiterSystem()

	var numGen int = 100000
	var time float64 = 1
	var imgWidth int = 500
	var outputFilename string = "out.png"
	var animOutputFile string = "animated-jupiter"
	var frameRate int = 1000

	evolution := BarnesHut(jupiter, numGen, time, 0.5)

	fmt.Println("evolution is complete")

	// write out an animation of the universe
	frames := AnimateSystem(evolution, imgWidth, frameRate, 2)
	gifhelper.ImagesToGIF(frames, animOutputFile)

	// export the final frame as an png
	img := evolution[len(evolution)-1].DrawToCanvas(imgWidth, 2)
	var c Canvas = CreateNewCanvas(imgWidth, imgWidth)
	c.img = img
	c.SaveToPNG(outputFilename)
}

func GalaxySimulation() {
	g0 := InitializeGalaxy(500, 4e21, 7e22, 2e22)
	width := 1.0e23
	galaxies := []Galaxy{g0}
	initialUniverse := InitializeUniverse(galaxies, width)
	// now evolve the universe: feel free to adjust the following parameters.
	numGens := 50000
	time := 2e14
	theta := 0.5

	timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	fmt.Println("Simulation run. Now drawing images.")
	canvasWidth := 900
	frequency := 1000
	scalingFactor := 1e11 // a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, "galaxy")
	fmt.Println("GIF drawn.")
}

func CollisionSimulation() {
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

func main() {
	cmd := os.Args[1]
	if cmd == "jupiter" {
		JupiterSimulation()
	} else if cmd == "galaxy" {
		GalaxySimulation()
	} else if cmd == "collision" {
		CollisionSimulation()
	} else {
		panic("hey, please choose command form jupiter/galaxy/collision (case-sensitive).")
	}
}
