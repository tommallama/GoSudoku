package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func loadFromFile(filename string) []Puzzle {
	readFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	var tempName string
	var tempData string
	var puzzleLineCount = 1
	var puzzleSet []Puzzle

	for fileScanner.Scan() {
		data := fileScanner.Text()
		if strings.Contains(data, "Grid") {
			tempName = data
			tempData = ""
			puzzleLineCount = 1
		} else {
			tempData += data
			puzzleLineCount++
		}
		if len(tempData) == 81 {
			var tempPuzzle Puzzle = createPuzzle(tempName, tempData)
			puzzleSet = append(puzzleSet, tempPuzzle)
		}

	}

	readFile.Close()
	return puzzleSet
}

func createPuzzle(name, data string) Puzzle {
	var p Puzzle
	p.name = name
	for i := 0; i < 81; i++ {
		u, _ := strconv.ParseUint(string(data[i]), 10, 8)
		p.puzzleStart[i] = uint8(u)
	}
	// fmt.Println("Name: ", name)
	// printPuzzle(p.puzzleStart)
	return p
}
