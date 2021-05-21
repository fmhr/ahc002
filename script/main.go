package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {
	log.SetFlags(log.Lshortfile)
	P, _ = os.Getwd()
	//testRun()
	parallelRun()
}

var P string = ""

func testn(n int) {
	sumScore := 0
	for i := 0; i < n; i++ {
		fmt.Print("case=", i)
		score, loop := run(i)
		fmt.Printf(" score=%d loop=%d\n", score, loop)
		sumScore += score
	}
	fmt.Println("ALL SCORE = ", sumScore)
}

func testRun() {
	score, n := run(0)
	log.Printf("score=%d loop=%d\n", score, n)
	// vscore := vis(inputpaths[50], out)
}

func run(seed int) (int, int) {
	exe := P + "/solver"
	inFile := fmt.Sprintf("%s/tools/in/%s.txt", P, fmt.Sprintf("%04d", seed))
	outFile := fmt.Sprintf("%s/out/%s.out", P, fmt.Sprintf("%04d", seed))
	cmdStr := exe + " < " + inFile + " > " + outFile
	cmds := []string{"sh", "-c", cmdStr}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Println(cmds)
		log.Fatal(err)
	}
	cmd.Wait()
	score := parseScore(stderr.String())
	if score == 0 {
		log.Println(stderr.String())
	}
	loop := parseLoop(stderr.String())
	return score, loop
}

type Date struct {
	seed  int
	score int
	loop  int
}

func parallelRun() {
	CORE := 4
	maxSeed := 10
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE-1)
	datas := make([]Date, 0)
	for seed := 0; seed < maxSeed; seed++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			var d Date
			d.score, d.loop = run(seed)
			d.seed = seed
			mu.Lock()
			datas = append(datas, d)
			// fmt.Print(".")
			log.Printf("seed=%d score=%d loop=%d\n", d.seed, d.score, d.loop)
			mu.Unlock()
			wg.Done()
			<-sem
		}(seed)
	}
}

func parseScore(s string) int {
	ms := `score=([0-9]+)`
	re := regexp.MustCompile(ms)
	ma := re.FindString(s)
	score, err := strconv.Atoi(strings.Replace(ma, "score=", "", -1))
	if err != nil {
		log.Println(score)
	}
	return score
}

func parseLoop(s string) int {
	ms := `loop=([0-9]+)`
	re := regexp.MustCompile(ms)
	ma := re.FindString(s)
	n, err := strconv.Atoi(strings.Replace(ma, "loop=", "", -1))
	if err != nil {
		log.Println(n)
	}
	return n
}

func vis(input string, output string) (score int) {
	vispath := P + "/tools/target/release/vis"
	cmdStr := vispath + " " + input + " " + output
	cmds := []string{"sh", "-c", cmdStr}
	var out []byte
	var err error
	out, err = exec.Command(cmds[0], cmds[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}
	outs := strings.Split(string(out), "\n")
	score, err = strconv.Atoi(outs[0])
	if err != nil {
		panic(err)
	}
	return score
}
