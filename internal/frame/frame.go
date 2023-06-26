package frame

import (
	"image"

	"github.com/fogleman/gg"
)

type FrameControl struct {
	asset image.Image
	W     int
	H     int
}

func NewFrameControl(path string) (*FrameControl, error) {
	img, err := gg.LoadImage(path)
	if err != nil {
		return nil, err
	}
	return &FrameControl{
		asset: img,
		W:     400,
		H:     400,
	}, nil
}

func (fc *FrameControl) DrawFrame(degree int) image.Image {
	iw, ih := fc.asset.Bounds().Dx(), fc.asset.Bounds().Dy()
	dc := gg.NewContext(fc.W, fc.H)
	dc.SetHexColor("#fff")
	dc.Clear()
	// draw outline
	dc.SetHexColor("#ff0000")
	dc.SetLineWidth(1)
	dc.DrawRectangle(0, 0, float64(fc.W), float64(fc.H))
	dc.Stroke()
	// draw image with current matrix applied
	dc.SetHexColor("#0000ff")
	dc.SetLineWidth(2)
	dc.Rotate(gg.Radians(float64(degree)))
	dc.DrawRectangle(100, 0, float64(iw), float64(ih))
	dc.StrokePreserve()
	dc.Clip()
	dc.DrawImageAnchored(fc.asset, 100, 0, 0.0, 0.0)
	// dc.SavePNG("./out.png")
	return dc.Image()
}
