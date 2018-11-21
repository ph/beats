package debug

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
)

func TestIsGarbage(t *testing.T) {
	t.SkipNow()
	t.Run("return true if null byte is received", func(t *testing.T) {
		assert.True(t, IsNullByte(0x00))
	})

	t.Run("return false on anything other bytes", func(t *testing.T) {
		assert.False(t, IsNullByte(0x21)) // 0x21 is !
	})
}

type reporter struct {
	Pos int
	Raw []byte
}

func (r *reporter) log(pos int, raw []byte) {
	r.Pos = pos
	r.Raw = raw
}

func TestReader(t *testing.T) {
	t.Run("garbage byte found", testGarbage)
	t.Run("no garbage byte found", testNoGarbage)
	t.Run("consume all bytes", testConsumeAll)
	t.Run("empty buffer", testEmptyBuffer)
}

func testGarbage(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteByte(0x00)
	b.WriteString(" world")

	dup := make([]byte, b.Len())
	copy(dup, b.Bytes())

	r := &reporter{}

	debug, _ := NewReader(&b, 8, 20, IsNullByte, r.log)
	data := make([]byte, 9)
	n, err := debug.Read(data)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 9, n)
	assert.Equal(t, 5, r.Pos)
	assert.Equal(t, dup[:8], r.Raw)
}

func testNoGarbage(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("hello world")

	r := &reporter{}

	debug, _ := NewReader(&b, 8, 20, IsNullByte, r.log)
	data := make([]byte, 9)
	n, err := debug.Read(data)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 9, n)
	assert.Equal(t, 0, r.Pos)
	assert.Equal(t, 0, len(r.Raw))
}

func testConsumeAll(t *testing.T) {
	c, _ := common.RandomBytes(2000)
	var buf bytes.Buffer
	consumed := 0

	reader := bytes.NewReader(c)
	debug, _ := NewReader(reader, 8, 20, IsNullByte, nil)

	for consumed < 2000 {
		data := make([]byte, 33)
		n, _ := debug.Read(data)

		buf.Write(data[:n])
		consumed += n
	}

	assert.Equal(t, len(c), consumed)
	assert.Equal(t, c, buf.Bytes())
}

func testEmptyBuffer(t *testing.T) {
	var buf bytes.Buffer

	debug, _ := NewReader(&buf, 8, 20, IsNullByte, nil)

	data := make([]byte, 33)
	n, err := debug.Read(data)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
}
