package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

var up = rune("u"[0])
var down = rune("down"[0])
var right = rune("right"[0])
var left = rune("left"[0])

var reverser = map[rune]rune{
	up:    down,
	down:  up,
	right: left,
	left:  right,
}

type Config struct {
	startX int
	startY int
	endX   int
	endY   int
}

type Block struct {
	X             int
	Y             int
	DirectionFrom rune
	MustTurn      int
}

type BlockMap struct {
	BlocksCreated map[Block]int
	BlocksPath    *PriorityQueue
}

type Paths struct {
	CurrentBlock Block
	Value        int
}

type PriorityQueue []*Paths

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Value < pq[j].Value
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Paths)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func main() {

	inputFile, err := os.Open("seventeen.txt")
	if err != nil {
		fmt.Println("ABORT")
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	var lineString string
	var grid [][]int

	for scanner.Scan() {
		var line []int
		lineString = "0" + (scanner.Text()) + "0"
		for _, runeDigit := range lineString {
			digit, err := strconv.Atoi(string(runeDigit))
			if err != nil {
				fmt.Println("fail")
				os.Exit(1)
			}
			line = append(line, digit)
		}
		grid = append(grid, line)
	}

	config := Config{startX: 1, startY: 1, endX: len(grid[0]) - 2, endY: len(grid)}

	var emptyline = make([]int, len(grid[0]))

	grid = append([][]int{emptyline}, grid...)
	grid = append(grid, emptyline)

	priorityQueue := make(PriorityQueue, 0)
	heap.Init(&priorityQueue)

	blockMap := BlockMap{BlocksCreated: make(map[Block]int), BlocksPath: &priorityQueue}

	minValue := constructPaths(grid, blockMap, config)

	fmt.Println(minValue)
}

func constructPaths(grid [][]int, blockMap BlockMap, config Config) int {
	startBlock1 := Block{
		X:             config.startX,
		Y:             config.startY,
		DirectionFrom: right,
	}

	blockMap.BlocksCreated[startBlock1] = 0

	path1 := Paths{
		CurrentBlock: startBlock1,
	}
	heap.Push(blockMap.BlocksPath, &path1)

	for blockMap.BlocksPath.Len() > 0 {

		currentPath := heap.Pop(blockMap.BlocksPath).(*Paths)
		currentValue := currentPath.Value

		if currentPath.CurrentBlock.X == config.endX && currentPath.CurrentBlock.Y == config.endY {
			return currentValue
		}

		var nextX, nextY, turning int
		// right
		for _, directionFrom := range []rune{right, left, down, up} {
			turning = currentPath.CurrentBlock.MustTurn

			// check if reverse direction
			if currentPath.CurrentBlock.DirectionFrom == reverser[directionFrom] {
				continue
			}

			// check if must turn
			if currentPath.CurrentBlock.DirectionFrom == directionFrom && turning == 3 {
				continue
			} else if currentPath.CurrentBlock.DirectionFrom == directionFrom {
				turning++
			} else {
				turning = 1
			}

			if directionFrom == right {
				nextX = currentPath.CurrentBlock.X + 1
				nextY = currentPath.CurrentBlock.Y
			}
			if directionFrom == left {
				nextX = currentPath.CurrentBlock.X - 1
				nextY = currentPath.CurrentBlock.Y
			}
			if directionFrom == down {
				nextX = currentPath.CurrentBlock.X
				nextY = currentPath.CurrentBlock.Y + 1
			}
			if directionFrom == up {
				nextX = currentPath.CurrentBlock.X
				nextY = currentPath.CurrentBlock.Y - 1
			}

			// check if outside the grid
			if grid[nextY][nextX] == 0 {
				continue
			}

			var newBlockValue int
			newPath := Paths{CurrentBlock: Block{}, Value: grid[nextY][nextX] + currentValue}
			newPath.CurrentBlock.X = nextX
			newPath.CurrentBlock.Y = nextY
			newPath.CurrentBlock.DirectionFrom = directionFrom
			newPath.CurrentBlock.MustTurn = turning
			if _, ok := blockMap.BlocksCreated[newPath.CurrentBlock]; ok {
				newBlockValue = blockMap.BlocksCreated[newPath.CurrentBlock]
				if newBlockValue >= newPath.Value {
					blockMap.BlocksCreated[newPath.CurrentBlock] = newPath.Value
				} else {
					continue
				}
			} else {
				blockMap.BlocksCreated[newPath.CurrentBlock] = newPath.Value
			}
			heap.Push(blockMap.BlocksPath, &newPath)
		}
	}
	return 0
}
