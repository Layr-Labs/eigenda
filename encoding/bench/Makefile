build_cpu:
	go build -gcflags="all=-N -l" -ldflags="-s=false -w=false" -o bin/main main.go

build_icicle:
	go build -tags=icicle -gcflags="all=-N -l" -ldflags="-s=false -w=false" -o bin/main_icicle main.go

benchmark_default:
	go run main.go -cpuprofile cpu.prof -memprofile mem.prof

benchmark_icicle:
	go run -tags=icicle main.go -cpuprofile cpu.prof -memprofile mem.prof

cpu_profile:
	go tool pprof -http=:8080 cpu.prof

mem_profile:
	go tool pprof -http=:8080 mem.prof
