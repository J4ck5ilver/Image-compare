package algos

import (
	"ic/compare/src/utils"
	"image"
	"image/color"
	"math"
)

const contrastThreshold = 0.25

func ConstrastCompare(set utils.CompareSet) (float64, int, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	numMatches := 0
	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := set.ImageA.At(x, y).RGBA()
			grayA := utils.GetGrayValue(r, g, b)

			r, g, b, _ = set.ImageB.At(x, y).RGBA()
			grayB := utils.GetGrayValue(r, g, b)

			if math.Abs(grayA-grayB) > contrastThreshold {
				numFailed++
				c := color.Gray16{uint16(0xffff * math.Abs(grayA-grayB))}
				result.Set(x, y, c)
			} else {
				numMatches++
				result.Set(x, y, color.Black)
			}

		}
	}

	fraction := float64(numMatches) / float64(w*h)
	return fraction, numFailed, result
}
