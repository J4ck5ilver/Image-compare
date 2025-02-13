package main

import (
	"flag"
	"fmt"
	"os"
	"ic/shared"
)

var (
	comparison = flag.String("c", "all", "Optional: Comparison options, [pixel,contrast,quad,ssim,mse].")
	index      = flag.Float64("i", 1.0, "Optional: Index threshold.")
	numFailed  = flag.Int("n", 0, "Optional: Num failed points.")
	directory  = flag.String("d", "", "Optional: Path to directory to filter.")
)

func filterComparisons(comparisons []shared.Comparison) []shared.Comparison {
	filtered := []shared.Comparison{}

	comp := shared.GetComparisons(*comparison)

	for _, c := range comparisons {
		for _, r := range c.Results {
			compareMatch := false
			for _, c := range comp {
				if c == shared.ComparisonType(r.Comparison) {
					compareMatch = true
					break
				}
			}

			if !compareMatch {
				continue
			}

			if r.Index > *index {
				continue
			}

			if *numFailed != 0 {
				if r.NumFailed > *numFailed || r.NumFailed == -1 {
					continue
				}
			}

			filtered = append(filtered, c)
			break
		}
	}

	return filtered
}


func main() {
    flag.Parse()

    comparisons := shared.FindMetaFiles(*directory)
    comparisons = filterComparisons(comparisons)

	for _, c := range comparisons {
        fmt.Println(c.Location)
    }

    os.Exit(0)
}