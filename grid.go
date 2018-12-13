package gridview

import "sync"

type grid struct {
	gridId int
	row    int
	col    int

	// 网格矩形坐标
	minX int
	maxX int
	minY int
	maxY int

	// 集合
	palyers map[string]bool
	mutex   sync.RWMutex
}

func (g *grid) AddPlayer(playerId string) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	g.palyers[playerId] = true
}

func (g *grid) DeletePlayer(playerId string) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	delete(g.palyers, playerId)
}

func NewGrid(_gridId int,
	_minX int,
	_maxX int,
	_minY int,
	_maxY int,
	_row int,
	_col int) *grid {
	obj := &grid{
		gridId:  _gridId,
		minX:    _minX,
		maxX:    _maxX,
		minY:    _minY,
		maxY:    _maxY,
		row:     _row,
		col:     _col,
		palyers: make(map[string]bool, 0),
	}
	return obj
}
