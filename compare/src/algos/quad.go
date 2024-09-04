package algos

import (
	"fmt"
	"ic/compare/src/utils"
	"image"
	"image/color"
	"math"
)

const quadThreshold = 0.5

func QuadCompare(set utils.CompareSet) (float64, int, image.Image, error) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	if w%2 != 0 || h%2 != 0 {
		err := fmt.Errorf("quad comparison requires power of two resolution. resolution: %d x %d", w, h)
		return 0.0, 0, nil, err
	}

	numMatches := 0
	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x += 2 {
		for y := 0; y < h; y += 2 {

			avgGrayA := 0.0
			avgGrayB := 0.0

			r, g, b, _ := set.ImageA.At(x, y).RGBA()
			avgGrayA += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x+1, y).RGBA()
			avgGrayA += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x, y+1).RGBA()
			avgGrayA += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x+1, y+1).RGBA()
			avgGrayA += utils.GetGrayValue(r, g, b)

			r, g, b, _ = set.ImageB.At(x, y).RGBA()
			avgGrayB += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x+1, y).RGBA()
			avgGrayB += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x, y+1).RGBA()
			avgGrayB += utils.GetGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x+1, y+1).RGBA()
			avgGrayB += utils.GetGrayValue(r, g, b)

			avgGrayA /= 4.0
			avgGrayB /= 4.0

			if math.Abs(avgGrayA-avgGrayB) > quadThreshold {
				numFailed += 4
				c := color.Gray16{uint16(0xffff * math.Abs(avgGrayA-avgGrayB))}
				result.Set(x, y, c)
				result.Set(x+1, y, c)
				result.Set(x, y+1, c)
				result.Set(x+1, y+1, c)
			} else {
				numMatches += 4
				result.Set(x, y, color.Black)
				result.Set(x+1, y, color.Black)
				result.Set(x, y+1, color.Black)
				result.Set(x+1, y+1, color.Black)
			}

		}
	}

	fraction := float64(numMatches) / float64(w*h)
	return fraction, numFailed, result, nil
}
