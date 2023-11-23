package life_test

import (
	"github.com/k6mil6/game/pkg/life"
	"log"
	"testing"
)

func TestNewWorld(t *testing.T) {
	height := 10
	width := 4

	world, err := life.NewWorld(height, width)
	if err != nil {
		log.Fatal(err)
	}

	if world.Height != height {
		t.Errorf("expected height: %d, actual height: %d", height, world.Height)
	}

	if world.Width != width {
		t.Errorf("expected width: %d, actual width: %d", width, world.Width)
	}

	if len(world.Cells) != height {
		t.Errorf("expected height: %d, actual number of rows: %d", height, len(world.Cells))
	}

	for i, row := range world.Cells {
		if len(row) != width {
			t.Errorf("expected width: %d, actual row's %d len: %d", width, i, world.Width)
		}
	}
}
