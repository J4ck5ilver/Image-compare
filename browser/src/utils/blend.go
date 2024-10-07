package utils

import (
	"image/color"
	"math"
)

const (
	Blend_Alpha   = "Alpha"
	Blend_Lighten = "Lighten"
)

func BlendAlpha(base color.Color, overlay color.Color, alpha float64) color.Color {
	baseR, baseG, baseB, baseA := base.RGBA()
	overlayR, overlayG, overlayB, overlayA := overlay.RGBA()

	overlayAlpha := float64(overlayA) / 65535.0 * alpha

	r := uint8(math.Round((float64(baseR)/65535.0*(1-overlayAlpha) + float64(overlayR)/65535.0*overlayAlpha) * 255))
	g := uint8(math.Round((float64(baseG)/65535.0*(1-overlayAlpha) + float64(overlayG)/65535.0*overlayAlpha) * 255))
	b := uint8(math.Round((float64(baseB)/65535.0*(1-overlayAlpha) + float64(overlayB)/65535.0*overlayAlpha) * 255))
	a := uint8(math.Round((float64(baseA)/65535.0*(1-overlayAlpha) + float64(overlayA)/65535.0*overlayAlpha) * 255))

	return color.RGBA{r, g, b, a}
}

func BlendLighten(base color.Color, overlay color.Color) color.Color {
	baseR, baseG, baseB, _ := base.RGBA()
	overlayR, overlayG, overlayB, _ := overlay.RGBA()

	r := uint8(float64(max(baseR, overlayR)) / 65535.0 * 255)
	g := uint8(float64(max(baseG, overlayG)) / 65535.0 * 255)
	b := uint8(float64(max(baseB, overlayB)) / 65535.0 * 255)

	return color.RGBA{r, g, b, 255}
}
