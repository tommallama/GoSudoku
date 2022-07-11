package main

import (
	"fmt"
	"strconv"
)

// TODO: test these with a function instead of lookup tables for speed
// Equivalent function value=index/9
var cellToRowLookup = [81]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8}

// Equivalent function value=index%9
var cellToColLookup = [81]uint8{
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	0, 1, 2, 3, 4, 5, 6, 7, 8}

var cellToBoxLookup = [81]uint8{
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8}

var boxIndices = [][]uint8{
	{0, 1, 2, 9, 10, 11, 18, 19, 20},
	{3, 4, 5, 12, 13, 14, 21, 22, 23},
	{6, 7, 8, 15, 16, 17, 24, 25, 26},
	{27, 28, 29, 36, 37, 38, 45, 46, 47},
	{30, 31, 32, 39, 40, 41, 48, 49, 50},
	{33, 34, 35, 42, 43, 44, 51, 52, 53},
	{54, 55, 56, 63, 64, 65, 72, 73, 74},
	{57, 58, 59, 66, 67, 68, 75, 76, 77},
	{60, 61, 62, 69, 70, 71, 78, 79, 80}}

var rowIndices = [][]uint8{
	{0, 1, 2, 3, 4, 5, 6, 7, 8},
	{9, 10, 11, 12, 13, 14, 15, 16, 17},
	{18, 19, 20, 21, 22, 23, 24, 25, 26},
	{27, 28, 29, 30, 31, 32, 33, 34, 35},
	{36, 37, 38, 39, 40, 41, 42, 43, 44},
	{45, 46, 47, 48, 49, 50, 51, 52, 53},
	{54, 55, 56, 57, 58, 59, 60, 61, 62},
	{63, 64, 65, 66, 67, 68, 69, 70, 71},
	{72, 73, 74, 75, 76, 77, 78, 79, 80}}

var colIndices = [][]uint8{
	{0, 9, 18, 27, 36, 45, 54, 63, 72},
	{1, 10, 19, 28, 37, 46, 55, 64, 73},
	{2, 11, 20, 29, 38, 47, 56, 65, 74},
	{3, 12, 21, 30, 39, 48, 57, 66, 75},
	{4, 13, 22, 31, 40, 49, 58, 67, 76},
	{5, 14, 23, 32, 41, 50, 59, 68, 77},
	{6, 15, 24, 33, 42, 51, 60, 69, 78},
	{7, 16, 25, 34, 43, 52, 61, 70, 79},
	{8, 17, 26, 35, 44, 53, 62, 71, 80}}

type Puzzle struct {
	name        string
	puzzleStart [81]uint8
	puzzleSol   [81]uint8
	score       float32
}

type PuzzleDataset struct {
	name       string
	puzzleSol  [81]uint8
	possibles  [81][]uint8
	score      float32
	iterations uint
}

func solvePuzzle(p *Puzzle) {
	// Create working dataset
	var pds PuzzleDataset
	pds.name = p.name
	pds.puzzleSol = p.puzzleStart
	pds.score = calculateCompletionScore(pds.puzzleSol)
	var tempscore = pds.score

	if globalDebugMode {
		// We won't calc init score during scoring run for speed
		fmt.Printf("Puzzle Name: %s\n", pds.name)
		fmt.Printf("Initial Score: %.2f%%\n", pds.score)
		printPuzzle(pds.puzzleSol)
	}

	// This is essentially a do while loop
	for ok := true; ok; ok = pds.score < 100 {
		// Iterate over all cells
		for i := 0; i < 81; i++ {

			// Consider running this before each algorithm
			updateAllCellPossibles(&pds)

			// Algorithm 1: Only one possibility so assign the first and only in slice
			if (pds.puzzleSol[i] == 0) && (len(pds.possibles[i]) == 1) {
				pds.puzzleSol[i] = pds.possibles[i][0]
			}

			// Algorithm 2 Check row for unique possibles
			if pds.puzzleSol[i] == 0 {
				relatedRowIdices := rowIndices[cellToRowLookup[i]]

				// Iterate over possibles of current cells
				for _, cellPossible := range pds.possibles[i] {

					var candidate bool = true
					// Compare against all other row possibilities
					for _, cellIndex := range relatedRowIdices {
						// Skip if the same
						if cellIndex == uint8(i) {
							continue
						}
						// First match we find move on to the next possible of this cell
						if contains(pds.possibles[cellIndex], cellPossible) {
							candidate = false
							continue
						}
					}
					if candidate {
						pds.puzzleSol[i] = cellPossible
					}
				}
			}

			// Algorithm 3 Check column for unique possibles
			if pds.puzzleSol[i] == 0 {
				relatedColumnIdices := colIndices[cellToColLookup[i]]

				// Iterate over possibles of current cells
				for _, cellPossible := range pds.possibles[i] {

					var candidate bool = true
					// Compare against all other column possibilities
					for _, cellIndex := range relatedColumnIdices {
						// Skip if the same
						if cellIndex == uint8(i) {
							continue
						}
						// First match we find move on to the next possible of this cell
						if contains(pds.possibles[cellIndex], cellPossible) {
							candidate = false
							continue
						}
					}
					if candidate {
						pds.puzzleSol[i] = cellPossible
					}
				}
			}

			// Algorithm 4 Check box for unique possibles
			if pds.puzzleSol[i] == 0 {
				relatedBoxIdices := boxIndices[cellToBoxLookup[i]]

				// Iterate over possibles of current cells
				for _, cellPossible := range pds.possibles[i] {

					var candidate bool = true
					// Compare against all other column possibilities
					for _, cellIndex := range relatedBoxIdices {
						// Skip if the same
						if cellIndex == uint8(i) {
							continue
						}
						// First match we find move on to the next possible of this cell
						if contains(pds.possibles[cellIndex], cellPossible) {
							candidate = false
							continue
						}
					}
					if candidate {
						pds.puzzleSol[i] = cellPossible
					}
				}
			}
		}

		pds.score = calculateCompletionScore(pds.puzzleSol)

		// If we hit the same score as last time, we are stuck
		if pds.score == tempscore {
			fmt.Println("STUCK!!!!!!!")
			break
		} else {
			tempscore = pds.score
		}
	}

	if globalDebugMode {
		fmt.Printf("Final Score: %.2f%%\n", pds.score)
		printPuzzle(pds.puzzleSol)

		// for i, v := range pds.possibles {
		// 	fmt.Println("idx: ", i, "pos: ", v)
		// }
	}

	p.puzzleSol = pds.puzzleSol
	p.score = pds.score
}

