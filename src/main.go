package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var sc = bufio.NewScanner(os.Stdin)
var buff []byte

func nextInt() int {
	sc.Scan()
	i, err := strconv.Atoi(sc.Text())
	if err != nil {
		panic(err)
	}
	return i
}
func init() {
	sc.Split(bufio.ScanWords)
	sc.Buffer(buff, bufio.MaxScanTokenSize*1024)
	log.SetFlags(log.Lshortfile)
}

// https://golang.org/pkg/runtime/pprof/
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...
	solver()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

var N int
var targetScore float64

var Tile [50][50]int16 // 0~2500
var GetPoint [50][50]int16
var TileList [2500][]Point // tileのNo.から他の座標を取る

type Point struct {
	y, x int16
}

var start Point

func solver() {
	// input
	start.y = int16(nextInt())
	start.x = int16(nextInt())
	var i, j int16
	for i = 0; i < 50; i++ {
		for j = 0; j < 50; j++ {
			Tile[i][j] = int16(nextInt())
			TileList[Tile[i][j]] = append(TileList[Tile[i][j]], Point{y: i, x: j})
		}
	}
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			GetPoint[i][j] = int16(nextInt())
		}
	}
	beamsearch(start)
}

type State struct {
	scoreReal int
	score     int
	stepCount int
	move      [2000]uint8
	position  Point
	footprint [50][50]bool
}

func (s *State) Add(p Point) {
	tileNum := Tile[p.y][p.x]
	for _, tp := range TileList[tileNum] {
		s.footprint[tp.y][tp.x] = true
	}
}

func (s State) canMove(p Point) bool {
	if p.x < 0 || p.x > 49 || p.y < 0 || p.y > 49 {
		return false
	}
	if s.footprint[p.y][p.x] {
		return false
	}
	return true
}

func (s *State) Move(p Point, d int) {
	s.scoreReal += int(GetPoint[p.y][p.x])
	s.score += int(GetPoint[p.y][p.x])
	s.score += absInt(int(start.x-p.x)) + absInt(int(start.y-p.y))
	s.position = p
	s.move[s.stepCount] = uint8(d)
	s.stepCount++
	s.Add(p)
}

var dy = []int16{0, 1, 0, -1}
var dx = []int16{1, 0, -1, 0}
var direction = []string{"R", "D", "L", "U"}

var timeout bool
var TimeLimit time.Duration = 1000

func beamsearch(start Point) {
	go func() {
		time.Sleep(TimeLimit * time.Millisecond)
		timeout = true
	}()
	var startState Item
	startState.value.position = start
	startState.value.Add(startState.value.position)
	states := make([]PriorityQueue, 2001)
	for i := 0; i < 2000; i++ {
		states[i] = make(PriorityQueue, 0)
	}
	heap.Push(&states[0], &startState)
	loop := 0
	bestScore := 0
	var bestMove [2000]uint8
	stepCount := 0
	var i int
	for !timeout {
		loop++
		starti := maxInt(0, 0)
		for i = starti; i < 2000; i++ {
			if len(states[i]) == 0 {
				continue
			}
			nowState := heap.Pop(&states[i]).(*Item)
			addCount := 0
			for d := 0; d < 4; d++ {
				var nextPosition = nowState.value.position
				nextPosition.y += dy[d]
				nextPosition.x += dx[d]
				if nowState.value.canMove(nextPosition) {
					newState := &Item{
						value: nowState.value,
					}
					newState.value.Move(nextPosition, d)
					heap.Push(&states[i+1], newState)
				}
				addCount++
			}
			if addCount == 0 && nowState.value.scoreReal > bestScore {
				bestScore = nowState.value.scoreReal
				bestMove = nowState.value.move
				stepCount = nowState.value.stepCount
			}
		}

	}
	for i := 0; i < 2000; i++ {
		if len(states[i]) == 0 {
			continue
		}
		var j int
		for len(states[i]) != 0 && j < 20 {
			state := heap.Pop(&states[i]).(*Item)
			if state.value.scoreReal >= bestScore {
				bestScore = state.value.scoreReal
				bestMove = state.value.move
				stepCount = state.value.stepCount
			}
			j++
		}
	}

	var s string
	for i := 0; i < stepCount; i++ {
		fmt.Fprintln(os.Stderr, "-----BEGIN-----")
		s += fmt.Sprint(direction[bestMove[i]])
		fmt.Fprintln(os.Stderr, (s))
		fmt.Fprintln(os.Stderr, "-----END-----")
		fmt.Fprintln(os.Stderr, "index = 0: score = 0")
		fmt.Print(direction[bestMove[i]])
	}
	fmt.Println("")
	//	log.Printf("loop=%d score=%d\n", loop, bestScore)
}

// An Item is something we manage in a priority queue.
type Item struct {
	value State // The value of the item; arbitrary.
	index int   // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].value.score > pq[j].value.score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value State, priority int) {
	item.value = value
	heap.Fix(pq, item.index)
}
