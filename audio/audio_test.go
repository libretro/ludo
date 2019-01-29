package audio

import (
	"testing"
)

func Test_Sample(t *testing.T) {
	t.Run("Doesn't crash when called", func(t *testing.T) {
		Sample(30000, 30000)
	})
}
