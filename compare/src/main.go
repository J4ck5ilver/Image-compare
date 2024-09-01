package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

func validateArgs() (CompareData, error) {
	pathA := flag.String("A", "", "Filepath/directory A")
	pathB := flag.String("B", "", "Filepath/directory B")
	c := flag.String("c", "", "Optional: Comparison options, [pixel,contrast,quad] Default: all")
	o := flag.String("o", "", "Optional: output directory")

	flag.Parse()

	infoA, errA := os.Stat(*pathA)
	infoB, errB := os.Stat(*pathB)
	if errA != nil || errB != nil {
		return CompareData{}, errors.New(fmt.Sprintln("No sources provided."))
	}

	if (infoA.IsDir() && !infoB.IsDir()) || (!infoA.IsDir() && infoB.IsDir()) {
		return CompareData{}, errors.New(fmt.Sprintln("Sources differ, comparing file to directory."))
	}

	data := CompareData{}
	data.SourceA = *pathA
	data.SourceB = *pathB
	data.IsDir = infoA.IsDir()
	data.ExportDest = *o

	cOptions := strings.Split(*c, ",")
	if len(cOptions) == 0 {
		data.Comparisons = []comparisonType{Pixel, Contrast, Quad}
	} else {
		for _, cO := range cOptions {
			switch comparisonType(cO) {
			case Pixel:
				data.Comparisons = append(data.Comparisons, Pixel)
			case Contrast:
				data.Comparisons = append(data.Comparisons, Contrast)
			case Quad:
				data.Comparisons = append(data.Comparisons, Quad)
			}
		}
	}

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
	sets := []CompareSet{}

	if data.IsDir {
		return sets, errors.New(fmt.Sprintln("Directory load not supported"))
	} else {
		imgA, err := loadImage(data.SourceA)
		if err != nil {
			return sets, err
		}

		imgB, err := loadImage(data.SourceB)
		if err != nil {
			return sets, err
		}

		sets = []CompareSet{CompareSet{data, imgA, imgB}}
	}

	return sets, nil
}

func main() {
	compareData, err := validateArgs()
	if err != nil {
		log.Fatal(err)
		return
	}

	compareSets, err := load(compareData)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(compareSets) == 0 {
		log.Println("Zero valid comparisons loaded")
	}

	for _, s := range compareSets {
		err = Compare(compareData, s)
		if err != nil {
			log.Fatal(err)
		}
	}
}
