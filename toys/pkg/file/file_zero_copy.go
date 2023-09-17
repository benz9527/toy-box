//go:build linux || windows

package file

import (
	"bufio"
	"errors"
	"fmt"
	Ants "github.com/panjf2000/ants/v2"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type copyFileInfo struct {
	inFile *os.File
	inBuf  []byte
	isEOF  bool
}

func newCopyFileInfo(inFile *os.File, inBytes []byte, eof ...bool) *copyFileInfo {
	cfi := &copyFileInfo{
		inFile: inFile,
		inBuf:  inBytes,
	}
	if len(eof) > 0 {
		cfi.isEOF = eof[0]
	}
	return cfi
}

type copyChannel struct {
	copyCh   chan *copyFileInfo
	closeCh  chan struct{}
	isClosed atomic.Bool
}

func newCopyChannel(size int) *copyChannel {
	cc := &copyChannel{
		copyCh:   make(chan *copyFileInfo, size),
		isClosed: atomic.Bool{},
	}
	cc.isClosed.Store(false)
	return cc
}

func (cc *copyChannel) Close() {
	if !cc.isClosed.Load() {
		close(cc.copyCh)
		cc.isClosed.Store(true)
	}
}

func (cc *copyChannel) Get() <-chan *copyFileInfo {
	return cc.copyCh
}

func (cc *copyChannel) Add(infos ...*copyFileInfo) error {
	if cc.isClosed.Load() {
		return errors.New("copy channel has been closed")
	}
	for i := 0; i < len(infos); i++ {
		cc.copyCh <- infos[i]
	}
	return nil
}

func CopyNFilesUnderDir(outFilename, inDir string) error {
	var (
		outFile *os.File
		inFiles = make([]string, 0, 128)
	)
	if fi, err := os.Stat(outFilename); fi != nil && fi.IsDir() {
		return errors.New("outFile file is a dir")
	} else if os.IsNotExist(err) {
		f, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		outFile = f
	} else {
		f, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		_, _ = f.Seek(0, io.SeekEnd) // 追加模式
		outFile = f
	}
	defer func() {
		_ = outFile.Close()
	}()

	if fi, err := os.Stat(inDir); os.IsNotExist(err) || fi != nil && !fi.IsDir() {
		return errors.New("in dir is not a dir or not exist")
	}
	if err := filepath.Walk(inDir, func(path string, fi fs.FileInfo, err error) error {
		if fi != nil && !fi.IsDir() && fi.Mode()&fs.ModeTemporary == 0 {
			inFiles = append(inFiles, filepath.Join(inDir, fi.Name()))
		}
		return nil
	}); err != nil {
		return err
	}

	// TODO phase2 监控目录下是否有新文件创建

	cc := newCopyChannel(2000)
	p, err := Ants.NewPool(11,
		Ants.WithMaxBlockingTasks(100),
		Ants.WithPanicHandler(func(i any) {

		}))
	if err != nil {
		return err
	}
	n := len(inFiles)
	wg := sync.WaitGroup{}
	wg.Add(n + 1)
	_ = p.Submit(func() {
		count := 0
	ReadLoop:
		for {
			select {
			case info := <-cc.Get():
				if info.isEOF {
					count++
					if count == n {
						wg.Done()
						break ReadLoop
					}
					continue ReadLoop
				}

				_, err := outFile.Write(append(info.inBuf, '\n'))
				if err != nil {
					wg.Done()
					break ReadLoop
				}
			}
		}
	})
	for len(inFiles) > 0 {
		var files []string
		if len(inFiles) <= 10 && len(inFiles) > 0 {
			files = inFiles
			inFiles = []string{}
		} else {
			files = inFiles[:10]
			inFiles = inFiles[10:]
		}
		for i := 0; i < len(files); i++ {
			_f := files[i]
			_ = p.Submit(func() {
				f, err := os.OpenFile(_f, os.O_RDONLY, 0666)
				if err != nil {
					// syscall.EBUSY   正在忙碌的文件，不能打开
					// syscall.ENOENT  目录或者文件不存在，no such file or directory
					return
				}
				defer func() {
					_ = cc.Add(newCopyFileInfo(f, nil, true))
					_ = f.Close()
					wg.Done()
				}()
				// seek 操作文档的指针位置
				// pos, err := f.Seek(0, io.SeekCurrent) // 获取文档当前的位置
				// f.Seek(10, io.SeekStart)              // 跳过开头的 10 个字节
				// f.Seek(pos, io.SeekStart)             // 跳回到原来的位置
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					if err = cc.Add(newCopyFileInfo(f, scanner.Bytes())); err != nil {
						break
					}
				}
			})
		}
	}
	wg.Wait()
	cc.Close()
	return nil
}

