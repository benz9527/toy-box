package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// CopyBuffer
// A --> B --> C
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		// Write to copy method, avoids an allocation and a copy.
		return wt.WriteTo(dst)
	}
	if rf, ok := dst.(io.ReaderFrom); ok {
		// Read from copy method, avoids an allocation and a copy.
		return rf.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024 // 32k
		if lr, ok := src.(*io.LimitedReader); ok && int64(size) > lr.N {
			if lr.N < 1 {
				size = 1
			} else {
				size = int(lr.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nRead, readErr := src.Read(buf)
		if nRead > 0 {
			nWrite, writeErr := dst.Write(buf[:nRead])
			if nWrite > 0 {
				written += int64(nWrite)
			}
			if writeErr != nil {
				err = writeErr
				break
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				err = readErr
			}
			break
		}
	}
	return
}

func DownloadFile(url, fileName string) error {
	rsp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download file failed %w", err)
	}
	defer func() {
		if rsp != nil {
			_ = rsp.Body.Close()
		}
	}()
	// If we use the io.ReadAll(), the memory will be exhausted.
	downloadFile, err := os.Create(filepath.Join(os.TempDir(), fileName))
	defer func() {
		if downloadFile != nil {
			_ = downloadFile.Close()
		}
	}()
	wt := bufio.NewWriter(downloadFile)
	if _, err = io.Copy(wt, rsp.Body); err != nil {
		return fmt.Errorf("copy file failed %w", err)
	}
	return wt.Flush()
}

func Join(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (inCount, outCount int64, err error) {
	var wg = sync.WaitGroup{}
	var errors = make([]error, 2)
	pipe := func(to io.ReadWriteCloser, from io.ReadWriteCloser, count *int64, errIdx int) {
		defer to.Close()
		defer from.Close()
		defer wg.Done()

		buffer := GetBuffer(_16k)
		defer PutBuffer(buffer)
		*count, errors[errIdx] = io.CopyBuffer(to, from, buffer)
	}
	wg.Add(2)
	go pipe(c1, c2, &inCount, 0)
	go pipe(c2, c1, &outCount, 1)
	wg.Wait()
	if errors[0] != nil {
		err = fmt.Errorf("%w", errors[0])
	}
	if errors[1] != nil {
		if err != nil {
			err = fmt.Errorf("%w", err)
		}
		err = fmt.Errorf("%w", errors[1])
	}
	return
}
