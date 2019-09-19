package main

import (
	"lifego/life"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 7 {
		panic("Expecting 6 arguments")
	}
	path := os.Args[1]
	h, err1 := strconv.Atoi(os.Args[2])
	w, err2 := strconv.Atoi(os.Args[3])
	numH, err3 := strconv.Atoi(os.Args[4])
	numW, err4 := strconv.Atoi(os.Args[5])
	iter, err5 := strconv.Atoi(os.Args[6])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		panic("Invalid arguments")
	}

	img := life.NewGifCreator(w*numW, h*numH)
	g, err := life.NewGame(path, h, w, numH, numW, iter, img)
	if err != nil {
		panic(err)
	}
	g.Start()

}
