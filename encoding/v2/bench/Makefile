bench_all: bench_primitives bench_eigenda

# This downloads the SRS tables used by the prover benchmarks in bench_eigenda.
# This will download SRS tables for all cosets up to 512. Its ~500MB in size.
# Our benchmarks currently only use up to coset32, however, so if download is slow feel free to ctrl-c early.
# TODO(samlaf): now that we have this we can increase encoded_blob sizes benchmarked against up to 128MiB (8*16MiB).
download_srs_tables:
	cd ../../../tools/srs-utils && go run ./cmd/main.go download-tables --dimension dimE1024 --output-dir ../../resources/srs/SRSTables

# If running this benchmark for the first time, run `make download_srs_tables` first to get the SRS tables.
# This will greatly speed up the first run of the benchmark (which otherwise will generate and write the SRS tables to disk itself).
# Benchmark on ec2 g6.xlarge went from ~400s to ~120s downloading the srs tables first.
bench_eigenda:
	$(eval GOOS=$(shell go env GOOS))
	$(eval GOARCH=$(shell go env GOARCH))
	go test -benchmem -bench=. -run=^$$ benchmark_eigenda_test.go | tee results/golang_bench_eigenda_$(GOOS)_$(GOARCH).txt

bench_primitives:
	$(eval GOOS=$(shell go env GOOS))
	$(eval GOARCH=$(shell go env GOARCH))
# We could add the cpu info as well... but not sure how to get it reliably across platforms...
# `sysctl -n machdep.cpu.brand_string` will output "Apple M4 Pro" on macbook, but clearly this doesn't work on linux.
	go test -benchmem -bench=. -run=^$$ benchmark_primitives_test.go | tee results/golang_bench_primitives_$(GOOS)_$(GOARCH).txt

# benchstat pretty prints the results from the golang benchmarks
benchstat_results:
	benchstat results/*