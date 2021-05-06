package input

import (
	"testing"
)

func Test_getPressedReleased(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		var old = States{
			{false, true, false, false},
		}
		var new = States{
			{false, false, false, true},
		}
		pressed, released := getPressedReleased(new, old)
		wantpressed := States{{false, false, false, true}}
		wantreleased := States{{false, true, false, false}}
		if pressed != wantpressed {
			t.Errorf("got = %v, want %v", pressed, wantpressed)
		}
		if released != wantreleased {
			t.Errorf("got = %v, want %v", released, wantreleased)
		}
	})
}