func updateAllCellPossibles(pds *PuzzleDataset) {
	for i := range pds.possibles {
		cellPossibles := []uint8{}

		// No work required if solution in place, i.e. non zero value of 1-9
		if pds.puzzleSol[i] != 0 {
			pds.possibles[i] = cellPossibles
			continue
		}

		relatedIndices := getCellRelatedIndices(uint8(i))
		possiblesKeepMap := map[uint8]bool{
			1: true,
			2: true,
			3: true,
			4: true,
			5: true,
			6: true,
			7: true,
			8: true,
			9: true,
		}

		for _, value := range relatedIndices {
			possiblesKeepMap[pds.puzzleSol[value]] = false
		}

		for key, value := range possiblesKeepMap {
			if value {
				cellPossibles = append(cellPossibles, key)
			}
		}

		pds.possibles[i] = cellPossibles
	}
}

func getCellRelatedIndices(cellIdx uint8) []uint8 {
	// Get all the related indices
	relatedIndices := []uint8{}

	// Row indices will be first and therefore unique, no need to check for duplicates (9 total)
	relatedIndices = append(relatedIndices, rowIndices[cellToRowLookup[cellIdx]]...)

	// Add unique column indices (8 more)
	for _, v := range colIndices[cellToColLookup[cellIdx]] {
		if !contains(relatedIndices, v) {
			relatedIndices = append(relatedIndices, v)
		}
	}

	// Add unique box indices (4 more)
	for _, v := range boxIndices[cellToBoxLookup[cellIdx]] {
		if !contains(relatedIndices, v) {
			relatedIndices = append(relatedIndices, v)
		}
	}

	return relatedIndices
}

func contains(s []uint8, e uint8) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func printPuzzle(p [81]uint8) {
	for j := 0; j < 9; j++ {
		if j%3 == 0 {
			fmt.Println("-------------------------")
		}
		for k := 0; k < 3; k++ {
			var a string = "_"
			var b string = "_"
			var c string = "_"
			if p[j*9+k*3] != 0 {
				a = strconv.FormatUint(uint64(p[j*9+k*3]), 10)
			}
			if p[j*9+k*3+1] != 0 {
				b = strconv.FormatUint(uint64(p[j*9+k*3+1]), 10)
			}
			if p[j*9+k*3+2] != 0 {
				c = strconv.FormatUint(uint64(p[j*9+k*3+2]), 10)
			}

			fmt.Printf("| %s %s %s ", a, b, c)
		}
		fmt.Printf("|\n")
	}
	fmt.Println("-------------------------")
}

// func printDetailedCore(p PuzzleDataset) {
// 	var possibles
// }

func calculateCompletionScore(p [81]uint8) float32 {
	var count uint8 = 0
	for _, value := range p {
		if value != 0 {
			count++
		}
	}
	return float32(count) / 81.0 * 100.0
}
