package input

import (
	"testing"
)

func Test_reset(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		var state1 inputstate
		var state2 inputstate
		state1[0][0] = true
		state1[3][5] = true
		state1 = reset(state1)
		if state1 != state2 {
			t.Errorf("got = %v, want %v", state1, state2)
		}
	})
}

func Test_getPressedReleased(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		var old = inputstate{
			{false, true, false, false},
		}
		var new = inputstate{
			{false, false, false, true},
		}
		pressed, released := getPressedReleased(new, old)
		wantpressed := inputstate{{false, false, false, true}}
		wantreleased := inputstate{{false, true, false, false}}
		if pressed != wantpressed {
			t.Errorf("got = %v, want %v", pressed, wantpressed)
		}
		if released != wantreleased {
			t.Errorf("got = %v, want %v", released, wantreleased)
		}
	})
}
