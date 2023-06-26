package frame

import (
	"testing"
)

func TestFrameControl(t *testing.T) {
	fc, err := NewFrameControl("../../public/gopher.png")
	if err != nil {
		t.Error("failed")
	}
	fc.DrawFrame(300)
}
