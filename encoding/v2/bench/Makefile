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

# This downloads the SRS tables used by the prover benchmarks in golang_bench_eigenda.
# This will download SRS tables for all cosets up to 512. Its ~500MB in size.
# Our benchmarks currently only use up to coset32, however, so if download is slow feel free to ctrl-c early.
# TODO(samlaf): now that we have this we can increase encoded_blob sizes benchmarked against up to 128MiB (8*16MiB).
download_srs_tables:
	cd ../../../tools/srs-utils && go run ./cmd/main.go download-tables --dimension dimE8192 --output-dir ../../resources/srs/SRSTables

# If running this benchmark for the first time, run `make download_srs_tables` first to get the SRS tables.
# This will greatly speed up the first run of the benchmark (which otherwise will generate and write the SRS tables to disk itself).
# Benchmark on ec2 g6.xlarge went from ~400s to ~120s downloading the srs tables first.
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