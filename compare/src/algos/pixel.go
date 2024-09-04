package algos

import (
	"ic/compare/src/utils"
	"image"
	"image/color"
)

func PixelCompare(set utils.CompareSet) (float64, int, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	numMatches := 0
	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if set.ImageA.At(x, y) != set.ImageB.At(x, y) {
				numFailed++
				result.Set(x, y, color.White)
			} else {
				numMatches++
				result.Set(x, y, color.Black)
			}

		}
	}

	fraction := float64(numMatches) / float64(w*h)
	return fraction, numFailed, result
}
