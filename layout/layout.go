package layout

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/divan/graphx/graph"
	"gopkg.in/cheggaaa/pb.v1"
)

// stableThreshold determines the movement diff needed to
// call the system stable
const stableThreshold = 2.001

// Layout represents physical layout used to process graph.
type Layout interface {
	Nodes() map[string]*Object
	Calculate()
	CalculateN(n int)

	AddForce(Force)
}

// Layout3D implements Layout interface for force-directed 3D graph.
type Layout3D struct {
	objects map[string]*Object // node ID as a key
	links   []*graph.Link
	forces  []Force
}

// New initializes 3D layout with objects data and set of forces.
func New(g *graph.Graph, forces ...Force) *Layout3D {
	l := &Layout3D{
		objects: make(map[string]*Object),
		links:   g.Links(),
		forces:  forces,
	}

	l.initPositions(g)

	return l
}

// initPositions inits layout graph from the original graph data.
func (l *Layout3D) initPositions(g *graph.Graph) {
	for i, node := range g.Nodes() {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		x := int(radius * math.Cos(rollAngle))
		y := int(radius * math.Sin(rollAngle))
		z := int(radius * math.Sin(yawAngle))

		/* FIXME
		var weight int = 1
		if wnode, ok := node.(graph.WeightedNode); ok {
			weight = wnode.Weight()
		}
		if weight == 0 {
			weight = 1
		}
		*/

		o := NewObject(x, y, z)
		o.ID = node.ID()

		l.objects[node.ID()] = o
	}

	l.resetForces()
}

// Calculate runs positions' recalculations iteratively until the
// system minimizes it's energy.
func (l *Layout3D) Calculate() {
	// tx is the total movement, which should drop to the minimum
	// at the minimal energy state
	fmt.Println("Simulation started...")
	var (
		now    = time.Now()
		count  int
		prevTx float64
	)
	for tx := math.MaxFloat64; math.Abs(tx-prevTx) >= stableThreshold; {
		prevTx = tx
		tx = l.UpdatePositions()
		log.Println("PrevTx, tx:", tx, ", diff:", math.Abs(tx-prevTx))
		count++
		if count%1000 == 0 {
			since := time.Since(now)
			fmt.Printf("Iterations: %d, tx: %f, time: %v\n", count, tx, since)
		}
	}
	fmt.Printf("Simulation finished in %v, run %d iterations\n", time.Since(now), count)
}

// CalculateN run positions' recalculations exactly N times.
func (l *Layout3D) CalculateN(n int) {
	fmt.Println("Simulation started...")
	bar := pb.StartNew(n)
	for i := 0; i < n; i++ {
		l.UpdatePositions()
		bar.Increment()
	}
	bar.FinishPrint("Simulation finished")

}

// UpdatePositions recalculates nodes' positions, applying all the forces.
// It returns average amount of movement generated by this step.
func (l *Layout3D) UpdatePositions() float64 {
	l.resetForces()

	for _, force := range l.forces {
		apply := force.Rule()
		apply(force, l.objects, l.links)
	}

	return l.integrate()
}

func (l *Layout3D) resetForces() {
	// TODO FIXME
}

// AddForce adds force to the internal list of forces.
func (l *Layout3D) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

// Nodes returns nodes information.
func (l *Layout3D) Nodes() map[string]*Object {
	return l.objects
}

// Links returns graph data links.
func (l *Layout3D) Links() []*graph.Link {
	return l.links
}
