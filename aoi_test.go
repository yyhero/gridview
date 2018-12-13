package gridview

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var (
	Debug = false
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randCoord(min, max int) int {
	return min + rand.Intn(max-min)
}

func TestGrid(t *testing.T) {
	mapGrids := NewGridMgr(0, 4000, 0, 4000, 20, 20)

	// 初始化20000个player
	creatures := make(map[string]*creature, 200)
	for i := 1; i <= 20000; i++ {
		posX := rand.Intn(4000)
		posY := rand.Intn(4000)
		pid := strconv.Itoa(i)
		obj := newCreature(posX, posY, pid)
		creatures[pid] = obj
		mapGrids.Enter(posX, posY, pid)
	}

	// all creatures move one time
	t0 := time.Now()
	for _, obj := range creatures {
		oldX, oldY := obj.posX, obj.posY
		for {
			posX, posY := oldX+randCoord(-5, 5), oldY+randCoord(-5, 5)
			if mapGrids.ValidPos(posX, posY) {
				obj.posX, obj.posY = posX, posY
				break
			}
		}
		mapGrids.Move(oldX, oldY, obj.posX, obj.posY, obj.id)
	}
	sub := time.Now().Sub(t0)
	fmt.Printf("all creatures move one time used %v !\n", sub)

	// one creature move 20000 times
	t0 = time.Now()
	obj := creatures[strconv.Itoa(1+rand.Intn(20000))]
	for i := 0; i < 20000; i++ {
		oldX, oldY := obj.posX, obj.posY
		for {
			posX, posY := oldX+randCoord(-10, 10), oldY+randCoord(-10, 10)
			if mapGrids.ValidPos(posX, posY) {
				obj.posX, obj.posY = posX, posY
				break
			}
		}

		lostGrids, bornGrids, aoiGrids, _ := mapGrids.Move(oldX, oldY, obj.posX, obj.posY, obj.id)
		printInfo(lostGrids, bornGrids, aoiGrids, mapGrids, obj)
	}
	sub = time.Now().Sub(t0)
	fmt.Printf("one creature move 20000 times used %v !\n", sub)
}

func printInfo(lostGrids, bornGrids, aoiGrids map[int]*grid, mapMgr *gridMgr, obj *creature) {
	if !Debug {
		return
	}

	printInfo := func(grids map[int]*grid) {
		for _, grid := range grids {
			fmt.Printf("grid id:%v,row:%v,col:%v, range:(minX:%v,maxX:%v,minY:%v,maxY:%v)\n", grid.gridId, grid.row, grid.col, grid.minX, grid.maxX, grid.minY, grid.maxY)
		}
	}
	grid := mapMgr.GetGridByPos(obj.posX, obj.posY)

	fmt.Printf("x:%v, y:%v,lostGrids len:%v\n", obj.posX, obj.posY, len(lostGrids))
	printInfo(lostGrids)
	print("---------------------------------------\n")

	fmt.Printf("x:%v, y:%v,bornGrids len:%v\n", obj.posX, obj.posY, len(bornGrids))
	printInfo(bornGrids)
	print("---------------------------------------\n")

	fmt.Printf("row:%v,col:%v, aoiGrids len:%v\n", grid.row, grid.col, len(aoiGrids))
	printInfo(aoiGrids)
	print("****************************************************************\n")
}
