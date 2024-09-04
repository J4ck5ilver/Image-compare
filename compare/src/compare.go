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

func export(data utils.CompareData, images []image.Image, results []shared.ResultData) error {
	_, err := os.Stat(data.ExportDest)
	if err != nil {
		return err
	}

	os.Mkdir(data.ExportDest, os.ModePerm)

	for i, r := range results {
		base := filepath.Base(data.SourceA)
		ext := filepath.Ext(base)
		filename := strings.TrimSuffix(base, ext)

		filename += r.Comparison + "_diff.png"

		f, err := os.Create(filepath.Join(data.ExportDest, filename))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := png.Encode(f, images[i]); err != nil {
			return err
		}
	}

	copy(data.SourceA, filepath.Join(data.ExportDest, filepath.Base(data.SourceA)))
	copy(data.SourceB, filepath.Join(data.ExportDest, filepath.Base(data.SourceB)))

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		err = fmt.Errorf("error marshaling json: %v", err)
		return err
	}

	err = os.WriteFile(filepath.Join(data.ExportDest, "meta.json"), jsonData, 0644)
	if err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
		return err
	}

	return nil
}

func Compare(set utils.CompareSet) ([]shared.ResultData, error) {
	if len(set.Data.Comparisons) == 0 {
		return []shared.ResultData{}, fmt.Errorf("no comparison type set.")
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
			result = shared.ResultData{string(shared.Pixel), index, numFailed, set.Data.ExportDest}
		case shared.Contrast:
			index, numFailed, img = algos.ConstrastCompare(set)
			result = shared.ResultData{string(shared.Contrast), index, numFailed, set.Data.ExportDest}
		case shared.Quad:
			index, numFailed, img, err = algos.QuadCompare(set)
			if err != nil {
				fmt.Println(err)
				continue
			}
			result = shared.ResultData{string(shared.Quad), index, numFailed, set.Data.ExportDest}
		case shared.SSIM:
			index, numFailed, img = algos.SSIM(set)
			result = shared.ResultData{string(shared.SSIM), index, numFailed, set.Data.ExportDest}
		case shared.MSE:
			index, numFailed, img = algos.MSE(set)
			result = shared.ResultData{string(shared.MSE), index, numFailed, set.Data.ExportDest}
		default:
			return []shared.ResultData{}, fmt.Errorf("comparison type \"%v\" not supported", c)
		}

		if debug {
			fmt.Printf("%s comparison: %f\n", result.Comparison, result.Index)
		}
		results = append(results, result)
		images = append(images, img)
	}

	if len(set.Data.ExportDest) > 0 {
		if err := export(set.Data, images, results); err != nil {
			return []shared.ResultData{}, err
		}
	}

	return results, nil
}