func CopyNFilesUnderDirBySplice(outFilename, inDir string) error {
	var (
		outFile *os.File
		inFiles = make([]string, 0, 128)
	)
	if fi, err := os.Stat(outFilename); fi != nil && fi.IsDir() {
		return errors.New("outFile file is a dir")
	} else if os.IsNotExist(err) {
		f, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		outFile = f
	} else {
		f, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		_, _ = f.Seek(0, io.SeekEnd) // 追加模式
		outFile = f
	}
	defer func() {
		_ = outFile.Close()
	}()

	if fi, err := os.Stat(inDir); os.IsNotExist(err) || fi != nil && !fi.IsDir() {
		return errors.New("in dir is not a dir or not exist")
	}
	if err := filepath.Walk(inDir, func(path string, fi fs.FileInfo, err error) error {
		if fi != nil && !fi.IsDir() && fi.Mode()&fs.ModeTemporary == 0 {
			inFiles = append(inFiles, filepath.Join(inDir, fi.Name()))
		}
		return nil
	}); err != nil {
		return err
	}

	// TODO phase2 监控目录下是否有新文件创建

	cc := newCopyChannel(2000)
	p, err := Ants.NewPool(11,
		Ants.WithMaxBlockingTasks(100),
		Ants.WithPanicHandler(func(i any) {

		}))
	if err != nil {
		return err
	}
	n := len(inFiles)
	wg := sync.WaitGroup{}
	wg.Add(n + 1)
	_ = p.Submit(func() {
		count := 0
	ReadLoop:
		for {
			select {
			case info := <-cc.Get():
				if info.isEOF {
					count++
					_ = info.inFile.Close()
					if count == n {
						wg.Done()
						break ReadLoop
					}
					continue ReadLoop
				}

				// SPLICE_F_NONBLOCK 使用非阻塞的 IO 模型
				// SPLICE_F_MORE 本次 splice 不是最后一次操作，让 linux kernel 优化合并多个 splice
				// SPLICE_F_MOVE 直接移动数据
				// SPLICE_F_GIFT 把数据所有权移交给目标文件
				// splice 每次调用都会涉及到临时的 pipe 创建和销毁
				// syscall.EBADF bad file descriptor
				// syscall.EINVAL invalid argument
				//
				// [ user space ]
				//
				//              -----------------------------------------------------------------------------------------------
				//             |                                           application                                         |
				//              -----------------------------------------------------------------------------------------------
				// ········································|····················|·····················|··································
				//                                         | splice()           | pipe()              | splice()
				//[ kernel space ]                         |                    |                     |
				//              -----------------        ”copy“       -----------------------      ”copy“       -----------------
				//             |  socket buffer  |· · · · · · · · · >| pipe writefd & readfd |· · · · · · · · >|  socket buffer  |
				//              -----------------                     -----------------------                   -----------------
				//                     | copy                                                                         |
				// ····················|··············································································|
				//[ hardware sapce ]   |                                                                              |
				//              -----------------------------------------------------------------------------------------------
				//             |                                         network interface                                     |
				//              -----------------------------------------------------------------------------------------------
				// splice 需要保证缓存在 kernel space 中流动
				// _, err := syscall.Splice(int(info.inFile.Fd()), &info.inStartOffset, int(outFile.Fd()), &woffset, info.inLen, unix.SPLICE_F_NONBLOCK|unix.SPLICE_F_MORE)
				// io.Copy() 其中的 read 函数实现的 splice pool，在 interal/poll 包下面
				_, err := io.Copy(outFile, info.inFile)
				if err != nil {
					wg.Done()
					break ReadLoop
				}
			}
		}
	})
	for len(inFiles) > 0 {
		var files []string
		if len(inFiles) <= 10 && len(inFiles) > 0 {
			files = inFiles
			inFiles = []string{}
		} else {
			files = inFiles[:10]
			inFiles = inFiles[10:]
		}
		for i := 0; i < len(files); i++ {
			_f := files[i]
			_ = p.Submit(func() {
				f, err := os.OpenFile(_f, os.O_RDONLY, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer func() {
					_ = cc.Add(newCopyFileInfo(f, nil, true))
					wg.Done()
				}()
				if err = cc.Add(newCopyFileInfo(f, nil)); err != nil {
					fmt.Println(err)
				}
			})
		}
	}
	wg.Wait()
	cc.Close()
	return nil
}
