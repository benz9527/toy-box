package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// 测试文件生成参考 dd if=/dev/zero of=test bs=1M count=1000

func TestCopyNFilesUnderDir(t *testing.T) {
	type args struct {
		outFilename string
		inDir       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				outFilename: "/tmp/splice.txt",
				inDir:       "/root/splice",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyNFilesUnderDir(tt.args.outFilename, tt.args.inDir); (err != nil) != tt.wantErr {
				t.Errorf("CopyNFilesUnderDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyNFilesUnderDirBySplice(t *testing.T) {
	type args struct {
		outFilename string
		inDir       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				outFilename: "/tmp/splice-2.txt",
				inDir:       "/root/splice",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyNFilesUnderDirBySplice(tt.args.outFilename, tt.args.inDir); (err != nil) != tt.wantErr {
				t.Errorf("CopyNFilesUnderDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// go test -run none -bench BenchmarkCopyNFilesUnderDir -benchtime=1000x -benchmem --parallel 2
// 在benchmark测试中,通常会多次执行相同的操作来获取平均时间。
// 下面的操作会导致生成的文件体积是正常的两倍。

func BenchmarkCopyNFilesUnderDir(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		err := CopyNFilesUnderDir(filepath.Join(os.TempDir(), fmt.Sprintf("splice%d.txt", i)), filepath.Join(os.TempDir(), "splice"))
		b.StopTimer()
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}

// io.Copy() 的性能是上面 naive approach 的两倍

func BenchmarkCopyNFilesUnderDirBySplice(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		err := CopyNFilesUnderDirBySplice(filepath.Join(os.TempDir(), fmt.Sprintf("splice-2-%d.txt", i)), filepath.Join(os.TempDir(), "splice"))
		b.StopTimer()
		assert.NoError(b, err)
	}
	b.ReportAllocs()
}
