package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	OPTION    = "Enter a sorting option:"
	DUPLICATE = "Check for duplicates?"
	DELETE    = "Delete files?"
	FORMAT    = "Enter file format:"
	DIRECTORY = "Directory is not specified"
	MENU      = "Size sorting options:\n1. Descending\n2. Ascending"
)

func main() {
	directoryName := getDirectoryName(os.Args)
	sizeIndex := indexFiles(directoryName)
	scanner := bufio.NewScanner(os.Stdin)
	fileExt := readInput(FORMAT, scanner)
	sizeIndex = selectExtension(sizeIndex, fileExt)
	sizeIndex = selectMultiple(sizeIndex)
	fmt.Println(MENU)
	option := readOption(scanner, OPTION, "1", "2")
	sortedSizes := sortSize(sizeIndex, option)
	showFiles(sizeIndex, sortedSizes)
	checkAnswer := readOption(scanner, DUPLICATE, "yes", "no")
	if "yes" == checkAnswer {
		duplicateFiles := filterDuplicateFiles(sizeIndex)
		orderedDuplicates := showDuplicateFiles(duplicateFiles, sortedSizes)
		deleteAnswer := readOption(scanner, DELETE, "yes", "no")
		if "yes" == deleteAnswer {
			fileNumbers := readFileNumbers(scanner, orderedDuplicates)
			freedSpace := deleteFiles(orderedDuplicates, fileNumbers)
			showFreedSpace(freedSpace)
		}
	}
}

func getDirectoryName(args []string) string {
	if len(args) < 2 {
		fmt.Println(DIRECTORY)
		os.Exit(0)
	}
	return args[1]
}

func indexFiles(directory string) map[int][]string {
	index := make(map[int][]string)
	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		if _, ok := index[int(info.Size())]; !ok {
			index[int(info.Size())] = make([]string, 0)
		}
		index[int(info.Size())] = append(index[int(info.Size())], path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return index
}

func readInput(prompt string, scanner *bufio.Scanner) string {
	fmt.Println(prompt)
	scanner.Scan()
	return scanner.Text()
}

func selectExtension(index map[int][]string, ext string) map[int][]string {
	if len(ext) == 0 {
		return index
	}
	newIndex := make(map[int][]string, len(index))
	for key, value := range index {
		for i := range value {
			extension := strings.TrimLeft(filepath.Ext(value[i]), ".")
			if extension == ext {
				newIndex[key] = append(newIndex[key], value[i])
			}
		}
	}
	return newIndex
}

func selectMultiple[T comparable](index map[T][]string) map[T][]string {
	for k := range index {
		if len(index[k]) < 2 {
			delete(index, k)
		}
	}
	return index
}

func readOption(scanner *bufio.Scanner, prompt string, options ...string) (option string) {
	for {
		option = readInput(prompt, scanner)
		for _, opt := range options {
			if option == opt {
				return
			}
		}
		fmt.Println("Wrong option")
	}
}

func sortSize(index map[int][]string, option string) []int {
	sizes := make([]int, len(index))
	i := 0
	for key := range index {
		sizes[i] = key
		i++
	}
	sort.Slice(sizes, func(i, j int) bool {
		if option == "1" {
			return sizes[i] > sizes[j]
		} else {
			return sizes[i] < sizes[j]
		}
	})
	return sizes
}

func showFiles(sizeIndex map[int][]string, sortedSizes []int) {
	for _, size := range sortedSizes {
		if paths, ok := sizeIndex[size]; ok {
			fmt.Printf("%d bytes\n", size)
			for _, path := range paths {
				fmt.Println(path)
			}
			fmt.Println()
		}
	}
}

func filterDuplicateFiles(index map[int][]string) map[int]map[string][]string {
	duplicateFiles := make(map[int]map[string][]string, len(index))
	for key, values := range index {
		hashes := findHashes(values)
		hashes = selectMultiple(hashes)
		if len(hashes) > 0 {
			duplicateFiles[key] = hashes
		}
	}
	return duplicateFiles
}

func findHashes(paths []string) map[string][]string {
	hashes := make(map[string][]string, len(paths))
	for _, path := range paths {
		hashString := hashFile(path)
		hashes[hashString] = append(hashes[hashString], path)
	}
	return hashes
}

func hashFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func showDuplicateFiles(duplicates map[int]map[string][]string, sizes []int) map[int]string {
	count := 1
	orderedFiles := make(map[int]string)
	for _, size := range sizes {
		if hashes, ok := duplicates[size]; ok {
			fmt.Printf("%d bytes\n", size)
			for hash, paths := range hashes {
				fmt.Printf("Hash: %s\n", hash)
				for _, path := range paths {
					orderedFiles[count] = path
					fmt.Printf("%d. %s\n", count, path)
					count++
				}
			}
		}
	}
	return orderedFiles
}

func readFileNumbers(scanner *bufio.Scanner, orderedDuplicates map[int]string) []int {
	for {
		fileNumbers := make([]int, 0)
		numberString := readInput("Enter file numbers to delete:", scanner)
		tokens := strings.Split(numberString, " ")
		mustRepeat := len(tokens) == 0
		for _, token := range tokens {
			number, err := strconv.Atoi(token)
			if err != nil {
				mustRepeat = true
				break
			}
			if _, ok := orderedDuplicates[number]; ok {
				fileNumbers = append(fileNumbers, number)
			} else {
				mustRepeat = true
				break
			}
		}
		if mustRepeat {
			fmt.Println("Wrong format")
			continue
		}
		return fileNumbers
	}
}

func deleteFiles(orderedDuplicates map[int]string, fileNumbers []int) (freedSpace int64) {
	for _, fileNumber := range fileNumbers {
		if path, ok := orderedDuplicates[fileNumber]; ok {
			fileInfo, err := os.Stat(path)
			if err != nil {
				log.Fatal(err)
			}
			fileSize := fileInfo.Size()
			freedSpace += fileSize
			err = os.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return
}

func showFreedSpace(space int64) {
	fmt.Printf("Total freed up space: %d bytes\n", space)
}
