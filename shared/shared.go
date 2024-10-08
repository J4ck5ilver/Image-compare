package shared

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

type ComparisonType string

const (
	Pixel    ComparisonType = "pixel"
	Contrast ComparisonType = "contrast"
	Quad     ComparisonType = "quad"
	SSIM     ComparisonType = "ssim"
	MSE      ComparisonType = "mse"
)

type Comparison struct {
	Location string       `json:"location"`
	SourceA  string       `json:"source_a"`
	SourceB  string       `json:"source_b"`
	Results  []ResultData `json:"results"`
}

type ResultData struct {
	Comparison string  `json:"comparison"`
	Index      float64 `json:"index"`
	NumFailed  int     `json:"numfailed"`
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
	return loadImage(path, 1.0)
}

func LoadImageScaled(path string, scale float64) (image.Image, error) {
	return loadImage(path, scale)
}

func loadImage(path string, scale float64) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if scale != 1.0 {
		img = scaleImage(img, scale)
	}

	return img, nil
}

func scaleImage(img image.Image, scale float64) image.Image {
	bounds := img.Bounds()
	w := float64(bounds.Max.X) * scale
	h := float64(bounds.Max.Y) * scale

	dst := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	return dst
}

func FindMetaFiles(dir string) []Comparison {
	comparisons := []Comparison{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "meta.json" {
			var r Comparison

			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}

			err = json.Unmarshal(data, &r)
			if err != nil {
				return fmt.Errorf("error unmarshalling json: %v", err)
			}

			comparisons = append(comparisons, r)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return comparisons
}
