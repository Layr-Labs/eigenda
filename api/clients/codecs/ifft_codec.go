package codecs

import "fmt"

type IFFTCodec struct {
	writeCodec BlobCodec
}

var _ BlobCodec = IFFTCodec{}

func NewIFFTCodec(writeCodec BlobCodec) IFFTCodec {
	return IFFTCodec{
		writeCodec: writeCodec,
	}
}

func (v IFFTCodec) EncodeBlob(data []byte) ([]byte, error) {
	var err error
	data, err = v.writeCodec.EncodeBlob(data)
	if err != nil {
		// this cannot happen, because EncodeBlob never returns an error
		return nil, fmt.Errorf("error encoding data: %w", err)
	}

	return IFFT(data)
}

func (v IFFTCodec) DecodeBlob(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("blob has length 0, meaning it is malformed")
	}
	var err error
	data, err = FFT(data)
	if err != nil {
		return nil, fmt.Errorf("error FFTing data: %w", err)
	}

	return GenericDecodeBlob(data)
}
