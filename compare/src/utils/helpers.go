package utils

import (
	"image"
	"image/color"
)

func GetGrayValue(r uint32, g uint32, b uint32) float64 {
	gray := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
	return gray / float64(0xffff)
}

func ConvertToGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray.SetGray(x, y, color.Gray{Y: uint8((r*299 + g*587 + b*114 + 500) / 1000 >> 8)})
		}
	}
	return gray
}

func Mean(gray *image.Gray) float64 {
	var sum float64
	bounds := gray.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sum += float64(gray.GrayAt(x, y).Y)
		}
	}
	return sum / float64(bounds.Dx()*bounds.Dy())
}

func Variance(gray *image.Gray, mean float64) float64 {
	var sum float64
	bounds := gray.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			diff := float64(gray.GrayAt(x, y).Y) - mean
			sum += diff * diff
		}
	}
	return sum / float64(bounds.Dx()*bounds.Dy())
}

func Covariance(image1, image2 *image.Gray, mean1, mean2 float64) (float64, []float64) {
	var sum float64
	pixelSum := []float64{}
	bounds := image1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			diff1 := float64(image1.GrayAt(x, y).Y) - mean1
			diff2 := float64(image2.GrayAt(x, y).Y) - mean2

			d := diff1 * diff2
			pixelSum = append(pixelSum, d)
			sum += d
		}
	}
	return sum / float64(bounds.Dx()*bounds.Dy()), pixelSum
}
