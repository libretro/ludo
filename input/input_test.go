package input

import (
	"testing"
)

func Test_getPressedReleased(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		var old = States{
			{0, 1, 0, 0},
		}
		var new = States{
			{0, 0, 0, 1},
		}
		pressed, released := getPressedReleased(new, old)
		wantpressed := States{{0, 0, 0, 1}}
		wantreleased := States{{0, 1, 0, 0}}
		if pressed != wantpressed {
			t.Errorf("got = %v, want %v", pressed, wantpressed)
		}
		if released != wantreleased {
			t.Errorf("got = %v, want %v", released, wantreleased)
		}
	})
}
