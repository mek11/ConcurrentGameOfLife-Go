package life

// function to determine whether to increase the counter in the computation of the next value of a cell
type predicate func(int) bool

// function to determine the next value of a cell given a counter and the current state
type increase func(int, int) int

// representation of the state of a (portion of) the game board
type grid struct {
	rows, cols int
	grid       []int
	pred       predicate
	fun        increase
}

func newGrid(rows, cols int, pred predicate, fun increase) *grid {
	// create a new grid
	g := make([]int, (rows+2)*(cols+2))
	return &grid{rows: rows, cols: cols, grid: g, pred: pred, fun: fun}
}

func (g *grid) set(r, c int, val int) {
	g.grid[(g.cols+2)*(r+1)+(c+1)] = val
}

func (g *grid) at(r, c int) int {
	return g.grid[(g.cols+2)*(r+1)+(c+1)]
}

// next value of the cell in the specified position
func (g *grid) next(r, c int) int {
	counter := 0
	for i := r - 1; i <= r+1; i++ {
		for j := c - 1; j <= c+1; j++ {
			if (i != r || j != c) && g.pred(g.at(i, j)) {
				counter++
			}
		}
	}
	return g.fun(counter, g.at(r, c))
}
