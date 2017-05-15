package main

import (
	"fmt"
	"time"
	"encoding/json"
)

type World struct {
	Width  int
	Height int
	Matrix [][]int
}

func (w World) RunGeneration() ([]byte, time.Time, error){
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

	data, err := w.getJson()
	return data, time.Now(), err
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

func (w World) readFromJson(data []byte) (error) {
	if err := json.Unmarshal(data, &world); err != nil {
		fmt.Println("I failed at unmarhsalling the json ", err)
		return err
	}
	return nil
}

func (w World) getJson() ([]byte, error) {

	data, err := json.Marshal(&world)
	if err != nil {
		fmt.Println("I failed at marhsalling the json ", err)
	}
	return data, err

}



