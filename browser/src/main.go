package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"ic/browser/src/utils"
	"ic/shared"
)

type C = layout.Context
type D = layout.Dimensions

var testImage image.Image
var testList widget.List

var imageButton1 widget.Clickable
var imageButton2 widget.Clickable

var comparisonList widget.List
var comparisons []shared.Comparison
var comparisonButtons []widget.Clickable

var (
	directory = flag.String("d", "", "Path to directory to load")
)

func setupDefaults() {
	comparisonList = widget.List{
		List: layout.List{
			Axis: layout.Horizontal,
		},
	}
}

func main() {
	setupDefaults()

	flag.Parse()
	comparisons = shared.FindMetaFiles(*directory)
	for _, _ = range comparisons {
		comparisonButtons = append(comparisonButtons, widget.Clickable{})
	}

	testImage, _ = shared.LoadImage("../../testAssets/screenA.png")

	testList = widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Image-compare browser"))
		w.Option(app.Size(unit.Dp(1280), unit.Dp(720)))

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
		Bg:         color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF},
		Fg:         color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF},
		ContrastBg: color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xFF},
		ContrastFg: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return tabs(gtx, th)
				}),
				layout.Flexed(1, func(gtx C) D {
					return mainContent(gtx, th)
				}),
			)

			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func tabs(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			utils.FillLayout(gtx, th.Palette.Bg)
			ls := material.List(th, &comparisonList)

			return ls.Layout(gtx, len(comparisons), func(gtx C, index int) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
					if comparisonButtons[index].Pressed() {
						setComparison(comparisons[index])
					}
					lbl := material.Button(th, &comparisonButtons[index], comparisons[index].SourceA)
					return lbl.Layout(gtx)
				})
			})
		}),
	)
}

func mainContent(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return toolSidebar(gtx, th)
		}),
		layout.Flexed(1, func(gtx C) D {
			return pictureViewer(gtx)
		}),
	)
}

func toolSidebar(gtx C, th *material.Theme) D {
	gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			utils.FillLayout(gtx, color.NRGBA{R: 0x08, G: 0x08, B: 0x08, A: 0xFF})
			ls := material.List(th, &testList)

			return ls.Layout(gtx, 10, func(gtx C, index int) D {
				lbl := material.Body1(th, "item")
				return lbl.Layout(gtx)
			})
		}),
	)
}

func pictureViewer(gtx C) D {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			utils.FillLayout(gtx, color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF})
			return D{Size: utils.DrawImage(gtx, testImage)}
		}),
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(5)}.Layout,
		),
		layout.Rigid(func(gtx C) D {
			return pictureBrowser(gtx)
		}),
	)
}

func pictureBrowser(gtx C) D {
	gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))
	gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return imageButton1.Layout(gtx, func(gtx C) D {
				if imageButton1.Pressed() {
					fmt.Println("A")
				}
				utils.FillLayout(gtx, color.NRGBA{R: 0x08, G: 0x08, B: 0x08, A: 0xFF})
				return D{Size: utils.DrawImage(gtx, testImage)}
			})
		}),
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(5)}.Layout,
		),
		layout.Rigid(func(gtx C) D {
			return imageButton2.Layout(gtx, func(gtx C) D {
				if imageButton2.Pressed() {
					fmt.Println("B")
				}

				utils.FillLayout(gtx, color.NRGBA{R: 0x08, G: 0x08, B: 0x08, A: 0xFF})
				return D{Size: utils.DrawImage(gtx, testImage)}
			})
		}),
	)
}

func setComparison(comparison shared.Comparison) {
	testImage, _ = shared.LoadImage(comparison.Location + "/" + comparison.SourceA)
}
