package life

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Game represent an istance of a game
type Game struct {
	width, height int        // of a part
	numW, numH    int        // number of parts (horiz. and vert.)
	parts         [][]*part  // parts
	chDraw        chan pixel // channel to send pixels to draw
	iter          int        // number of iterations
	img           ImgCreator
	pred          predicate
	fun           increase
}

// NewGame creates a new istance of a game
func NewGame(path string, height, width, numH, numW, iter int, img ImgCreator) (*Game, error) {

	totHeight, totWidth := height*numH, width*numW

	g := new(Game)
	g.width, g.height, g.numW, g.numH = width, height, numW, numH
	g.iter = iter
	g.img = img

	//functions to use for the computation
	g.pred = func(c int) bool { return c == 1 }
	g.fun = func(i int, c int) int {
		if i == 3 || i == 2 && c == 1 {
			return 1
		}
		if c == 1 {
			return 2
		}
		return 0
	}

	//creation of channels
	chs := make([][]chan ghostCell, numH)
	for i := range chs {
		chs[i] = make([]chan ghostCell, numW)
		for j := range chs[i] {
			chs[i][j] = make(chan ghostCell, 2*height+2*width+4)
		}
	}
	g.chDraw = make(chan pixel, totHeight*totWidth)

	//creation of parts
	g.parts = make([][]*part, numH)
	for i := range g.parts {
		g.parts[i] = make([]*part, numW)
	}
	for i := range g.parts {
		for j := range g.parts[i] {
			g.parts[i][j] = newPart(height, width, j*width, i*height, chs[i][j], g.chDraw, g.pred, g.fun)
		}
	}

	//create links between adjacent parts
	linksParts(g.parts, chs)

	//set values
	if path == "random" {
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < (totWidth * totHeight / 4); i++ {
			r, c, v := rand.Intn(numH*height), rand.Intn(numW*width), 1
			g.parts[r/height][c/width].set(r%height, c%width, v)
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		scan := bufio.NewScanner(f)

		for scan.Scan() {
			values := strings.Fields(scan.Text())
			if len(values) != 3 {
				return nil, errors.New("Invalid file content")
			}
			r, err1 := strconv.Atoi(values[0])
			c, err2 := strconv.Atoi(values[1])
			v, err3 := strconv.Atoi(values[2])
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, errors.New("Invalid file content")
			}
			if r >= height*numH || c >= width*numW {
				return nil, errors.New("Dimensions too small")
			}
			g.parts[r/height][c/width].set(r%height, c%width, v)
		}
	}
	return g, nil
}

func linksParts(p [][]*part, chs [][]chan ghostCell) {
	h := len(p)
	w := len(p[0])

	// creation of links to (8) channels for each part
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			p[i][j].othCh = make([]chan<- ghostCell, 8)

			p[i][j].othCh[top] = chs[(i-1+h)%h][j]
			p[i][j].othCh[topRight] = chs[(i-1+h)%h][(j+1)%w]
			p[i][j].othCh[right] = chs[i][(j+1)%w]
			p[i][j].othCh[bottomRight] = chs[(i+1)%h][(j+1)%w]
			p[i][j].othCh[bottom] = chs[(i+1)%h][j]
			p[i][j].othCh[bottomLeft] = chs[(i+1)%h][(j-1+w)%w]
			p[i][j].othCh[left] = chs[i][(j-1+w)%w]
			p[i][j].othCh[topLeft] = chs[(i-1+h)%h][(j-1+w)%w]
		}
	}
}

// Start the game
func (g *Game) Start() {

	chDone := make(chan struct{})
	go g.img.create(g.chDraw, chDone)

	g.initialize()
	for it := 0; it < g.iter; it++ {
		g.advance()
	}

	close(g.chDraw)
	<-chDone // attendi fine di create()
}

// sets up the game before starting
func (g *Game) initialize() {
	var wg sync.WaitGroup
	for _, ps := range g.parts {
		for _, p := range ps {
			wg.Add(1)
			go p.drawPart(&wg)
			wg.Add(1)
			go p.updateGhostCells(&wg)
		}
	}
	wg.Wait()
}

// advances the game by one step
func (g *Game) advance() {
	var wg sync.WaitGroup
	for _, ps := range g.parts {
		for _, p := range ps {
			wg.Add(1)
			go p.advance(&wg)
		}
	}
	wg.Wait()
}
