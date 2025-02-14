package utils

import (
	"ic/shared"
	"image"
)

type CompareData struct {
	SourceA     string
	SourceB     string
	IsDir       bool
	Comparisons []shared.ComparisonType
	ExportDest  string
	Threads int
}

type CompareSet struct {
	Data   CompareData
	ImageA image.Image
	ImageB image.Image
	ImageAPath string
    ImageBPath string
}
