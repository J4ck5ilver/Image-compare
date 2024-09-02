package main

import (
	"flag"
	"fmt"
	"ic/shared"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Pair struct {
	a, b interface{}
}

func validateArgs(args []string) (CompareData, error) {
	fs := flag.NewFlagSet("image-compare", flag.ContinueOnError)

	pathA := fs.String("A", "", "Filepath/directory A.")
	pathB := fs.String("B", "", "Filepath/directory B.")
	c := fs.String("c", "all", "Optional: Comparison options, [pixel,contrast,quad].")
	o := fs.String("o", "", "Optional: output directory.")

	if err := fs.Parse(args); err != nil {
		return CompareData{}, err
	}

	infoA, errA := os.Stat(*pathA)
	infoB, errB := os.Stat(*pathB)
	if errA != nil || errB != nil {
		return CompareData{}, fmt.Errorf("no sources provided")
	}

	if (infoA.IsDir() && !infoB.IsDir()) || (!infoA.IsDir() && infoB.IsDir()) {
		return CompareData{}, fmt.Errorf("sources differ, comparing file to directory")
	}

	data := CompareData{}
	data.SourceA = *pathA
	data.SourceB = *pathB
	data.IsDir = infoA.IsDir()
	data.ExportDest = *o
	data.Comparisons = shared.GetComparisons(*c)

	return data, nil
}

func loadImage(path string) (image.Image, error) {
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

func load(data CompareData) ([]CompareSet, error) {
	pairs := []Pair{}
	sets := []CompareSet{}

	orgExportDest := data.ExportDest

	if data.IsDir {
		itemsA, _ := os.ReadDir(data.SourceA)
		itemsB, _ := os.ReadDir(data.SourceB)
		for _, iA := range itemsA {
			if !iA.IsDir() {
				for _, iB := range itemsB {
					if !iA.IsDir() {
						if iA.Name() == iB.Name() {
							pairs = append(pairs, Pair{data.SourceA + "/" + iA.Name(), data.SourceB + "/" + iB.Name()})
						}
					}
				}
			}
		}
	} else {
		pairs = append(pairs, Pair{data.SourceA, data.SourceB})
	}

	for _, p := range pairs {
		imgA, err := loadImage(p.a.(string))
		if err != nil {
			return sets, err
		}

		imgB, err := loadImage(p.b.(string))
		if err != nil {
			return sets, err
		}

		data.SourceA = p.a.(string)
		data.SourceB = p.b.(string)

		if len(orgExportDest) > 1 {

			base := filepath.Base(p.a.(string))
			ext := filepath.Ext(base)
			dir := filepath.Join(orgExportDest, strings.TrimSuffix(base, ext))
			os.Mkdir(dir, os.ModePerm)

			data.ExportDest = dir
		}

		sets = append(sets, CompareSet{data, imgA, imgB})
	}

	return sets, nil
}

func run(args []string) []shared.ResultData {
	compareData, err := validateArgs(args)
	if err != nil {
		log.Fatal(err)
	}

	compareSets, err := load(compareData)
	if err != nil {
		log.Fatal(err)
	}

	if len(compareSets) == 0 {
		log.Fatal("Zero valid comparisons loaded")
	}

	results := []shared.ResultData{}
	for _, s := range compareSets {
		r, err := Compare(s)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, r...)
	}

	return results
}

func main() {
	results := run(os.Args[1:])
	for _, r := range results {
		fmt.Println(r.Location)
	}
}
