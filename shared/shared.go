package shared

import "strings"

type ComparisonType string

const (
	Pixel    ComparisonType = "pixel"
	Contrast ComparisonType = "contrast"
	Quad     ComparisonType = "quad"
)

type ResultData struct {
	Comparison string  `json:"comparison"`
	Fraction   float64 `json:"fraction"`
	NumFailed  int     `json:"numfailed"`
	Location   string  `json:"location"`
}

func GetComparisons(compString string) []ComparisonType {
	comparisons := []ComparisonType{}

	cOptions := strings.Split(compString, ",")
	if cOptions[0] == "all" {
		comparisons = []ComparisonType{Pixel, Contrast, Quad}
	} else {
		for _, cO := range cOptions {
			switch ComparisonType(cO) {
			case Pixel:
				comparisons = append(comparisons, Pixel)
			case Contrast:
				comparisons = append(comparisons, Contrast)
			case Quad:
				comparisons = append(comparisons, Quad)
			}
		}
	}

	return comparisons
}
