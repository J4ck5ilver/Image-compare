package utils

import (
	"image/color"
	"math"
)

func BlendAlpha(base color.Color, overlay color.Color, alpha float64) color.Color {
	baseR, baseG, baseB, baseA := base.RGBA()
	overlayR, overlayG, overlayB, overlayA := overlay.RGBA()

	overlayAlpha := float64(overlayA) / 65535.0 * alpha

	r := uint8(math.Round((float64(baseR)*(1-overlayAlpha) + float64(overlayR)*overlayAlpha) / 256))
	g := uint8(math.Round((float64(baseG)*(1-overlayAlpha) + float64(overlayG)*overlayAlpha) / 256))
	b := uint8(math.Round((float64(baseB)*(1-overlayAlpha) + float64(overlayB)*overlayAlpha) / 256))
	a := uint8(math.Round((float64(baseA)*(1-overlayAlpha) + float64(overlayA)*overlayAlpha) / 256))

	return color.RGBA{r, g, b, a}
}

func BlendLighten(base color.Color, overlay color.Color) color.Color {
	baseR, baseG, baseB, _ := base.RGBA()
	overlayR, overlayG, overlayB, _ := overlay.RGBA()

	r := uint8(max(baseR, overlayR) / 256)
	g := uint8(max(baseG, overlayG) / 256)
	b := uint8(max(baseB, overlayB) / 256)

	return color.RGBA{r, g, b, 255}
}
