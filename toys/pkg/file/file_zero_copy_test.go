package file_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/benz9527/toy-box/toys/pkg/file"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestZeroCopyFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zero Copy File Suite",
		types.SuiteConfig{
			LabelFilter:     "ZeroCopy",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

// 测试文件生成参考 dd if=/dev/zero of=test bs=1M count=1000

type copyNFilesUnderDirTestCase struct {
	args struct {
		outFilename string
		inDir       string
	}
	wantErr bool
}

var _ = DescribeTable("Copy N files under dir",
	Label("ZeroCopy"),
	func(tc *copyNFilesUnderDirTestCase) {
		err := os.MkdirAll(tc.args.inDir, os.ModePerm)
		assert.NoError(GinkgoT(), err)
		dir, _ := filepath.Split(tc.args.outFilename)
		err = os.MkdirAll(dir, os.ModePerm)
		assert.NoError(GinkgoT(), err)

		DeferCleanup(func() {
			assert.NoError(GinkgoT(), os.RemoveAll(tc.args.inDir))
			assert.NoError(GinkgoT(), os.RemoveAll(dir))
		})

		err = CopyNFilesUnderDir(tc.args.outFilename, tc.args.inDir)
		if tc.wantErr {
			assert.Error(GinkgoT(), err)
		} else {
			assert.NoError(GinkgoT(), err)
		}
	},
	Entry("1", &copyNFilesUnderDirTestCase{
		args: struct {
			outFilename string
			inDir       string
		}{
			outFilename: "/tmp/n-files-merge/merged.txt",
			inDir:       "/tmp/n-files",
		},
		wantErr: false,
	}),
)

func TestZeroCopyFileBySplice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zero Copy File By Splice Suite",
		types.SuiteConfig{
			LabelFilter:     "ZeroCopySplice",
			ParallelTotal:   1,
			ParallelProcess: 1,
			GracePeriod:     5 * time.Second,
		},
		types.ReporterConfig{
			Verbose: true,
		},
	)
}

type copyNFilesUnderDirBySpliceTestCase struct {
	copyNFilesUnderDirTestCase
}

var _ = DescribeTable("Copy N files under dir by splice",
	Label("ZeroCopySplice"),
	func(tc *copyNFilesUnderDirBySpliceTestCase) {
		err := os.MkdirAll(tc.args.inDir, os.ModePerm)
		assert.NoError(GinkgoT(), err)
		dir, _ := filepath.Split(tc.args.outFilename)
		err = os.MkdirAll(dir, os.ModePerm)
		assert.NoError(GinkgoT(), err)

		DeferCleanup(func() {
			assert.NoError(GinkgoT(), os.RemoveAll(tc.args.inDir))
			assert.NoError(GinkgoT(), os.RemoveAll(dir))
		})

		err = CopyNFilesUnderDirBySplice(tc.args.outFilename, tc.args.inDir)
		if tc.wantErr {
			assert.Error(GinkgoT(), err)
		} else {
			assert.NoError(GinkgoT(), err)
		}
	},
	Entry("1", &copyNFilesUnderDirBySpliceTestCase{
		copyNFilesUnderDirTestCase{
			args: struct {
				outFilename string
				inDir       string
			}{
				outFilename: "/tmp/n-files-splice-merge/splice-merge.txt",
				inDir:       "/tmp/n-files-splice",
			},
			wantErr: false,
		},
	}),
)

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
