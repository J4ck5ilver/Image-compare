package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"ic/shared"
	"log"
	"os"
	"path/filepath"
)

var (
	comparison = flag.String("c", "all", "Optional: Comparison options, [pixel,contrast,quad,ssim,mse].")
	index      = flag.Float64("i", 1.0, "Optional: Index threshold.")
	numFailed  = flag.Int("n", 0, "Optional: Num failed points.")
	directory  = flag.String("d", "", "Optional: Path to directory to filter.")
)

func findMetaFiles(dirs []string) []shared.ResultData {
	results := []shared.ResultData{}

	for _, d := range dirs {
		err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Base(path) == "meta.json" {
				var r []shared.ResultData

				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file: %v", err)
				}

				err = json.Unmarshal(data, &r)
				if err != nil {
					return fmt.Errorf("error unmarshalling json: %v", err)
				}

				results = append(results, r...)
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	return results
}

func filterResults(results []shared.ResultData) []shared.ResultData {
	filtered := []shared.ResultData{}

	comp := shared.GetComparisons(*comparison)

	for _, r := range results {
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

		filtered = append(filtered, r)
	}

	return filtered
}

func main() {
	flag.Parse()

	directories := []string{}

	if len(*directory) < 1 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			directories = append(directories, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading input: %v\n", err)
		}
	} else {
		directories = append(directories, *directory)
	}

	results := findMetaFiles(directories)

	results = filterResults(results)

	for _, r := range results {
		fmt.Println(r.Location)
	}
}
