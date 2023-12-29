package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXRingBuffer_uint64(t *testing.T) {
	rb := NewXRingBuffer[uint64](1024)
	rb.StoreEntry(0, 1)
	e, ok := rb.LoadEntry(0)
	assert.True(t, ok)
	assert.Equal(t, uint64(1), e.GetValue())

	rb.StoreEntry(1023, 100)
	e, ok = rb.LoadEntry(1023)
	assert.True(t, ok)
	assert.Equal(t, uint64(100), e.GetValue())

	rb.StoreEntry(1024, 1000)
	e, ok = rb.LoadEntry(1024)
	assert.True(t, ok)
	assert.Equal(t, uint64(1000), e.GetValue())
}
