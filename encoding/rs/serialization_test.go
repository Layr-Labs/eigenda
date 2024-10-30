package rs

import (
	"bytes"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func TestFrameSerialization(t *testing.T) {
	// Helper to create frame with specific values
	newTestFrame := func(values []uint64) *Frame {
		frame := &Frame{
			Coeffs: make([]fr.Element, len(values)),
		}
		for i, v := range values {
			frame.Coeffs[i].SetUint64(v)
		}
		return frame
	}

	tests := []struct {
		name    string
		frame   *Frame
		wantErr bool
	}{
		{
			name:    "empty frame",
			frame:   newTestFrame([]uint64{}),
			wantErr: false,
		},
		{
			name:    "single coefficient",
			frame:   newTestFrame([]uint64{42}),
			wantErr: false,
		},
		{
			name:    "multiple coefficients",
			frame:   newTestFrame([]uint64{1, 2, 3, 4, 5}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Serialize
			serialized, err := Serialize(tt.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Frame.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify length
			expectedLen := len(tt.frame.Coeffs) * BYTES_PER_SYMBOL
			if len(serialized) != expectedLen {
				t.Errorf("Frame.Serialize() wrong length = %v, want %v", len(serialized), expectedLen)
			}

			// Test Deserialize
			deserialized, err := Deserialize(serialized)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare original and deserialized frames
			if len(deserialized.Coeffs) != len(tt.frame.Coeffs) {
				t.Errorf("Deserialize() wrong length = %v, want %v", len(deserialized.Coeffs), len(tt.frame.Coeffs))
			}

			for i := range tt.frame.Coeffs {
				if !bytes.Equal(tt.frame.Coeffs[i].Marshal(), deserialized.Coeffs[i].Marshal()) {
					t.Errorf("Coefficient mismatch at index %d", i)
				}
			}
		})
	}
}

func TestDeserializeErrors(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "nil data",
			data:    nil,
			wantErr: false, // Empty frame is valid
		},
		{
			name:    "invalid length",
			data:    make([]byte, BYTES_PER_SYMBOL+1), // Not multiple of BYTES_PER_SYMBOL
			wantErr: true,
		},
		{
			name:    "truncated data",
			data:    make([]byte, BYTES_PER_SYMBOL-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Deserialize(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
