package utils

import (
	"image/color"
	"math"
)

const (
	Blend_Alpha      = "Alpha"
	Blend_Lighten    = "Lighten"
	Blend_Darken     = "Darken"
	Blend_Difference = "Difference"
)

func BlendAlpha(base *color.Color, overlay color.Color, alpha float64) color.Color {
	baseR, baseG, baseB, _ := (*base).RGBA()
	overlayR, overlayG, overlayB, overlayA := overlay.RGBA()

	overlayAlpha := float64(overlayA) / 65535.0 * alpha

	const map8bit = 1.0 / 65535.0 * 255
	r := uint8((float64(baseR)*(1-overlayAlpha) + float64(overlayR)*overlayAlpha) * map8bit)
	g := uint8((float64(baseG)*(1-overlayAlpha) + float64(overlayG)*overlayAlpha) * map8bit)
	b := uint8((float64(baseB)*(1-overlayAlpha) + float64(overlayB)*overlayAlpha) * map8bit)

	return color.RGBA{r, g, b, 255}
}

func BlendLighten(base *color.Color, overlay color.Color) color.Color {
	baseR, baseG, baseB, _ := (*base).RGBA()
	overlayR, overlayG, overlayB, _ := overlay.RGBA()

	const map8bit = 1.0 / 65535.0 * 255
	r := uint8(float64(max(baseR, overlayR)) * map8bit)
	g := uint8(float64(max(baseG, overlayG)) * map8bit)
	b := uint8(float64(max(baseB, overlayB)) * map8bit)

	return color.RGBA{r, g, b, 255}
}

func BlendDarken(base *color.Color, overlay color.Color) color.Color {
	baseR, baseG, baseB, _ := (*base).RGBA()
	overlayR, overlayG, overlayB, _ := overlay.RGBA()

	const map8bit = 1.0 / 65535.0 * 255
	r := uint8(float64(min(baseR, overlayR)) * map8bit)
	g := uint8(float64(min(baseG, overlayG)) * map8bit)
	b := uint8(float64(min(baseB, overlayB)) * map8bit)

	return color.RGBA{r, g, b, 255}
}

func BlendDifference(base *color.Color, overlay color.Color) color.Color {
	baseR, baseG, baseB, _ := (*base).RGBA()
	overlayR, overlayG, overlayB, _ := overlay.RGBA()

	const map8bit = 1.0 / 65535.0 * 255
	r := uint8(math.Abs(float64(baseR)-float64(overlayR)) * map8bit)
	g := uint8(math.Abs(float64(baseG)-float64(overlayG)) * map8bit)
	b := uint8(math.Abs(float64(baseB)-float64(overlayB)) * map8bit)

	return color.RGBA{r, g, b, 255}
}
