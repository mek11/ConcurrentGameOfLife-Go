package life

import "sync"

// Posizioni per la scelta del canale da utilizzare
const (
	top         = 0
	topRight    = 1
	right       = 2
	bottomRight = 3
	bottom      = 4
	bottomLeft  = 5
	left        = 6
	topLeft     = 7
)

// used to communicate values for ghost cells
type ghostCell struct {
	pos, index, val int
}

// step of the game of a portion of the board with the previous and current state
type part struct {
	rows, cols       int
	xCorner, yCorner int
	prev, curr       *grid
	myCh             <-chan ghostCell
	othCh            []chan<- ghostCell
	chDraw           chan<- pixel
}

func newPart(height, width, xCorner, yCorner int, myCh <-chan ghostCell, chDraw chan<- pixel, pred predicate, fun increase) *part {
	return &part{
		rows:    height,
		cols:    width,
		xCorner: xCorner,
		yCorner: yCorner,
		prev:    newGrid(height, width, pred, fun),
		curr:    newGrid(height, width, pred, fun),
		myCh:    myCh,
		chDraw:  chDraw}
}

// advances the part by one step
func (p *part) advance(wg *sync.WaitGroup) {
	// swap current grid
	p.prev, p.curr = p.curr, p.prev
	// update the grid
	for r := 0; r < p.rows; r++ {
		for c := 0; c < p.cols; c++ {
			newCell := p.prev.next(r, c)
			p.set(r, c, newCell)
			p.chDraw <- pixel{c + p.xCorner, r + p.yCorner, newCell}
		}
	}
	// update ghost cells
	p.updateGhostCells(wg)
}

func (p *part) updateGhostCells(wg *sync.WaitGroup) {
	go p.sendGhostCells()
	p.receiveGhostCells()
	wg.Done()
}

func (p *part) sendGhostCells() {
	// first and last row
	for i := 0; i < p.cols; i++ {
		p.othCh[top] <- ghostCell{bottom, i, p.at(0, i)}
		p.othCh[bottom] <- ghostCell{top, i, p.at(p.rows-1, i)}
	}
	// first and last column
	for i := 0; i < p.rows; i++ {
		p.othCh[left] <- ghostCell{right, i, p.at(i, 0)}
		p.othCh[right] <- ghostCell{left, i, p.at(i, p.cols-1)}
	}
	// corners
	p.othCh[bottomRight] <- ghostCell{topLeft, 0, p.at(p.rows-1, p.cols-1)}
	p.othCh[bottomLeft] <- ghostCell{topRight, 0, p.at(p.rows-1, 0)}
	p.othCh[topLeft] <- ghostCell{bottomRight, 0, p.at(0, 0)}
	p.othCh[topRight] <- ghostCell{bottomLeft, 0, p.at(0, p.cols-1)}
}

func (p *part) receiveGhostCells() {
	for i := 0; i < 2*p.rows+2*p.cols+4; i++ {
		gc := <-p.myCh
		switch gc.pos {
		case top:
			p.set(-1, gc.index, gc.val)
		case bottom:
			p.set(p.rows, gc.index, gc.val)
		case left:
			p.set(gc.index, -1, gc.val)
		case right:
			p.set(gc.index, p.cols, gc.val)
		case topLeft:
			p.set(-1, -1, gc.val)
		case topRight:
			p.set(-1, p.cols, gc.val)
		case bottomRight:
			p.set(p.rows, p.cols, gc.val)
		case bottomLeft:
			p.set(p.rows, -1, gc.val)
		}
	}
}

func (p *part) at(r, c int) int {
	return p.curr.at(r, c)
}

func (p *part) set(r, c int, val int) {
	p.curr.set(r, c, val)
}

func (p *part) drawPart(wg *sync.WaitGroup) {
	for r := 0; r < p.rows; r++ {
		for c := 0; c < p.cols; c++ {
			p.chDraw <- pixel{c + p.xCorner, r + p.yCorner, p.at(r, c)}
		}
	}
	wg.Done()
}
