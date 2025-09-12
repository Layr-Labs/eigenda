package kzg

type KzgConfig struct {
	G1Path          string
	G2Path          string
	G2TrailingPath  string
	CacheDir        string
	NumWorker       uint64
	SRSOrder        uint64 // Order is the total size of SRS
	SRSNumberToLoad uint64 // Number of points to be loaded from the beginning
	Verbose         bool
	PreloadEncoder  bool
	LoadG2Points    bool
}
