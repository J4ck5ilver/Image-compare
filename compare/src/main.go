package main

import (
	"flag"
	"fmt"
	"ic/compare/src/utils"
	"ic/shared"
	"log"
	"os"
	"path/filepath"
    "strings"
    "sync"
)

type Pair struct {
	a, b interface{}
}

func validateArgs(args []string) (utils.CompareData, error) {
	fs := flag.NewFlagSet("image-compare", flag.ContinueOnError)

	pathA := fs.String("A", "", "Filepath/directory A.")
	pathB := fs.String("B", "", "Filepath/directory B.")
	o := fs.String("o", "", "Optional: output directory.")
	c := fs.String("c", "all", "Optional: Comparison options, [pixel,contrast,quad,ssim,mse].")
    t := fs.Int("t", 1, "Number of threads to use.")

	if err := fs.Parse(args); err != nil {
		return utils.CompareData{}, err
	}

	infoA, errA := os.Stat(*pathA)
	infoB, errB := os.Stat(*pathB)
	if errA != nil || errB != nil {
		return utils.CompareData{}, fmt.Errorf("no sources provided")
	}

	if (infoA.IsDir() && !infoB.IsDir()) || (!infoA.IsDir() && infoB.IsDir()) {
		return utils.CompareData{}, fmt.Errorf("sources differ, comparing file to directory")
	}

	data := utils.CompareData{}
	data.SourceA = *pathA
	data.SourceB = *pathB
	data.IsDir = infoA.IsDir()
	data.ExportDest = *o
	data.Comparisons = shared.GetComparisons(*c)
    data.Threads = *t

	return data, nil
}


func load(data utils.CompareData) ([]utils.CompareSet, error) {
    if isFileComparison(data) {
        return handleFileComparison(data)
    }
    return handleDirectoryComparison(data)
}

func isFileComparison(data utils.CompareData) bool {
    infoA, errA := os.Stat(data.SourceA)
    infoB, errB := os.Stat(data.SourceB)
    return errA == nil && errB == nil && !infoA.IsDir() && !infoB.IsDir()
}

func handleFileComparison(data utils.CompareData) ([]utils.CompareSet, error) {
    pairs := []Pair{{data.SourceA, data.SourceB}}
    orgExportDest := data.ExportDest

    if len(orgExportDest) > 0 {
        os.MkdirAll(orgExportDest, os.ModePerm)
    }

    sets := []utils.CompareSet{}
    for _, p := range pairs {
        localData := data
        localData.SourceA = p.a.(string)
        localData.SourceB = p.b.(string)
        localData.ExportDest = orgExportDest
    
        sets = append(sets, utils.CompareSet{
            Data:       localData,
            ImageAPath: p.a.(string),
            ImageBPath: p.b.(string),
        })
    }
    return sets, nil
    
}

