package gridview

import (
	"math"
	"sync"
)

type gridMgr struct {
	allGrids map[int]*grid
	mutex    sync.RWMutex

	// 地图大小
	minX int
	maxX int
	minY int
	maxY int

	// 行列
	row int
	col int

	// 格子大小
	lenX int
	lenY int
}

func (g *gridMgr) ValidPos(posX, posY int) bool {
	if posX <= g.minX || posX > g.maxX || posY <= g.minY || posY > g.maxY {
		return false
	}
	return true
}

func (g *gridMgr) Enter(posX, posY int, playerId string) map[int]*grid {
	grid := g.GetGridByPos(posX, posY)
	grid.AddPlayer(playerId)

	return g.GetInterestAreaByPos(posX, posY)
}

func (g *gridMgr) Leave(posX, posY int, playerId string) map[int]*grid {
	grid := g.GetGridByPos(posX, posY)
	grid.DeletePlayer(playerId)

	return g.GetInterestAreaByPos(posX, posY)
}

func (g *gridMgr) Move(oldPosX, oldPosY, posX, posY int, playerId string) (lostGrids, bornGrids, aoiGrids map[int]*grid, isCrossGrid bool) {
	oldGrid := g.GetGridByPos(oldPosX, oldPosY)
	curGrid := g.GetGridByPos(posX, posY)

	if oldGrid != curGrid {
		oldGrid.DeletePlayer(playerId)
		curGrid.AddPlayer(playerId)

		oldArea := g.GetInterestAreaByPos(oldPosX, oldPosY)
		curArea := g.GetInterestAreaByPos(posX, posY)

		lostGrids = make(map[int]*grid)
		bornGrids = make(map[int]*grid)
		for _, obj := range oldArea {
			if _, exist := curArea[obj.gridId]; !exist {
				lostGrids[obj.gridId] = obj
			}
		}
		for _, obj := range curArea {
			if _, exist := oldArea[obj.gridId]; !exist {
				bornGrids[obj.gridId] = obj
			}
		}
	}

	aoiGrids = g.GetInterestAreaByPos(posX, posY)
	return
}

func (g *gridMgr) AddAoiByGridId(gridId int, playerId string) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	grid := g.allGrids[gridId]
	grid.AddPlayer(playerId)
}

func (g *gridMgr) AddAoiByPos(posX, posY int, playerId string) {
	grid := g.GetGridByPos(posX, posY)
	grid.AddPlayer(playerId)
}

func (g *gridMgr) GetGridByPos(posX, posY int) *grid {
	row, col := g.GetRowColByPos(posX, posY)
	result := g.GetGridByRowCol(row, col)
	if result == nil {
		print("\n GetGridByPos grid is nil:", row, "---", col)
	}
	return result
}

func (g *gridMgr) GetGridByRowCol(row, col int) *grid {
	gridId := col + (row-1)*g.col
	return g.GetGridById(gridId)
}

func (g *gridMgr) GetGridById(gridId int) *grid {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.allGrids[gridId]
}

func (g *gridMgr) GetRowColByPos(posX, posY int) (int, int) {
	col := int(math.Ceil(float64(posX) / float64(g.lenX)))
	row := int(math.Ceil(float64(posY) / float64(g.lenY)))
	if col == 0 {
		col += 1
	}
	if row == 0 {
		row += 1
	}
	return row, col
}

func (g *gridMgr) GetInterestAreaByPos(posX, posY int) map[int]*grid {
	row, col := g.GetRowColByPos(posX, posY)

	curGrid := g.GetGridByRowCol(row, col)
	grids := make(map[int]*grid, 9)
	midGrids := make(map[int]*grid, 3)
	grids[curGrid.gridId] = curGrid
	midGrids[curGrid.gridId] = curGrid

	// 中间左边
	if col > 1 {
		temp := g.GetGridByRowCol(row, col-1)
		grids[temp.gridId] = temp
		midGrids[temp.gridId] = temp
	}
	// 中间右边
	if col < g.col {
		temp := g.GetGridByRowCol(row, col+1)
		grids[temp.gridId] = temp
		midGrids[temp.gridId] = temp
	}

	// 遍历中间，搜索上下行
	for _, grid := range midGrids {
		if grid.row > 1 {
			temp := g.GetGridByRowCol(grid.row-1, grid.col)
			grids[temp.gridId] = temp
		}
		if grid.row < g.row {
			temp := g.GetGridByRowCol(grid.row+1, grid.col)
			grids[temp.gridId] = temp
		}
	}
	return grids
}

func NewGridMgr(_minX int,
	_maxX int,
	_minY int,
	_maxY int,
	_lenX int,
	_lenY int,
) *gridMgr {
	obj := &gridMgr{
		minX:     _minX,
		maxX:     _maxX,
		minY:     _minY,
		maxY:     _maxY,
		lenX:     _lenX,
		lenY:     _lenY,
		allGrids: make(map[int]*grid),
	}

	// 初始化网格
	col := int(math.Ceil(float64((_maxX - _minX) / _lenX)))
	row := int(math.Ceil(float64((_maxY - _minY) / _lenY)))
	for i := 1; i <= col; i++ {
		for j := 0; j < row; j++ {
			id := i + j*col
			gridMinX := i*_lenX - _lenX
			gridMaxX := i * _lenX
			girdMinY := j * _lenY
			girdMaxY := j*_lenY + _lenY
			grid := NewGrid(id, gridMinX, gridMaxX, girdMinY, girdMaxY, j+1, i)
			obj.allGrids[id] = grid
		}
	}
	obj.col = col
	obj.row = row

	return obj
}
