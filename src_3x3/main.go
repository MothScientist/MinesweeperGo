package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"slices"
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
			printOriginalField(field)
			break
		}
	}

}

func fillSquare(minesCount int) [3][3]int {
	field := [3][3]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	var minesPositions [][2]int
	// Determine their location within a 3x3 field
	// If we have a NxN field, where column i and row j are the location of the first mine,
	// then the second has k, n, where k != i && n != j
	x1, y1 := rand.Intn(3), rand.Intn(3)
	minesPositions = append(minesPositions, [2]int{x1, y1})
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
		minesPositions = append(minesPositions, [2]int{x2, y2})
	}

	// Mine = -3
	// Closed cell without mine = (0, 2)
	// When cell is opened: 0 -> -10, 1 -> -1, 2 -> -2

	// First, indicate the location of the mines separately.
	for i := range minesPositions {
		field[minesPositions[i][0]][minesPositions[i][1]] = -3
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	// Now you can run goroutines in the same cycle as above
	for i := range minesPositions {
		wg.Add(1)
		go searchMines(&field, minesPositions[i], &wg, &mutex)
	}
	wg.Wait()


	return field
}

// Filling the fields with a number showing the number of mines in its perimeter
func searchMines(field *[3][3]int, pos [2]int, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			searchX := pos[0] + k
			searchY := pos[1] + l
			if searchX < 0 || searchX > 2 || searchY < 0 || searchY > 2 {
				continue
			} else if field[searchX][searchY] != -3 {
				// (about != -3) Don`t forget that in this condition we can stumble upon not only the current mine,
				// but also another one located within the perimeter of the current one
				mutex.Lock() // Before making changes, be sure to block access to shared data
				field[searchX][searchY]++
				mutex.Unlock() // Can't use defer because we are in a loop - it will lead to a deadlock
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
	if slices.Contains(openPoints, userPosition) {
		return "This point opened"
	}
	for _, val := range availablePoints {
		if userPosition == val {
			return ""
		}
	}
	return "Not available point"
}

// Transform the field after the user opens a new field without a mine
func changeField(field *[3][3]int, i int, j int) {
	if field[i][j] == 0 {
		field[i][j] = -10
	} else if field[i][j] > 0 {
		field[i][j] = -field[i][j]
	}
}

// Displays the location of mines on the map to the user at the end of the game
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
