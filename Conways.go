package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type lifeGroundState struct {
	groundMatrixNow, groundMatrixNext [][]bool
	weight, height int
}

// set initial lives
func initialGound(weight, height int, modeSelect string) *lifeGroundState {
	// initial the state now and next to matrix at the size of the ground as [y][x]
	groundMatrixNow := make([][]bool, height)
	for i := range groundMatrixNow {
		groundMatrixNow[i] = make([]bool, weight)
	}

	groundMatrixNext := make([][]bool, height)
	for i := range groundMatrixNext {
		groundMatrixNext[i] = make([]bool, weight)
	}

	// begin as selected mode
	switch modeSelect {
	case "R":
		rand.Seed(time.Now().Unix())
		for i := 0; i < (weight * height / 6); i++ {	// adjust the initial number of lives
			groundMatrixNow[rand.Intn(height)][rand.Intn(weight)] = true
		}
	case "S":
		// Block
		groundMatrixNow[1][1] = true
		groundMatrixNow[1][2] = true
		groundMatrixNow[2][1] = true
		groundMatrixNow[2][2] = true

		// Bee-hive
		groundMatrixNow[1][7] = true
		groundMatrixNow[1][8] = true
		groundMatrixNow[2][6] = true
		groundMatrixNow[2][9] = true
		groundMatrixNow[3][7] = true
		groundMatrixNow[3][8] = true

		// Loaf
		groundMatrixNow[1][14] = true
		groundMatrixNow[1][15] = true
		groundMatrixNow[2][13] = true
		groundMatrixNow[2][16] = true
		groundMatrixNow[3][14] = true
		groundMatrixNow[3][16] = true
		groundMatrixNow[4][15] = true

		// Boat
		groundMatrixNow[1][20] = true
		groundMatrixNow[1][21] = true
		groundMatrixNow[2][20] = true
		groundMatrixNow[2][22] = true
		groundMatrixNow[3][21] = true

		// Tub
		groundMatrixNow[1][26] = true
		groundMatrixNow[2][25] = true
		groundMatrixNow[2][27] = true
		groundMatrixNow[3][26] = true

		// Blinker
		groundMatrixNow[2][31] = true
		groundMatrixNow[2][32] = true
		groundMatrixNow[2][33] = true

		// Toad
		groundMatrixNow[3][38] = true
		groundMatrixNow[3][39] = true
		groundMatrixNow[3][40] = true
		groundMatrixNow[4][37] = true
		groundMatrixNow[4][38] = true
		groundMatrixNow[4][39] = true

		// Beacon
		groundMatrixNow[2][44] = true
		groundMatrixNow[2][45] = true
		groundMatrixNow[3][44] = true
		groundMatrixNow[4][47] = true
		groundMatrixNow[5][46] = true
		groundMatrixNow[5][47] = true
	case "G":
		groundMatrixNow[1][3] = true
		groundMatrixNow[2][1] = true
		groundMatrixNow[2][3] = true
		groundMatrixNow[3][2] = true
		groundMatrixNow[3][3] = true
	case "M":
		// Light- weight spaceship (LWSS)
		groundMatrixNow[1][1] = true
		groundMatrixNow[1][4] = true
		groundMatrixNow[2][5] = true
		groundMatrixNow[3][1] = true
		groundMatrixNow[3][5] = true
		groundMatrixNow[4][2] = true
		groundMatrixNow[4][3] = true
		groundMatrixNow[4][4] = true
		groundMatrixNow[4][5] = true

		// Middle- weight spaceship (MWSS)
		groundMatrixNow[14][4] = true
		groundMatrixNow[15][2] = true
		groundMatrixNow[15][6] = true
		groundMatrixNow[16][1] = true
		groundMatrixNow[17][1] = true
		groundMatrixNow[17][6] = true
		groundMatrixNow[18][1] = true
		groundMatrixNow[18][2] = true
		groundMatrixNow[18][3] = true
		groundMatrixNow[18][4] = true
		groundMatrixNow[18][5] = true
	default:
		fmt.Println("Wrong Input, Switch to Random Mode in 5 sec")
		time.Sleep(time.Second * 5)
		rand.Seed(time.Now().Unix())
		for i := 0; i < (weight * height / 6); i++ {
			groundMatrixNow[rand.Intn(height)][rand.Intn(weight)] = true
		}
	}

	return &lifeGroundState {groundMatrixNow: groundMatrixNow, groundMatrixNext: groundMatrixNext, weight: weight, height: height}
}

// using the state now to calculate next state
func (lGS *lifeGroundState) nextState() {
	for y := 0; y < lGS.height; y ++ {	// go through each point
		for x := 0; x < lGS.weight; x ++ {
			roundLives := 0
			for i := -1; i <= 1; i ++ {	// go through all points around
				for j := -1; j <= 1; j ++ {
					if i == 0 && j == 0 {
						continue
					}
					if lGS.groundMatrixNow[(y + j + lGS.height) % lGS.height][(x + i + lGS.weight) % lGS.weight] == true {	// due the ground as a circle
						roundLives ++	// calculate the lives around
					}
				}
			}
			if lGS.groundMatrixNow[y][x] == true {	// calculate the lives of next state by rules
				if roundLives < 2 || roundLives > 3{
					lGS.groundMatrixNext[y][x] = false
				} else {
					lGS.groundMatrixNext[y][x] = true
				}
			} else {
				if roundLives == 3 {
					lGS.groundMatrixNext[y][x] = true
				} else {
					lGS.groundMatrixNext[y][x] = false
				}
			}
		}
	}

	// update the state
	for i := 0; i < lGS.height; i++ {
		for j := 0; j < lGS.weight; j++ {
			lGS.groundMatrixNow[i][j] = lGS.groundMatrixNext[i][j]
		}
	}
}

func main() {
	// let user choose the initial mode
	var modeSelect string
	fmt.Println("You should End the Game Manually")
	fmt.Println("Input Mode:")
	fmt.Println("Random(R) / Still lifes & Oscillators(S) / Glider(G) / Other Spaceships(M)")
	fmt.Scan(&modeSelect)

	// set the size of the ground and initial it
	weight := 50
	height := 25
	play := initialGound(weight, height, modeSelect)

	// transfer the state to string to print
	var show bytes.Buffer
	for ;; {	// user should end the game manually
		for y := 0; y < height; y ++ {
			for x := 0; x < weight; x ++ {
				switch play.groundMatrixNow[y][x] {
				case true:
					show.WriteByte('*')
				case false:
					show.WriteByte(' ')
				}
			}
			if y < height - 1 {
				show.WriteByte('\n')
			}
		}

		// clean the terminal to show next state
		cmdClear := exec.Command("clear")
		switch runtime.GOOS {
		case "windows":
			cmdClear = exec.Command("cmd", "/c", "cls")
		default:
			cmdClear = exec.Command("clear")
		}
		cmdClear.Stdout = os.Stdout
		cmdClear.Run()

		// show the state now for some time before update
		fmt.Println(show.String())
		show.Reset()
		time.Sleep(time.Second / 5)
		play.nextState()
	}
}
