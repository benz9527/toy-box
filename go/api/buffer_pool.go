package api

import "sync"

var (
	bufferPool16k sync.Pool
	bufferPool4k  sync.Pool
	bufferPool2k  sync.Pool
	bufferPool1k  sync.Pool
	bufferPool    sync.Pool
)

const (
	_1k  = 1024
	_2k  = 2 * _1k
	_4k  = 2 * _2k
	_16k = 4 * _4k
)

func GetBuffer(size int) []byte {
	var _buf any
	switch {
	case size >= _16k:
		_buf = bufferPool16k.Get()
	case size >= _4k:
		_buf = bufferPool4k.Get()
	case size >= _2k:
		_buf = bufferPool2k.Get()
	case size >= _1k:
		_buf = bufferPool1k.Get()
	case size < _1k:
		_buf = bufferPool.Get()
	}
	if _buf == nil {
		return make([]byte, size)
	}
	buf := _buf.([]byte)
	if cap(buf) < size {
		return make([]byte, size)
	}
	return buf[:size]
}

func PutBuffer(buffer []byte) {
	size := cap(buffer)
	switch {
	case size >= _16k:
		bufferPool16k.Put(buffer)
	case size >= _4k:
		bufferPool4k.Put(buffer)
	case size >= _2k:
		bufferPool2k.Put(buffer)
	case size >= _1k:
		bufferPool1k.Put(buffer)
	case size < _1k:
		bufferPool.Put(buffer)
	}
}
