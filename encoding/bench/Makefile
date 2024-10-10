build_cpu:
	go build -gcflags="all=-N -l" -ldflags="-s=false -w=false" -o bin/main_cpu main.go

build_gpu:
	go build -tags=gpu -gcflags="all=-N -l" -ldflags="-s=false -w=false" -o bin/main_gpu main.go

benchmark_cpu:
	go run main.go -cpuprofile cpu.prof -memprofile mem.prof

benchmark_gpu:
	go run -tags=gpu main.go -cpuprofile cpu.prof -memprofile mem.prof

cpu_profile:
	go tool pprof -http=:8080 cpu.prof

mem_profile:
	go tool pprof -http=:8080 mem.prof
