package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"ic/shared"
)

type C = layout.Context
type D = layout.Dimensions

var testImage image.Image
var testList widget.List

func main() {
	testImage, _ = shared.LoadImage("../../testAssets/screenA.png")

	testList = widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Image-compare browser"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func draw(w *app.Window) error {
	var ops op.Ops

	th := material.NewTheme()
	th.Palette = material.Palette{
		Bg:         color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xFF},
		Fg:         color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF},
		ContrastBg: color.NRGBA{R: 0x1E, G: 0x1E, B: 0x1E, A: 0xFF},
		ContrastFg: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layoutSidebar(gtx, th)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layoutMainContent(gtx, th)
				}),
			)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func layoutSidebar(gtx C, th *material.Theme) D {
	gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			fillLayout(gtx, th.Palette.ContrastBg)
			ls := material.List(th, &testList)

			return ls.Layout(gtx, 10, func(gtx C, index int) D {
				lbl := material.Body1(th, "item")
				return lbl.Layout(gtx)
			})
		}),
	)
}

func layoutMainContent(gtx C, th *material.Theme) D {
	fillLayout(gtx, th.Palette.Bg)
	drawImage(gtx, testImage)
	return D{}
}

func fillLayout(gtx C, col color.NRGBA) {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawImage(gtx C, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(gtx.Ops)

	scale := float32(gtx.Constraints.Max.X) / float32(img.Bounds().Max.X)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
