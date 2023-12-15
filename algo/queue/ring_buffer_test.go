package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXRingBuffer_uint64(t *testing.T) {
	rb := NewXRingBuffer[uint64](1024)
	rb.StoreElement(0, 1)
	e, ok := rb.LoadElement(0)
	assert.True(t, ok)
	assert.Equal(t, uint64(1), e.GetValue())

	rb.StoreElement(1023, 100)
	e, ok = rb.LoadElement(1023)
	assert.True(t, ok)
	assert.Equal(t, uint64(100), e.GetValue())

	rb.StoreElement(1024, 1000)
	e, ok = rb.LoadElement(1024)
	assert.True(t, ok)
	assert.Equal(t, uint64(1000), e.GetValue())
}
