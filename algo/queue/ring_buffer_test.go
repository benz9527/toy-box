package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXRingBuffer_uint64(t *testing.T) {
	rb := NewXRingBuffer[uint64](1024)
	e := rb.LoadEntryByCursor(0)
	e.Store(0, 1)
	assert.Equal(t, uint64(1), e.GetValue())

	e = rb.LoadEntryByCursor(1023)
	e.Store(1023, 100)
	assert.Equal(t, uint64(100), e.GetValue())

	e = rb.LoadEntryByCursor(1024)
	e.Store(1024, 1000)
	assert.Equal(t, uint64(1000), e.GetValue())
}
