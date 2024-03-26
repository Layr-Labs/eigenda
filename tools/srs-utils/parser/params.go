package parser

type Params struct {
	NumPoint         uint64
	NumTotalG1Points uint64
	G1Size           uint64
	G2Size           uint64

	G1StartByte uint64
	G2StartByte uint64
}

func (p *Params) SetG1StartBytePos(startPoint uint64) {
	p.G1StartByte = startPoint*p.G1Size + OffsetToG1
}

func (p *Params) SetG2StartBytePos(startPoint uint64) {
	p.G2StartByte = startPoint*p.G2Size + OffsetToG1 + p.NumTotalG1Points*p.G1Size
}

func (p *Params) GetG1EndBytePos() uint64 {
	return p.G1StartByte + uint64(p.NumPoint*p.G1Size)
}

func (p *Params) GetG2EndBytePos() uint64 {
	return p.G2StartByte + uint64(p.NumPoint*p.G2Size)
}
