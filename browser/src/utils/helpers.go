package utils

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type C = layout.Context
type D = layout.Dimensions

func FillLayout(gtx C, col color.NRGBA) {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func DrawImage(gtx C, img image.Image) image.Point {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(gtx.Ops)

	scale := float32(math.Min(float64(gtx.Constraints.Max.X)/float64(img.Bounds().Max.X), float64(gtx.Constraints.Max.Y)/float64(img.Bounds().Max.Y)))
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return image.Point{X: int(float32(img.Bounds().Max.X) * scale), Y: int(float32(img.Bounds().Max.Y) * scale)}
}
