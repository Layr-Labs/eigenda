//go:build icicle
package bench_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	iciclecore "github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ecntt"
	iciclebn254Msm "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/msm"
	iciclebn254Ntt "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicleruntime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"

	gnarkbn254fft "github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
)

// The benchmarks in this file are meant to test primitives in isolation: FFTFr, FFTG1, MSMG1.
// These should be compared to the gnark-crypto (CPU) implementations in benchmark_primitives_test.go
// TODO: The current implementations use async APIs but are written in a blocking sync way.
// To get optimal performance out of a GPU we would need to batch and pipeline multiple operations.

// deviceType should be one of "CUDA", "METAL", "CPU".
//
// CPU:
// Afaiu there is no point in using CPU device other than for testing the code witout a GPU.
// CPU icicle code will always be slower than gnark-crypto code running on CPU,
// since it requires some data conversions (e.g. field elements are stored in montgomery form in
// gnark-crypto, but not in icicle).
//
// METAL:
// Only works on macos, and requires github.com/ingonyama-zk/icicle/v3 v3.9.0.
// Install icicle dynamic libraries following https://dev.ingonyama.com/setup,
// and make them available using (/usr/local/icicle/lib is the recommended install location):
// export CGO_LDFLAGS="-L/usr/local/icicle/lib -lstdc++ -Wl,-rpath,/usr/local/icicle/lib"
//
// CUDA: TODO (not tested yet)
const deviceType = "METAL"

func BenchmarkIcicleFFTFr(b *testing.B) {
	icicleruntime.LoadBackendFromEnvOrDefault()
	device := icicleruntime.CreateDevice(deviceType, 0)

	for _, numFrsPowerOf2 := range []uint8{9, 14, 19, 22} {
		b.Run(fmt.Sprintf("2^%d_Points", numFrsPowerOf2), func(b *testing.B) {
			// We have to do this inside b.Run() to make sure all DeviceSlices are on the same device.
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			icicleruntime.SetDevice(&device)

			cfgBn254 := iciclebn254Ntt.GetDefaultNttConfig()
			cfgBn254.IsAsync = true
			streamBn254, _ := icicleruntime.CreateStream()
			cfgBn254.StreamHandle = streamBn254

			numScalars := 1 << numFrsPowerOf2
			scalarsBn254 := iciclebn254.GenerateScalars(numScalars)

			cfgInitDomainBls := iciclecore.GetDefaultNTTInitDomainConfig()
			rouMontBn254, _ := gnarkbn254fft.Generator(uint64(numScalars))
			rouBn254 := rouMontBn254.Bits()
			rouIcicleBn254 := iciclebn254.ScalarField{}
			limbsBn254 := iciclecore.ConvertUint64ArrToUint32Arr(rouBn254[:])
			rouIcicleBn254.FromLimbs(limbsBn254)
			iciclebn254Ntt.InitDomain(rouIcicleBn254, cfgInitDomainBls)

			var nttResultBn254 iciclecore.DeviceSlice
			_, e := nttResultBn254.MallocAsync(scalarsBn254.SizeOfElement(), numScalars, streamBn254)
			require.Equal(b, icicleruntime.Success, e, fmt.Sprint("Bn254 Malloc failed: ", e))

			for b.Loop() {
				err := iciclebn254Ntt.Ntt(scalarsBn254, iciclecore.KForward, &cfgBn254, nttResultBn254)
				require.Equal(b, icicleruntime.Success, err, fmt.Sprint("bn254 Ntt failed: ", err))
				nttResultBn254Host := make(iciclecore.HostSlice[iciclebn254.ScalarField], numScalars)
				nttResultBn254Host.CopyFromDeviceAsync(&nttResultBn254, streamBn254)
				icicleruntime.SynchronizeStream(streamBn254)
			}
			nttResultBn254.FreeAsync(streamBn254)
			icicleruntime.SynchronizeStream(streamBn254)
		})
	}
}