func walkSubdirectories(dirA, dirB, outDir string, allPairs *[]Pair) error {
    thesePairs := compareFilesInDirectories(dirA, dirB, outDir)
    *allPairs = append(*allPairs, thesePairs...)

    subdirsA, err := os.ReadDir(dirA)
    if err != nil {
        return err
    }
    subdirsB, err := os.ReadDir(dirB)
    if err != nil {
        return err
    }

    subdirsBMap := mapSubdirectories(subdirsB, dirB)

    for _, entryA := range subdirsA {
        if entryA.IsDir() {
            if matchingDirB, exists := subdirsBMap[entryA.Name()]; exists {
                subOutDir := filepath.Join(outDir, entryA.Name())

                if len(subOutDir) > 0 {
                    if err := os.MkdirAll(subOutDir, os.ModePerm); err != nil {
                        return err
                    }
                }

                err := walkSubdirectories(
                    filepath.Join(dirA, entryA.Name()),
                    matchingDirB,
                    subOutDir,
                    allPairs,
                )
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

func handleDirectoryComparison(data utils.CompareData) ([]utils.CompareSet, error) {
    pairs := []Pair{}

    if len(data.ExportDest) > 0 {
        os.MkdirAll(data.ExportDest, os.ModePerm)
    }

    err := walkSubdirectories(data.SourceA, data.SourceB, data.ExportDest, &pairs)
    if err != nil {
        return nil, err
    }

    return loadPairsIntoSets(pairs, data)
}

func mapSubdirectories(subdirs []os.DirEntry, basePath string) map[string]string {
    subdirsMap := make(map[string]string)
    for _, dir := range subdirs {
        if dir.IsDir() {
            subdirsMap[dir.Name()] = filepath.Join(basePath, dir.Name())
        }
    }
    return subdirsMap
}


func isImageFile(fileName string) bool {
    ext := strings.ToLower(filepath.Ext(fileName))
    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
        return true
    default:
        return false
    }
}

func compareFilesInDirectories(dirA, dirB, outputDir string) []Pair {
    filesA, _ := os.ReadDir(dirA)
    filesB, _ := os.ReadDir(dirB)
    
    filesBMap := make(map[string]string)
    for _, fileB := range filesB {
        if !fileB.IsDir() {
            filesBMap[fileB.Name()] = filepath.Join(dirB, fileB.Name())
        }
    }

    var pairs []Pair
    for _, fileA := range filesA {
        if !fileA.IsDir() && isImageFile(fileA.Name()) {
            if matchingFileB, exists := filesBMap[fileA.Name()]; exists {
                os.MkdirAll(outputDir, os.ModePerm)
                
                pairs = append(pairs, Pair{
                    a: filepath.Join(dirA, fileA.Name()),
                    b: matchingFileB,
                })
                
            }
        }
    }
    return pairs
}

func loadPairsIntoSets(pairs []Pair, data utils.CompareData) ([]utils.CompareSet, error) {
    sets := []utils.CompareSet{}

    originalSourceA, err := filepath.Abs(data.SourceA)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path for %s: %v", data.SourceA, err)
    }

    for _, p := range pairs {
        //imgA, err := shared.LoadImage(p.a.(string))
        if err != nil {
            return nil, fmt.Errorf("failed to load image %s: %v", p.a, err)
        }

        //imgB, err := shared.LoadImage(p.b.(string))
        if err != nil {
            return nil, fmt.Errorf("failed to load image %s: %v", p.b, err)
        }

        absA, err := filepath.Abs(p.a.(string))
        if err != nil {
            return nil, fmt.Errorf("failed to get absolute path for %s: %v", p.a, err)
        }

        relativePath, err := filepath.Rel(originalSourceA, absA)
        if err != nil {
            return nil, fmt.Errorf("failed to compute relative path for %s: %v", p.a, err)
        }

        dirPart := filepath.Dir(relativePath)
        filePart := filepath.Base(relativePath)
        baseName := filePart[:len(filePart)-len(filepath.Ext(filePart))]
        finalExportPath := filepath.Join(data.ExportDest, dirPart, baseName)

        if err := os.MkdirAll(finalExportPath, os.ModePerm); err != nil {
            return nil, fmt.Errorf("failed to create directory %s: %v", finalExportPath, err)
        }

        localData := data
        localData.SourceA = p.a.(string)
        localData.SourceB = p.b.(string)
        localData.ExportDest = finalExportPath

        sets = append(sets, utils.CompareSet{
            Data:       localData,
            ImageAPath: p.a.(string),
            ImageBPath: p.b.(string),
        })
        
    }

    return sets, nil
}




func run(args []string) []shared.Comparison {
    compareData, err := validateArgs(args)
    if err != nil {
        log.Fatal(err)
    }

    compareSets, err := load(compareData)
    if err != nil {
        log.Fatal(err)
    }

    sem := make(chan struct{}, compareData.Threads)

    comparisons := make([]shared.Comparison, len(compareSets))

    var wg sync.WaitGroup

    for i, s := range compareSets {
        wg.Add(1)
        sem <- struct{}{}

        go func(i int, s utils.CompareSet) {
            defer wg.Done()
            defer func() { <-sem }()

            imgA, err := shared.LoadImage(s.ImageAPath)
            if err != nil {
                log.Fatal(err)
            }
            imgB, err := shared.LoadImage(s.ImageBPath)
            if err != nil {
                log.Fatal(err)
            }

            s.ImageA = imgA
            s.ImageB = imgB

            c, err := Compare(s)
            if err != nil {
                log.Fatal(err)
            }

            comparisons[i] = c

        }(i, s)
    }

    wg.Wait()

    return comparisons
}

func main() {
	comparisons := run(os.Args[1:])
	for _, c := range comparisons {
		if len(c.Location) > 1 {
			fmt.Println(c.Location)
		} else {
			fmt.Println(c.Results)
		}
	}
}
