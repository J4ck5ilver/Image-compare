package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"ic/browser/src/utils"
	"ic/shared"
)

type ClickableImage struct {
	Image  image.Image
	Widget widget.Clickable
	Label  string
}

type C = layout.Context
type D = layout.Dimensions

var imageMap = make(map[string]image.Image)

var imageView draw.Image
var imagesActive []image.Image

var imageBrowser []ClickableImage
var imageBrowserList widget.List

var comparisonList widget.List
var comparisons []shared.Comparison
var comparisonButtons []widget.Clickable

var (
	directory = flag.String("d", "", "Path to directory to load")
)

func setupDefaults() {
	imageBrowserList = widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

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

	if len(comparisons) == 0 {
		fmt.Println("No comparison data found, make sure directory contains a meta.json file")
	} else {
		for _, _ = range comparisons {
			comparisonButtons = append(comparisonButtons, widget.Clickable{})
		}

		setComparison(comparisons[0])
	}

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Image-compare browser"))
		w.Option(app.Size(unit.Dp(1280), unit.Dp(720)))

		if err := drawApp(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func drawApp(w *app.Window) error {
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
			return pictureViewer(gtx, th)
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
			//ls := material.List(th, &testList)

			//return ls.Layout(gtx, 10, func(gtx C, index int) D {
			//	lbl := material.Body1(th, "item")
			//	return lbl.Layout(gtx)
			//})
			return D{}
		}),
	)
}

func pictureViewer(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			utils.FillLayout(gtx, color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF})
			if imageView != nil {
				return D{Size: utils.DrawImage(gtx, imageView)}
			} else {
				return D{}
			}

		}),
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(5)}.Layout,
		),
		layout.Rigid(func(gtx C) D {
			return pictureBrowser(gtx, th)
		}),
	)
}

func pictureBrowser(gtx C, th *material.Theme) D {
	gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))
	gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			ls := material.List(th, &imageBrowserList)

			return ls.Layout(gtx, len(imageBrowser), func(gtx C, index int) D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						l := material.Body1(th, imageBrowser[index].Label)
						if isActive(imageBrowser[index].Image) {
							l.Color = color.NRGBA{R: 0x2, G: 0xAA, B: 0x2, A: 0xFF}
						}
						return l.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(100))
						return imageBrowser[index].Widget.Layout(gtx, func(gtx C) D {
							if imageBrowser[index].Widget.Pressed() {
								updateViewImage(imageBrowser[index].Image)
							}
							utils.FillLayout(gtx, color.NRGBA{R: 0x08, G: 0x08, B: 0x08, A: 0xFF})
							return D{Size: utils.DrawImage(gtx, imageBrowser[index].Image)}
						})
					}),
					layout.Rigid(
						layout.Spacer{Height: unit.Dp(5)}.Layout,
					),
				)
			})
		}),
	)
}

func setComparison(comparison shared.Comparison) {
	imagesActive = []image.Image{}
	imageBrowser = []ClickableImage{}

	filepaths := []string{
		comparison.Location + "/" + comparison.SourceA,
		comparison.Location + "/" + comparison.SourceB,
	}
	for _, r := range comparison.Results {
		filepaths = append(filepaths, comparison.Location+"/"+r.Comparison+".png")
	}

	for _, p := range filepaths {
		if _, exists := imageMap[p]; !exists {
			img, err := shared.LoadImage(p)

			if err == nil {
				imageMap[p] = img
				imageBrowser = append(imageBrowser, ClickableImage{Image: img})
			} else {
				fmt.Printf("Could not load image: %s\n", p)
			}
		} else {
			imageBrowser = append(imageBrowser, ClickableImage{Image: imageMap[p], Label: filepath.Base(p)})
		}

	}

	if len(imageBrowser) != 0 {
		updateViewImage(imageBrowser[0].Image)
	}

}
func updateViewImage(newImage image.Image) {
	imageFound := false
	for i, img := range imagesActive {
		if img == newImage {
			imagesActive = append(imagesActive[:i], imagesActive[i+1:]...)
			imageFound = true
		}
	}

	if !imageFound {
		imagesActive = append(imagesActive, newImage)
	}

	imageView = nil
	if len(imagesActive) > 0 {
		imageView = image.NewRGBA(imagesActive[0].Bounds())

		//alpha := 1.0
		for _, img := range imagesActive {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					overlayColor := img.At(x, y)
					baseColor := imageView.At(x, y)

					blendedColor := utils.BlendLighten(baseColor, overlayColor)

					imageView.Set(x, y, blendedColor)
				}
			}
		}
	}
}

func isActive(img image.Image) bool {
	for _, i := range imagesActive {
		if i == img {
			return true
		}
	}

	return false
}
