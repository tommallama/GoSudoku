package main

import (
	"fmt"
)

//TODO: Turn this bool flag into a command line arg with default set to false
var globalDebugMode = true

func main() {
	if globalDebugMode {
		fmt.Println("\n\nGo Sudoku Started in Debug Mode! \n For Speed Run, globalDebugMode = false\n ")
	}

	// File definition
	filename := "p096_sudoku.txt"

	var puzzleSet []Puzzle = loadFromFile(filename)

	// solvePuzzle(&puzzleSet[24])

	// prettyPrintPuzzle(puzzleSet[24])
	var successCount int = 0
	var averageScore float64 = 0.0

	for _, pz := range puzzleSet {
		solvePuzzle(&pz)

		if globalDebugMode {
			fmt.Printf("Puzzle: %v\t\tScore: %v\n", pz.name, pz.score)
			prettyPrintPuzzle(pz)
			fmt.Print("\n\n")
		}

		if pz.score == 100.0 {
			successCount++
		}
		averageScore += float64(pz.score)
	}

	averageScore = averageScore / float64(len(puzzleSet))
	fmt.Printf("Success: %v out of %v\n", successCount, len(puzzleSet))
	fmt.Printf("Average Score: %v%%\n", averageScore)
}
