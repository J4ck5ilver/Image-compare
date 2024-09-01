package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const contrastThreshold = 0.25
const quadThreshold = 0.5

type comparisonType string

const (
	Pixel    comparisonType = "pixel"
	Contrast comparisonType = "contrast"
	Quad     comparisonType = "quad"
)

type CompareData struct {
	SourceA     string
	SourceB     string
	IsDir       bool
	Comparisons []comparisonType
	ExportDest  string
}

type CompareSet struct {
	Data   CompareData
	ImageA image.Image
	ImageB image.Image
}

type ResultData struct {
	Comparison string  `json:"name"`
	Fraction   float64 `json:"value"`
}

func getGrayValue(r uint32, g uint32, b uint32) float64 {
	gray := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
	return gray / float64(0xffff)
}

func copy(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func export(data CompareData, img image.Image, result ResultData) error {
	_, err := os.Stat(data.ExportDest)
	if err != nil {
		return err
	}

	dest := filepath.Join(data.ExportDest, result.Comparison)
	os.Mkdir(dest, os.ModePerm)

	base := filepath.Base(data.SourceA)
	ext := filepath.Ext(base)
	filename := strings.TrimSuffix(base, ext)

	filename += "_diff.png"

	f, err := os.Create(filepath.Join(dest, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	copy(data.SourceA, filepath.Join(dest, filepath.Base(data.SourceA)))
	copy(data.SourceB, filepath.Join(dest, filepath.Base(data.SourceB)))

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		err = errors.New(fmt.Sprintln("Error marshaling JSON:", err))
		return err
	}

	err = os.WriteFile(filepath.Join(dest, "meta.json"), jsonData, 0644)
	if err != nil {
		err = errors.New(fmt.Sprintln("Error writing to file:", err))
		return err
	}

	return nil
}

func Compare(set CompareSet) ([]ResultData, error) {
	if len(set.Data.Comparisons) == 0 {
		return []ResultData{}, errors.New(fmt.Sprintln("No comparison type set."))
	}

	results := []ResultData{}
	for _, c := range set.Data.Comparisons {
		var frac float64
		var img image.Image
		var err error
		var result ResultData

		switch c {
		case Pixel:
			frac, img = PixelCompare(set)
			result = ResultData{"Pixel", frac}
		case Contrast:
			frac, img = ConstrastCompare(set)
			result = ResultData{"Contrast", frac}
		case Quad:
			frac, img, err = QuadCompare(set)
			if err != nil {
				return []ResultData{}, err
			}
			result = ResultData{"Quad", frac}
		default:
			return []ResultData{}, errors.New(fmt.Sprintf("Comparison type \"%v\" not supported.\n", c))
		}

		fmt.Printf("%s comparison: %f\n", result.Comparison, result.Fraction)
		results = append(results, result)

		if len(set.Data.ExportDest) > 0 {
			if err := export(set.Data, img, result); err != nil {
				return []ResultData{}, err
			}
		}
	}

	return results, nil
}

func PixelCompare(set CompareSet) (float64, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if set.ImageA.At(x, y) != set.ImageB.At(x, y) {
				numFailed++
				result.Set(x, y, color.White)
			} else {
				result.Set(x, y, color.Black)
			}

		}
	}

	fraction := float64(numFailed) / float64(w*h)
	return fraction, result
}

func ConstrastCompare(set CompareSet) (float64, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := set.ImageA.At(x, y).RGBA()
			grayA := getGrayValue(r, g, b)

			r, g, b, _ = set.ImageB.At(x, y).RGBA()
			grayB := getGrayValue(r, g, b)

			if math.Abs(grayA-grayB) > contrastThreshold {
				numFailed++
				c := color.Gray16{uint16(0xffff * math.Abs(grayA-grayB))}
				result.Set(x, y, c)
			} else {
				result.Set(x, y, color.Black)
			}

		}
	}

	fraction := float64(numFailed) / float64(w*h)
	return fraction, result
}

func QuadCompare(set CompareSet) (float64, image.Image, error) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	if w%2 != 0 || h%2 != 0 {
		err := errors.New(fmt.Sprintf("Quad comparison requires power of two resolution. Resolution: %d x %d\n", w, h))
		return 0.0, nil, err
	}

	numFailed := 0
	result := image.NewNRGBA(bounds)

	for x := 0; x < w; x += 2 {
		for y := 0; y < h; y += 2 {

			avgGrayA := 0.0
			avgGrayB := 0.0

			r, g, b, _ := set.ImageA.At(x, y).RGBA()
			avgGrayA += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x+1, y).RGBA()
			avgGrayA += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x, y+1).RGBA()
			avgGrayA += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageA.At(x+1, y+1).RGBA()
			avgGrayA += getGrayValue(r, g, b)

			r, g, b, _ = set.ImageB.At(x, y).RGBA()
			avgGrayB += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x+1, y).RGBA()
			avgGrayB += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x, y+1).RGBA()
			avgGrayB += getGrayValue(r, g, b)
			r, g, b, _ = set.ImageB.At(x+1, y+1).RGBA()
			avgGrayB += getGrayValue(r, g, b)

			avgGrayA /= 4.0
			avgGrayB /= 4.0

			if math.Abs(avgGrayA-avgGrayB) > quadThreshold {
				numFailed++
				c := color.Gray16{uint16(0xffff * math.Abs(avgGrayA-avgGrayB))}
				result.Set(x, y, c)
				result.Set(x+1, y, c)
				result.Set(x, y+1, c)
				result.Set(x+1, y+1, c)
			} else {
				result.Set(x, y, color.Black)
				result.Set(x+1, y, color.Black)
				result.Set(x, y+1, color.Black)
				result.Set(x+1, y+1, color.Black)
			}

		}
	}

	fraction := float64(numFailed) / float64(w*h/4)
	return fraction, result, nil
}
