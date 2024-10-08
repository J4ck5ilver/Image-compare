package main

import (
	"ic/browser/src/utils"
	"ic/shared"
	"testing"

	"gioui.org/widget"
)

func BenchmarkImageProcessing(b *testing.B) {
	img, _ := shared.LoadImage("../../testAssets/screenA.png")

	imagesActive = append(imagesActive, ImageSettings{
		Label:     "Test",
		Image:     img,
		BlendMode: widget.Enum{Value: utils.Blend_Alpha},
		Alpha:     widget.Float{Value: 1.0},
		R:         widget.Float{Value: 1.0},
		G:         widget.Float{Value: 1.0},
		B:         widget.Float{Value: 1.0},
	})

	for i := 0; i < b.N; i++ {
		updateViewImage()
	}

}
