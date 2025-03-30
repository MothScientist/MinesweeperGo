package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

func main() {
	// Determine the number of mines (1-2)
	minesCount := rand.Intn(2) + 1
	// Create a field
	field := fillSquare(minesCount)
	var userPosition int
	var openPoints []int
	var availablePoints = [9]int{11, 12, 13, 21, 22, 23, 31, 32, 33}
	userMsg := "" // Since the console is cleared at the beginning of each cycle, the warning is displayed at the beginning of the next cycle
	for {
		clearConsole()
		if userMsg != "" {
			fmt.Println(userMsg)
			userMsg = ""
		}
		printField(field)
		fmt.Println("Enter position in format XY, example 11/23/31 etc.")
		_, err := fmt.Fscan(os.Stdin, &userPosition)
		if err != nil {
			return
		}

		warning := validate(availablePoints, userPosition, openPoints)
		if warning != "" {
			userMsg = warning
			continue
		} else {
			openPoints = append(openPoints, userPosition)
		}

		i, j := userPosition/10-1, userPosition%10-1
		if field[i][j] == -3 {
			clearConsole()
			fmt.Println("BOOM!")
			printOriginalField(field)
			break
		} else {
			changeField(&field, i, j)
		}

		if len(openPoints) == len(availablePoints)-minesCount {
			clearConsole()
			fmt.Println("Congratulation!")
			break
		}
	}

}

func fillSquare(minesCount int) [3][3]int {
	field := [3][3]int{}
	// Determine their location within a 3x3 field
	// If we have a NxN field, where column i and row j are the location of the first mine,
	// then the second has k, n, where k != i && n != j
	x1, y1 := rand.Intn(3), rand.Intn(3)
	var x2 int
	var y2 int
	if minesCount == 2 {
		// Go to the left so that when we get into the negative area, we invert it into a positive one.
		// That is, we consider going beyond the left border as an "entrance" to the field from the right border
		x2, y2 = x1-(rand.Intn(2)+1), y1-(rand.Intn(2)+1)
		if x2 < 0 {
			x2 += 3
		}
		if y2 < 0 {
			y2 += 3
		}
	}
	// Mine = -3
	// Closed cell without mine = (0, 2)
	// When cell is opened: 0 -> -10, 1 -> -1, 2 -> -2

	// We will fill it line by line, each nested array is 1 line of the field, starting from the bottom
	for i := range field {
		for j := range field[i] {
			if (i == x1 && j == y1) || (minesCount == 2 && (i == x2 && j == y2)) {
				field[i][j] = -3
			} else {
				field[i][j] = 0 // Pre-fill with zeros, because in the current cycle we do not know the exact location of the mines
			}
		}
	}

	// Now for each element we find the number of mines in its perimeter
	wg := sync.WaitGroup{}
	for i := range field {
		for j := range field[i] {
			if field[i][j] != -3 {
				wg.Add(1)
				go searchMines(&field, i, j, &wg)
			}
		}
	}
	wg.Wait()

	return field
}

func searchMines(field *[3][3]int, i int, j int, wg *sync.WaitGroup) {
	defer wg.Done()
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			searchX := i + k
			searchY := j + l
			if searchX < 0 || searchX > 2 || searchY < 0 || searchY > 2 {
				continue
			} else if field[searchX][searchY] == -3 {
				field[i][j] += 1
			}
		}
	}
}

// Console clearing function
func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return
	}
}

func printField(field [3][3]int) {
	var val string
	for i := range len(field) {
		for j := range field[i] {
			if field[i][j] < 0 && field[i][j] != -3 {
				if field[i][j] == -10 {
					field[i][j] = 0
				} else {
					field[i][j] = -field[i][j]
				}
				val = strconv.Itoa(field[i][j]) // int -> string
			} else {
				val = "?"
			}
			fmt.Print(val, " ")
		}
		fmt.Println()
	}
}

func validate(availablePoints [9]int, userPosition int, openPoints []int) string {
	for _, val := range openPoints {
		if userPosition == val {
			return "This point opened"
		}
	}
	for _, val := range availablePoints {
		if userPosition == val {
			return ""
		}
	}
	return "Not available point"
}

func changeField(field *[3][3]int, i int, j int) {
	if field[i][j] == 0 {
		field[i][j] = -10
	} else if field[i][j] > 0 {
		field[i][j] = -field[i][j]
	}
}

func printOriginalField(field [3][3]int) {
	for i := range field {
		for j := range field[i] {
			if field[i][j] == -3 {
				fmt.Print("*", " ")
			} else {
				fmt.Print("?", " ")
			}
		}
		fmt.Println()
	}
}
