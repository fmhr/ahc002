build: src/main.go
	cd src &&\
		go build


"example:
"	./src/solver < tools/example.in > example.out

vis:
	cd tools &&\
	cargo run --release --bin vis example.in example.out


buildcmd:
	cd cmd &&\
		go build

cmdtest:
	make buildcmd
	./cmd/script


pprof:
	./src/solver -cpuprofile cpu.prof < tools/in/0000.txt
	pprof -http=localhost:8080 cpu.prof
