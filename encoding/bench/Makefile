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

#### Golang Benchmarks

golang_benchmarks: golang_bench_primitives golang_bench_eigenda

golang_bench_eigenda:
	$(eval GOOS=$(shell go env GOOS))
	$(eval GOARCH=$(shell go env GOARCH))
	go test -benchmem -bench=. -run=^$$ benchmark_eigenda_test.go | tee results/golang_bench_eigenda_$(GOOS)_$(GOARCH).txt

golang_bench_primitives:
	$(eval GOOS=$(shell go env GOOS))
	$(eval GOARCH=$(shell go env GOARCH))
# We could add the cpu info as well... but not sure how to get it reliably across platforms...
# `sysctl -n machdep.cpu.brand_string` will output "Apple M4 Pro" on macbook, but clearly this doesn't work on linux.
	go test -benchmem -bench=. -run=^$$ benchmark_primitives_test.go | tee results/golang_bench_primitives_$(GOOS)_$(GOARCH).txt

# benchstat pretty prints the results from the golang benchmarks
benchstat_results:
	benchstat results/*