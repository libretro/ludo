package main

import "testing"

func TestAutoSwapInterval(t *testing.T) {
	tests := []struct {
		name      string
		refreshHz float64
		coreFPS   float64
		want      int
	}{
		{name: "60Hz stays at swap 1", refreshHz: 60, coreFPS: 59.73, want: 1},
		{name: "120Hz prefers swap 2 for 60fps", refreshHz: 120, coreFPS: 60, want: 2},
		{name: "144Hz prefers swap 3 for 48fps", refreshHz: 144, coreFPS: 48, want: 3},
		{name: "144Hz keeps swap 1 for 60fps", refreshHz: 144, coreFPS: 60, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := autoSwapInterval(tt.refreshHz, tt.coreFPS); got != tt.want {
				t.Fatalf("autoSwapInterval(%v, %v) = %v, want %v", tt.refreshHz, tt.coreFPS, got, tt.want)
			}
		})
	}
}

func TestBlocksOnSwap(t *testing.T) {
	if got := blocksOnSwap(60.0, 59.73); !got {
		t.Fatalf("blocksOnSwap() = %v, want %v", got, true)
	}

	if got := blocksOnSwap(144.0, 60.0); got {
		t.Fatalf("blocksOnSwap() = %v, want %v", got, false)
	}
}
