build: src/main.go
	go build


example:
	./src/solver < tools/example.in > example.out

vis:
	cd tools &&\
	cargo run --release --bin vis example.in example.out


buildcmd:
	cd script &&\
		go build

cmdtest:
	make buildcmd
	./script/script


pprof:
	.solver -cpuprofile cpu.prof < tools/in/0000.txt
	pprof -http=localhost:8080 cpu.prof
