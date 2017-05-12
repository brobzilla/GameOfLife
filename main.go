package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
)

type World struct {
	Width  int
	Height int
	Matrix [][]int
}

func (w World) RunGeneration() {
	tmp := make([][]int, w.Width)
	for i := range w.Matrix {
		tmp[i] = make([]int, w.Height)
		for j := range w.Matrix[i] {
			num_neighbors := w.calculateNeighbors(i,j)
			if w.Matrix[i][j] == 0 { // dead Case
				if num_neighbors == 3 {
					tmp[i][j] = 1
				} else {
					tmp[i][j] = 0
				}
			} else { // live case
				if num_neighbors < 2 {
					tmp[i][j] = 0 //Dies underpopulation
				} else if num_neighbors >=2 && num_neighbors <= 3 {
					tmp[i][j] = 1  // Lives to see another day.
				} else {
					tmp[i][j] = 0 // Dies overpopulation
				}
			}
		}
	}

	// Copy tmp world over to master.
	for i := range w.Matrix {
		for j := range w.Matrix[i] {
			w.Matrix[i][j] = tmp[i][j]
		}
	}
}

func (w World) calculateNeighbors(x int, y int) (int) {
	drow := [...]int{ 0,  1, 1, 1, 0, -1, -1, -1}
	dcol := [...]int{-1, -1, 0, 1, 1,  1,  0, -1}
	var cnt int
	for k := 0; k < 8; k++ {
		if x + drow[k] < 0 || y + dcol[k] < 0 || x + drow[k] >= w.Width || y + dcol[k] >= w.Height {
			continue
		}
		if w.Matrix[x+drow[k]][y+dcol[k]] == 1 {
			cnt++
		}

	}
	return cnt
}

func (w World) Print() {
	fmt.Println("")
	for i := range w.Matrix {
		for j := range w.Matrix[i] {
			if w.Matrix[i][j] == 1 {
				fmt.Print(" X ")
			} else {
				fmt.Print("   ")
			}

		}
		fmt.Println("|")
	}
	fmt.Println("")

}

func main() {
	f, err := os.OpenFile("gameoflife.log", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error case ", err)
		os.Exit(1)
	}

	log.Println(string(data))

	var world World

	if err := json.Unmarshal(data, &world); err != nil {
		fmt.Println("I failed at unmarhsalling the json ", err)
	}

	log.Println("Width = ", world.Width)
	log.Println("Height = ", world.Height)


	for {
		world.Print()
		world.RunGeneration()
		time.Sleep(100 * time.Millisecond)
	}
}
