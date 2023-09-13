package api

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
	"math/rand"
	"sync"
	"syscall"
	"testing"
	"time"
)

type myEpoll struct {
	fd int
}

func newMyEpoll() (*myEpoll, error) {
	epollFd, err := unix.EpollCreate(1024)
	if err != nil {
		return nil, err
	}
	return &myEpoll{
		fd: epollFd,
	}, nil
}

func (epoll *myEpoll) add(eventFd int) error {
	ee := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLHUP, // 监听可读事件
		Fd:     int32(eventFd),
	}
	return unix.EpollCtl(epoll.fd, syscall.EPOLL_CTL_ADD, eventFd, ee)
}

func (epoll *myEpoll) remove(eventFd int) error {
	return unix.EpollCtl(epoll.fd, syscall.EPOLL_CTL_DEL, eventFd, nil)
}

func (epoll *myEpoll) wait(waitEvents []unix.EpollEvent) (int, error) {
	return unix.EpollWait(epoll.fd, waitEvents, 1000)
}

func TestEpoll(t *testing.T) {
	fdSize := 15000
	epoller, err := newMyEpoll()
	assert.NoError(t, err)

	fds := make([]int, fdSize)
	for i := 0; i < fdSize; i++ {
		// 事件类型 fd，就是专门用于事件通知的文件描述符
		// 可以用于进程间的通信、用户态和内核态的通信
		// eventfd 内保存的是一个计数器:
		// 当计数值不为 0 就是有可读事件的发生
		// write 会让 eventfd 的计数器递增
		// read 会让 eventfd 的计数器清零
		// linux 下 /proc 进程都有进程 id 号的目录，下面可以看到使用资源信息
		// fd 就是其中一个
		// /proc/<pid>/fd --> anon_inode:[eventfd]
		// eventfd 的使用 api:
		// 1. read; 2. write; 3. watch; 4: close
		fd, err := unix.Eventfd(0, unix.EFD_NONBLOCK) // 非阻塞类型
		assert.NoError(t, err)
		fds[i] = fd
		err = epoller.add(fd)
		assert.NoError(t, err)
	}

	writeTimes := 50
	wg := sync.WaitGroup{}
	wg.Add(writeTimes)
	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// 写入必须是 8 字节大小
		var val uint64 = 1
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, val)
		for i := 0; i < writeTimes; i++ {
			time.Sleep(200 * time.Millisecond)
			_fd := fds[r.Intn(fdSize)]
			n, err := unix.Write(_fd, bs)
			assert.NoError(t, err)
			t.Logf("g1 write, event fd = %d, res = %d", _fd, n)
		}
	}()

	go func() {
		for {
			waitEvents := make([]unix.EpollEvent, fdSize)
			n, err := epoller.wait(waitEvents)
			assert.NoError(t, err)
			t.Logf("g2 wait, num = %d, total fd num = %d, epoll fd = %d",
				n, fdSize, epoller.fd)
			for j := 0; j < n; j++ {
				// 读取也必须是 8 字节大小
				res := make([]byte, 8)
				rn, err := unix.Read(int(waitEvents[j].Fd), res)
				assert.NoError(t, err)
				t.Logf("g2 read, event fd = %d, res = %d, %+v",
					waitEvents[j].Fd, rn, res)
				wg.Done()
			}
		}
	}()
	wg.Wait()
	for i := 0; i < fdSize; i++ {
		unix.CloseOnExec(fds[i])
	}
}
