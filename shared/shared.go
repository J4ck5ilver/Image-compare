package shared

import (
	"image"
	"os"
	"strings"
)

type ComparisonType string

const (
	Pixel    ComparisonType = "pixel"
	Contrast ComparisonType = "contrast"
	Quad     ComparisonType = "quad"
	SSIM     ComparisonType = "ssim"
	MSE      ComparisonType = "mse"
)

type ResultData struct {
	Comparison string  `json:"comparison"`
	Index      float64 `json:"index"`
	NumFailed  int     `json:"numfailed"`
	Location   string  `json:"location"`
}

func GetComparisons(compString string) []ComparisonType {
	comparisons := []ComparisonType{}

	cOptions := strings.Split(compString, ",")
	if cOptions[0] == "all" {
		comparisons = []ComparisonType{Pixel, Contrast, Quad, SSIM, MSE}
	} else {
		for _, cO := range cOptions {
			switch ComparisonType(cO) {
			case Pixel:
				comparisons = append(comparisons, Pixel)
			case Contrast:
				comparisons = append(comparisons, Contrast)
			case Quad:
				comparisons = append(comparisons, Quad)
			case SSIM:
				comparisons = append(comparisons, SSIM)
			case MSE:
				comparisons = append(comparisons, MSE)
			}
		}
	}

	return comparisons
}

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
