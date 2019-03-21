package main

type Matrix []Row

type Row []int

func initMatrix() *Matrix {
	// Creates a 20 by 10 empty Matrix
	newMatrix := &Matrix{}
	for i := 0; i < 20; i++ {
		// 10 wide row
		row := &Row{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		*newMatrix = append(*newMatrix, *row)
	}
	return newMatrix
}

func (m *Matrix) updateMatrix() {
}
