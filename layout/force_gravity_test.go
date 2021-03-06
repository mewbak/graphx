package layout

import (
	"testing"

	"github.com/divan/graphx/generation/basic"
)

func TestGravity(t *testing.T) {
	t.Run("2_bodies_eachoneach", func(t *testing.T) {
		graph := basic.NewLineGenerator(2).Generate()
		gravity := NewGravityForce(-1.0, EachOnEach)
		l := NewWithForces(graph, gravity)

		l.objects["0"].SetPosition(0, 0, 0)
		l.objects["1"].SetPosition(100, 100, 100)

		for i := 0; i < 10000; i++ {
			l.UpdatePositions()
			checkDistance(t, l.objects["0"], l.objects["1"], 0, 100)
		}
	})
	t.Run("2_bodies_barnehut", func(t *testing.T) {
		graph := basic.NewLineGenerator(2).Generate()
		gravity := NewGravityForce(-1.0, BarneHutMethod)
		l := NewWithForces(graph, gravity)

		l.objects["0"].SetPosition(0, 0, 0)
		l.objects["1"].SetPosition(100, 100, 100)

		for i := 0; i < 1000; i++ {
			l.UpdatePositions()
			checkDistance(t, l.objects["0"], l.objects["1"], 0, 100)
		}
	})
	t.Run("2_bodies_barne_hut_bug", func(t *testing.T) {
		graph := basic.NewLineGenerator(2).Generate()
		gravity := NewGravityForce(-1.0, BarneHutMethod)
		l := NewWithForces(graph, gravity)

		l.objects["0"].SetPosition(-541.38, -541.38, -541.38)
		l.objects["1"].SetPosition(641.38, 641.38, 641.38)

		l.UpdatePositions()
		checkDistance(t, l.objects["0"], l.objects["1"], 0, 100)
	})
	t.Run("2_bodies_drag", func(t *testing.T) {
		graph := basic.NewLineGenerator(2).Generate()
		gravity := NewGravityForce(-1.0, EachOnEach)
		drag := NewDragForce(1, ForEachNode)
		l := NewWithForces(graph, gravity)
		l_drag := NewWithForces(graph, gravity, drag)

		l.objects["0"].SetPosition(0, 0, 0)
		l.objects["1"].SetPosition(10, 10, 10)

		l_drag.objects["0"].SetPosition(0, 0, 0)
		l_drag.objects["1"].SetPosition(10, 10, 10)

		for i := 0; i < 100; i++ {
			l.UpdatePositions()
			l_drag.UpdatePositions()
			checkDistance(t, l.objects["0"], l.objects["1"], 0, 10)
			checkDistance(t, l_drag.objects["0"], l_drag.objects["1"], 0, 10)
		}
		right1 := l.objects["1"]
		right2 := l_drag.objects["1"]

		// X for right2 should be much smaller than right1, because drag force applied and repeated 100 times
		if right1.X()/right2.X() < 2 {
			t.Fatalf("Expect right2 X values be signifantly smaller than right1 X (%.5f < %.5f)", right2.X(), right1.X())
		}
	})
}

// checkDistance checks left and right distances from their initial positions (x0 and x1). they should be equal.
func checkDistance(t *testing.T, left, right *Object, x0, x1 float64) {
	dx1, dx2 := x0-left.X(), right.X()-x1
	dy1, dy2 := x0-left.Y(), right.Y()-x1
	dz1, dz2 := x0-left.Z(), right.Z()-x1
	if dx1-dx2 > 0.0001 {
		t.Fatalf("Expect dX be equal for left and right, but got (%.2f, %.2f)", dx1, dx2)
	}
	if dy1-dy2 > 0.0001 {
		t.Fatalf("Expect dY be equal for left and right, but got (%.2f, %.2f)", dy1, dy2)
	}
	if dz1-dz2 > 0.0001 {
		t.Fatalf("Expect dZ be equal for left and right, but got (%.2f, %.2f)", dz1, dz2)
	}
}
