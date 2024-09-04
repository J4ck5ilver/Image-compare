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
}

type CompareSet struct {
	Data   CompareData
	ImageA image.Image
	ImageB image.Image
}
