package debug

import (
	"bytes"
	"io"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
)

// Config is the config struct for the debug reader.
type Config struct {
	BufferSize   int `config:"buffer_size"`
	MaxExecution int `config:"max_execution"`
}

// DefaultConfig is the default configuration for the debug reader.
var DefaultConfig = Config{
	BufferSize:   16 * 1024,
	MaxExecution: 100,
}

const selector = "debug_reader"

// MatchByteFunc receive a byte and returns true if the byte match the predicate.
type MatchByteFunc func(b byte) bool

// ReporterFunc call called when a byte is match against the predicate.
type ReporterFunc func(pos int, raw []byte)

// Reader is a debug reader used to check the values of specific bytes from an io.Reader,
// Is is useful is you want to detect if you have received garbage from a network volume.
type Reader struct {
	log                 *logp.Logger
	reader              io.Reader
	buffer              bytes.Buffer
	bufferSize          int
	maxMatchedExecution int
	predicate           MatchByteFunc
	executionCount      int
	startDebugReader    sync.Once
	stopDebugReader     sync.Once
	reporter            ReporterFunc
}

// NewReader returns a debug reader.
func NewReader(
	reader io.Reader,
	bufferSize int,
	maxMatchedExecution int,
	predicate MatchByteFunc,
	reporter ReporterFunc,
) (*Reader, error) {
	r := &Reader{
		log:                 logp.NewLogger(selector),
		bufferSize:          bufferSize,
		reader:              reader,
		maxMatchedExecution: maxMatchedExecution,
		predicate:           predicate,
	}

	if reporter == nil {
		r.reporter = r.logReporter
	} else {
		r.reporter = reporter
	}
	return r, nil
}

// Read will proxy the read call to the original reader and will periodically checks the values of
// bytes in the buffer.
func (r *Reader) Read(p []byte) (int, error) {
	if r.executionCount > r.maxMatchedExecution {
		r.stopDebugReader.Do(func() {
			// cleanup any remaining bytes in the buffer.
			r.checkPendingBytes()
			r.log.Info("Stopping debug reader, max execution reached.")
		})
		return r.reader.Read(p)
	}

	n, err := r.reader.Read(p)

	r.startDebugReader.Do(func() {
		r.log.Info("Starting debug reader")
	})

	// We need to consume all the bytes in chunk from the original reader.
	if n != 0 {
		byteToConsume := n
		var consumeAt, upTo, available int
		for byteToConsume > 0 {
			// optimize what we can consume per iteration.
			available = r.bufferSize - r.buffer.Len()

			upTo = consumeAt + available
			if upTo > len(p) {
				upTo = len(p)
			}

			r.buffer.Write(p[consumeAt:upTo])

			byteToConsume -= upTo - consumeAt + 1
			consumeAt = upTo

			if r.buffer.Len() == r.bufferSize {
				if r.executionCount < r.maxMatchedExecution && r.checkPredicate() {
					r.executionCount++
				}
				r.buffer.Reset()
			}
		}
	}
	return n, err
}

func (r *Reader) checkPendingBytes() {
	if r.buffer.Len() > 0 {
		r.checkPredicate()
		r.buffer.Reset()
	}
}

func (r *Reader) logReporter(pos int, raw []byte) {
	r.log.Info("Matching byte found, position %d raw: %+v", pos, raw)
}

func (r *Reader) checkPredicate() bool {
	for idx, b := range r.buffer.Bytes() {
		if found := r.predicate(b); found {
			r.reporter(idx, r.buffer.Bytes())
			return true
		}
	}
	return false
}

// IsNullByte returns true if we receive a 0 or null byte.
// NOTES: we do not take into consideration UTF-8 code points.
func IsNullByte(b byte) bool {
	return b == 0x00
}