func BenchmarkIcicleMSMG1(b *testing.B) {
	icicleruntime.LoadBackendFromEnvOrDefault()
	device := icicleruntime.CreateDevice(deviceType, 0)

	for _, numG1PointsPowOf2 := range []uint8{12, 15, 19} {
		b.Run(fmt.Sprintf("2^%d_Points", numG1PointsPowOf2), func(b *testing.B) {
			// We have to do this inside b.Run() to make sure all DeviceSlices are on the same device.
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			icicleruntime.SetDevice(&device)

			cfgBn254 := iciclecore.GetDefaultMSMConfig()
			cfgBn254.IsAsync = true
			streamBn254, _ := icicleruntime.CreateStream()
			cfgBn254.StreamHandle = streamBn254

			msmResultBn254Host := make(iciclecore.HostSlice[iciclebn254.Projective], 1)
			var msmResultBn254 iciclecore.DeviceSlice
			_, e := msmResultBn254.MallocAsync(msmResultBn254Host.AsPointer().Size(), 1, streamBn254)
			require.Equal(b, icicleruntime.Success, e, fmt.Sprint("Bn254 Malloc failed: ", e))

			numG1Points := 1 << numG1PointsPowOf2
			scalarsBn254 := iciclebn254.GenerateScalars(numG1Points)
			pointsBn254 := iciclebn254.GenerateAffinePoints(numG1Points)

			for b.Loop() {
				err := iciclebn254Msm.Msm(scalarsBn254, pointsBn254, &cfgBn254, msmResultBn254)
				require.Equal(b, icicleruntime.Success, err, fmt.Sprint("bn254 Msm failed: ", err))
				msmResultBn254Host.CopyFromDeviceAsync(&msmResultBn254, streamBn254)
				icicleruntime.SynchronizeStream(streamBn254)
			}
			msmResultBn254.FreeAsync(streamBn254)
			icicleruntime.SynchronizeStream(streamBn254)
		})
	}
}

// ECNTT is not implemented on METAL. Only available on CUDA.
func BenchmarkIcicleFFTG1(b *testing.B) {
	icicleruntime.LoadBackendFromEnvOrDefault()
	device := icicleruntime.CreateDevice(deviceType, 0)

	for _, sizePowOf2 := range []uint8{13, 14} {
		b.Run(fmt.Sprintf("2^%d_Points", sizePowOf2), func(b *testing.B) {
			// We have to do this inside b.Run() to make sure all DeviceSlices are on the same device.
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			icicleruntime.SetDevice(&device)

			cfgBn254 := iciclebn254Ntt.GetDefaultNttConfig()
			cfgBn254.IsAsync = true
			streamBn254, _ := icicleruntime.CreateStream()
			cfgBn254.StreamHandle = streamBn254

			numG1Points := 1 << sizePowOf2
			pointsBn254 := iciclebn254.GenerateAffinePoints(numG1Points)

			var nttResultBn254 iciclecore.DeviceSlice
			_, e := nttResultBn254.MallocAsync(pointsBn254.SizeOfElement(), numG1Points, streamBn254)
			require.Equal(b, icicleruntime.Success, e, fmt.Sprint("Bn254 Malloc failed: ", e))

			for b.Loop() {
				err := ecntt.ECNtt(pointsBn254, iciclecore.KForward, &cfgBn254, nttResultBn254)
				require.Equal(b, icicleruntime.Success, err, fmt.Sprint("bn254 Ntt failed: ", err))
				nttResultBn254Host := make(iciclecore.HostSlice[iciclebn254.Affine], numG1Points)
				nttResultBn254Host.CopyFromDeviceAsync(&nttResultBn254, streamBn254)
				icicleruntime.SynchronizeStream(streamBn254)
			}
			nttResultBn254.FreeAsync(streamBn254)
			icicleruntime.SynchronizeStream(streamBn254)
		})
	}
}
