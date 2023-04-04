package ships

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Add - Example
func ExamplePoint_Add() {
	point := Point{X: 1, Y: 1}.Add(Point{X: 2, Y: 2})
	fmt.Println(point.X, point.Y)
	// Output: 3 3
}

func ExamplePoint_Add_negative() {
	point := Point{X: 4, Y: 4}.Add(Point{X: -2, Y: -2})
	fmt.Println(point.X, point.Y)
	// Output: 2 2
}

// MoveTo - 4 test cases
func TestShip_MoveTo(t *testing.T) {
	tests := []struct {
		name   string
		ship   Ship
		moveTo Point
		want   Ship
	}{
		{"X axis", []Point{{0, 0}, {1, 0}}, Point{1, 0}, []Point{{1, 0}, {2, 0}}},
		{"Y axis", []Point{{0, 0}, {0, 1}}, Point{0, 1}, []Point{{0, 1}, {0, 2}}},
		{"XnY axis", []Point{{0, 0}, {1, 0}}, Point{1, 1}, []Point{{1, 1}, {2, 1}}},
		{"empty point", []Point{{0, 0}, {1, 0}}, Point{}, []Point{{0, 0}, {1, 0}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			movedShip := test.ship.MoveTo(test.moveTo)

			assert.Equal(t, test.want, movedShip)
		})
	}
}
