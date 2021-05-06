package patch

import (
	"testing"
)

func Test_applyUPS(t *testing.T) {
	t.Run("Can detect a short patch", func(t *testing.T) {
		source := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		patch := []byte{'U', 'P', 'S', '1', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		got, err := applyUPS(patch, source)
		if err.Error() != "patch too small" {
			t.Errorf("applyUPS() = %v, want %v", err, "patch too small")
		}
		if got != nil {
			t.Errorf("applyUPS() = %v, want %v", got, nil)
		}
	})

	t.Run("Can detect a patch with a wrong header", func(t *testing.T) {
		source := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		patch := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
		got, err := applyUPS(patch, source)
		if err.Error() != "invalid patch header" {
			t.Errorf("applyUPS() = %v, want %v", err, "invalid patch header")
		}
		if got != nil {
			t.Errorf("applyUPS() = %v, want %v", got, nil)
		}
	})

	// t.Run("Can apply a valid UPS patch", func(t *testing.T) {
	// 	source := make([]byte, 16512)
	// 	patch := []byte{'U', 'P', 'S', '1',
	// 		0x00, 0x00, 0x80,
	// 		0x00, 0x00, 0x80,
	// 		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	// 	want := make([]byte, 16512)
	// 	got, err := applyUPS(patch, source)
	// 	fmt.Println(err)
	// 	if !reflect.DeepEqual(got, want) {
	// 		t.Errorf("applyUPS() = %v, want %v", got, want)
	// 	}
	// })
}
