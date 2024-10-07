package utils

import (
	"image"
)

func GetGrayValue(r uint32, g uint32, b uint32) float64 {
	gray := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
	return gray / float64(0xffff)
}

func ConvertToGray(img image.Image) []float64 {
	bounds := img.Bounds()
	graySlice := []float64{}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			graySlice = append(graySlice, GetGrayValue(r, g, b))
		}
	}
	return graySlice
}

func Mean(graySlice []float64) float64 {
	var sum float64
	for _, g := range graySlice {
		sum += g
	}

	return sum / float64(len(graySlice))
}

func Variance(graySlice []float64, mean float64) float64 {
	var sum float64
	for _, g := range graySlice {
		diff := g - mean
		sum += diff * diff
	}

	return sum / float64(len(graySlice))
}

func Covariance(graySlice1, graySlice2 []float64, mean1, mean2 float64) (float64, []float64) {
	var sum float64
	pixelSum := []float64{}
	for i, _ := range graySlice1 {
		diff1 := graySlice1[i] - mean1
		diff2 := graySlice2[i] - mean2

		d := diff1 * diff2
		pixelSum = append(pixelSum, d)
		sum += d
	}

	return sum / float64(len(graySlice1)), pixelSum
}
