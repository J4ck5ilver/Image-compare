package main

import (
	"encoding/json"
	"fmt"
	"ic/compare/src/algos"
	"ic/compare/src/utils"
	"ic/shared"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const debug = false

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

func export(data utils.CompareData, images []image.Image, comparison shared.Comparison) error {
	_, err := os.Stat(comparison.Location)
	if err != nil {
		return err
	}

	os.Mkdir(comparison.Location, os.ModePerm)

	copy(data.SourceA, filepath.Join(comparison.Location, filepath.Base(comparison.SourceA)))
	copy(data.SourceB, filepath.Join(comparison.Location, filepath.Base(comparison.SourceB)))

	for i, r := range comparison.Results {
		filename := r.Comparison + ".png"

		f, err := os.Create(filepath.Join(comparison.Location, filename))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := png.Encode(f, images[i]); err != nil {
			return err
		}
	}

	jsonData, err := json.MarshalIndent(comparison, "", "  ")
	if err != nil {
		err = fmt.Errorf("error marshaling json: %v", err)
		return err
	}

	err = os.WriteFile(filepath.Join(comparison.Location, "meta.json"), jsonData, 0644)
	if err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
		return err
	}

	return nil
}

func Compare(set utils.CompareSet) (shared.Comparison, error) {
	if len(set.Data.Comparisons) == 0 {
		return shared.Comparison{}, fmt.Errorf("no comparison type set")
	}

	results := []shared.ResultData{}
	images := []image.Image{}
	for _, c := range set.Data.Comparisons {
		var index float64
		var numFailed int
		var img image.Image
		var err error
		var result shared.ResultData

		switch c {
		case shared.Pixel:
			index, numFailed, img = algos.PixelCompare(set)
			result = shared.ResultData{Comparison: string(shared.Pixel), Index: index, NumFailed: numFailed}
		case shared.Contrast:
			index, numFailed, img = algos.ConstrastCompare(set)
			result = shared.ResultData{Comparison: string(shared.Contrast), Index: index, NumFailed: numFailed}
		case shared.Quad:
			index, numFailed, img, err = algos.QuadCompare(set)
			if err != nil {
				fmt.Println(err)
				continue
			}
			result = shared.ResultData{Comparison: string(shared.Quad), Index: index, NumFailed: numFailed}
		case shared.SSIM:
			index, numFailed, img = algos.SSIM(set)
			result = shared.ResultData{Comparison: string(shared.SSIM), Index: index, NumFailed: numFailed}
		case shared.MSE:
			index, numFailed, img = algos.MSE(set)
			result = shared.ResultData{Comparison: string(shared.MSE), Index: index, NumFailed: numFailed}
		default:
			return shared.Comparison{}, fmt.Errorf("comparison type \"%v\" not supported", c)
		}

		if debug {
			fmt.Printf("%s comparison: %f\n", result.Comparison, result.Index)
		}
		results = append(results, result)
		images = append(images, img)
	}

	comparison := shared.Comparison{
		Location: set.Data.ExportDest,
		SourceA:  filepath.Base(set.Data.SourceA),
		SourceB:  filepath.Base(set.Data.SourceB),
		Results:  results,
	}

	if comparison.SourceA == comparison.SourceB {
		ext := filepath.Ext(comparison.SourceA)
		comparison.SourceA = strings.TrimSuffix(comparison.SourceA, ext) + "_A" + ext

		ext = filepath.Ext(comparison.SourceB)
		comparison.SourceB = strings.TrimSuffix(comparison.SourceB, ext) + "_B" + ext
	}

	if len(set.Data.ExportDest) > 0 {
		if err := export(set.Data, images, comparison); err != nil {
			return shared.Comparison{}, err
		}
	}

	return comparison, nil
}
