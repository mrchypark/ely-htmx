package httpencoding

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// EncoderSet reuses at most one encoder per codec for a bounded, serial batch.
// It is not safe for concurrent use. Options must be normalized; every returned
// representation owns its bytes and remains valid after later Encode or Close
// calls. Close releases the encoder workspace when the batch is complete.
type EncoderSet struct {
	options Options
	buffer  bytes.Buffer
	gzip    *gzip.Writer
	brotli  *brotli.Writer
	zstd    *zstd.Encoder
	closed  bool
}

// NewEncoderSet creates a lazy batch encoder. No compression workspaces are
// allocated for codecs the batch never uses.
func NewEncoderSet(options Options) *EncoderSet {
	return &EncoderSet{options: options}
}

func (s *EncoderSet) Encode(data []byte, coding string) ([]byte, error) {
	if s.closed {
		return nil, fmt.Errorf("compression encoder set is closed")
	}
	s.buffer.Reset()
	switch coding {
	case "identity":
		return bytes.Clone(data), nil
	case "gzip":
		if s.gzip == nil {
			w, err := gzip.NewWriterLevel(&s.buffer, s.options.GzipLevel)
			if err != nil {
				return nil, err
			}
			s.gzip = w
		} else {
			s.gzip.Reset(&s.buffer)
		}
		if _, err := s.gzip.Write(data); err != nil {
			return nil, err
		}
		if err := s.gzip.Close(); err != nil {
			return nil, err
		}
	case "br":
		if s.brotli == nil {
			s.brotli = brotli.NewWriterLevel(&s.buffer, s.options.BrotliLevel)
		} else {
			s.brotli.Reset(&s.buffer)
		}
		if _, err := s.brotli.Write(data); err != nil {
			return nil, err
		}
		if err := s.brotli.Close(); err != nil {
			return nil, err
		}
	case "zstd":
		if s.zstd == nil {
			w, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(s.options.ZstdLevel)))
			if err != nil {
				return nil, err
			}
			s.zstd = w
		}
		return s.zstd.EncodeAll(data, nil), nil
	default:
		return nil, fmt.Errorf("unsupported content coding %q", coding)
	}
	return bytes.Clone(s.buffer.Bytes()), nil
}

// Close is idempotent. Gzip and Brotli streams are finalized by each Encode;
// dropping their references releases their reusable compression workspaces.
func (s *EncoderSet) Close() {
	if s.zstd != nil {
		s.zstd.Close()
	}
	s.gzip = nil
	s.brotli = nil
	s.zstd = nil
	s.buffer = bytes.Buffer{}
	s.closed = true
}
