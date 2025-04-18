package network_benchmark

type TestClient interface {
	GetData(size int64, seed int64) ([]byte, error)
}

type TestServer interface {
	SetRandomData(randomData *reusableRandomness)
}
