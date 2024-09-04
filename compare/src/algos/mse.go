package algos

import (
	"ic/compare/src/utils"
	"image"
	"image/color"
)

func MSE(set utils.CompareSet) (float64, int, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	var sumSquaredError float64
	result := image.NewNRGBA(bounds)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r1, g1, b1, _ := set.ImageA.At(x, y).RGBA()
			r2, g2, b2, _ := set.ImageB.At(x, y).RGBA()

			rf1, gf1, bf1 := float64(r1/0xffff), float64(g1/0xffff), float64(b1/0xffff)
			rf2, gf2, bf2 := float64(r2/0xffff), float64(g2/0xffff), float64(b2/0xffff)

			errR := rf1 - rf2
			errG := gf1 - gf2
			errB := bf1 - bf2
			sqe := (errR*errR + errG*errG + errB*errB) / 3

			result.Set(x, y, color.RGBA{uint8(sqe * 255), uint8(sqe * 255), uint8(sqe * 255), 255})

			if sqe != 0.0 {
				sumSquaredError += sqe
			}
		}
	}

	return 1.0 - (sumSquaredError / float64(w*h)), -1, result
}
