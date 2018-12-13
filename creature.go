package gridview

type creature struct {
	posX int
	posY int
	id   string
}

func newCreature(posX, posY int, id string) *creature {
	obj := &creature{
		posX: posX,
		posY: posY,
		id:   id,
	}
	return obj
}
