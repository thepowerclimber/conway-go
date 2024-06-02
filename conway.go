package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

const ROWS int = 48
const COLS int = 168
const SLEEP time.Duration = 50
const MUTATION bool = true
const MUTATION_CHANCE float32 = 0.001
const PROFILE bool = false

var current_state [ROWS][COLS]bool
var next_state [ROWS][COLS]bool

func display() {
	for x := 0; x < ROWS; x++ {
		for y := 0; y < COLS; y++ {
			if current_state[x][y] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Printf("\033[%dA\033[%dD", ROWS, COLS)
}

func mod(a int, b int) int {
	return (a%b + b) % b
}

func count_neighbors(x int, y int) int {
	var count int = 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if !(dx == 0 && dy == 0) {
				var nx int = mod(x+dx, ROWS)
				var ny int = mod(y+dy, COLS)
				if current_state[nx][ny] {
					count++
				}
			}
		}
	}
	return count
}

func calculate_next_state() {
	for x := 0; x < ROWS; x++ {
		for y := 0; y < COLS; y++ {
			var num_neighbors int = count_neighbors(x, y)
			if current_state[x][y] {
				next_state[x][y] = num_neighbors == 2 || num_neighbors == 3
			} else {
				next_state[x][y] = num_neighbors == 3
			}
			if MUTATION {
				var res int = rand.IntN(100000)
				if res < int(MUTATION_CHANCE*1000) {
					next_state[x][y] = true
				}
			}
		}
	}
}

func set_glider_state() {
	//   0 1 2
	// 0 . # .
	// 1 . . #
	// 2 # # #
	current_state[0][1] = true
	current_state[1][2] = true
	current_state[2][0] = true
	current_state[2][1] = true
	current_state[2][2] = true
}

func set_random_state() {
	for x := 0; x < ROWS; x++ {
		for y := 0; y < COLS; y++ {
			var res int = rand.IntN(100)
			if res > 80 {
				current_state[x][y] = true
			}
		}
	}
}

func cleanup() {
	fmt.Println("Keyboard interrupted...")
	if PROFILE {
		pprof.StopCPUProfile()
	}
}

func main() {
	if PROFILE {
		// Start profiling
		f, err := os.Create("conway.prof")
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.StartCPUProfile(f)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	set_glider_state()
	set_random_state()
	for {
		display()
		calculate_next_state()
		current_state = next_state
		time.Sleep(SLEEP * time.Millisecond)
	}
}
