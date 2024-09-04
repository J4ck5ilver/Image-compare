package algos

import (
	"ic/compare/src/utils"
	"image"
	"image/color"
	"math"
)

func SSIM(set utils.CompareSet) (float64, int, image.Image) {

	gray1 := utils.ConvertToGray(set.ImageA)
	gray2 := utils.ConvertToGray(set.ImageB)

	mean1 := utils.Mean(gray1)
	mean2 := utils.Mean(gray2)

	variance1 := utils.Variance(gray1, mean1)
	variance2 := utils.Variance(gray2, mean2)

	cov, pixels := utils.Covariance(gray1, gray2, mean1, mean2)

	c1 := 6.5025
	c2 := 58.5225

	ssim := ((2*mean1*mean2 + c1) * (2*cov + c2)) /
		((mean1*mean1 + mean2*mean2 + c1) * (variance1 + variance2 + c2))

	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	result := image.NewNRGBA(bounds)

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := math.Abs(((2*mean1*mean2 + c1) * (2*pixels[i] + c2)) /
				((mean1*mean1 + mean2*mean2 + c1) * (variance1 + variance2 + c2)))
			i++

			result.Set(x, y, color.RGBA{uint8(p * 255), uint8(p * 255), uint8(p * 255), 255})
		}
	}

	return math.Abs(ssim), -1, result
}
