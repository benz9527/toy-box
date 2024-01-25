package api

import (
	"errors"
	"io"
	"math/rand"
	"testing"
	"time"
)

type trickle struct {
	counter int
}

func (t *trickle) Read(buf []byte) (n int, err error) {
	if t.counter > 25 {
		return 0, io.EOF
	}
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	buf[0] = byte('a' + t.counter%26)
	t.counter++
	return 1, nil
}

func TestIOPipeToPrintAlphabet(t *testing.T) {
	src := &trickle{}
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer func() {
			if pipeWriter != nil {
				_ = pipeWriter.Close() // send done signal
			}
		}()
		if _, err := io.Copy(pipeWriter, src); err != nil && !errors.Is(err, io.EOF) {
			t.Logf("pipe copy error %s", err)
		}
	}()
	read := func(reader io.Reader) {
		buf := make([]byte, 1024)
		position := 0
		lastTs := time.Now()
		for {
			n, err := reader.Read(buf[position:])
			if n > 0 {
				position += n
			}
			if position >= len(buf) || time.Since(lastTs) > 100*time.Millisecond && position > 0 {
				t.Logf("read: %s\n", buf[:position])
				lastTs = time.Now()
				position = 0
			}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					t.Log(err)
				}
				break
			}
		}
		if position > 0 {
			t.Logf("read: %s\n", buf[:position])
		}
	}
	read(pipeReader) // must declare reader to read from Pipe
}
