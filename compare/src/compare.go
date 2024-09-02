package main

import (
	"encoding/json"
	"fmt"
	"ic/shared"
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
const debug = false

type CompareData struct {
	SourceA     string
	SourceB     string
	IsDir       bool
	Comparisons []shared.ComparisonType
	ExportDest  string
}

type CompareSet struct {
	Data   CompareData
	ImageA image.Image
	ImageB image.Image
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

func export(data CompareData, img image.Image, result shared.ResultData) error {
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
		err = fmt.Errorf("error marshaling json: %v", err)
		return err
	}

	err = os.WriteFile(filepath.Join(dest, "meta.json"), jsonData, 0644)
	if err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
		return err
	}

	return nil
}

func Compare(set CompareSet) ([]shared.ResultData, error) {
	if len(set.Data.Comparisons) == 0 {
		return []shared.ResultData{}, fmt.Errorf("no comparison type set.")
	}

	results := []shared.ResultData{}
	for _, c := range set.Data.Comparisons {
		var frac float64
		var numFailed int
		var img image.Image
		var err error
		var result shared.ResultData

		switch c {
		case shared.Pixel:
			frac, numFailed, img = PixelCompare(set)
			result = shared.ResultData{string(shared.Pixel), frac, numFailed, set.Data.ExportDest}
		case shared.Contrast:
			frac, numFailed, img = ConstrastCompare(set)
			result = shared.ResultData{string(shared.Contrast), frac, numFailed, set.Data.ExportDest}
		case shared.Quad:
			frac, numFailed, img, err = QuadCompare(set)
			if err != nil {
				return []shared.ResultData{}, err
			}
			result = shared.ResultData{string(shared.Quad), frac, numFailed, set.Data.ExportDest}
		default:
			return []shared.ResultData{}, fmt.Errorf("comparison type \"%v\" not supported", c)
		}

		if debug {
			fmt.Printf("%s comparison: %f\n", result.Comparison, result.Fraction)
		}
		results = append(results, result)

		if len(set.Data.ExportDest) > 0 {
			if err := export(set.Data, img, result); err != nil {
				return []shared.ResultData{}, err
			}
		}
	}

	return results, nil
}

func PixelCompare(set CompareSet) (float64, int, image.Image) {
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

func ConstrastCompare(set CompareSet) (float64, int, image.Image) {
	bounds := set.ImageA.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	numMatches := 0
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
				numMatches++
				result.Set(x, y, color.Black)
			}

		}
	}

	fraction := float64(numMatches) / float64(w*h)
	return fraction, numFailed, result
}

func QuadCompare(set CompareSet) (float64, int, image.Image, error) {
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
				numMatches++
				result.Set(x, y, color.Black)
				result.Set(x+1, y, color.Black)
				result.Set(x, y+1, color.Black)
				result.Set(x+1, y+1, color.Black)
			}

		}
	}

	fraction := float64(numMatches) / float64(w*h/4)
	return fraction, numFailed, result, nil
}
