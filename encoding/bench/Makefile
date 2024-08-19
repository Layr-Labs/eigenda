build_cpu:
	go build -gcflags="all=-N -l" -ldflags="-s=false -w=false" -o bin/main_cpu main.go

build_gpu:
	CGO_LDFLAGS="-L/home/ec2-user/github/icicle/icicle/build/lib -g" go build -gcflags="all=-N -l" -ldflags="-s=false -w=false" -tags=gpu -o bin/main_gpu main.go

benchmark_cpu:
	go run main.go -cpuprofile cpu.prof -memprofile mem.prof

benchmark_gpu:
	go run -tags=gpu main.go -cpuprofile cpu_gpu.prof -memprofile mem_gpu.prof

nvidia_profile_gpu: build_gpu
	nsys profile --stats=true ./bin/main_gpu  -cpuprofile gpu_onhost_comp.prof -memprofile gpu_onhost_comp.prof
