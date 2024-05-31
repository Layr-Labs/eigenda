package codecs

type NoIFFTCodec struct {
	writeCodec BlobCodec
}

var _ BlobCodec = NoIFFTCodec{}

func NewNoIFFTCodec(writeCodec BlobCodec) NoIFFTCodec {
	return NoIFFTCodec{
		writeCodec: writeCodec,
	}
}

func (v NoIFFTCodec) EncodeBlob(data []byte) ([]byte, error) {
	return v.writeCodec.EncodeBlob(data)
}

func (v NoIFFTCodec) DecodeBlob(data []byte) ([]byte, error) {
	return GenericDecodeBlob(data)
}
