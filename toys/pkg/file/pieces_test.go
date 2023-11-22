package file_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/benz9527/toy-box/toys/pkg/file"
)

var _ = Describe("A big size zero file creation test", func() {
	type zeroFileTestCase struct {
		pathToFile string
		size       PieceSize
	}

	It("should create a zero file with size 150MB by ZeroFile", func(ctx SpecContext) {
		By("Initialize the test case 1")
		tc1 := &zeroFileTestCase{
			pathToFile: filepath.Join(os.TempDir(), "zero", "tc1", "1.zero"),
			size:       150 * MB,
		}
		tc1Info, err := os.Stat(tc1.pathToFile)
		Expect(err).To(MatchError(os.ErrNotExist))
		Expect(tc1Info).To(BeNil())

		err = os.MkdirAll(filepath.Dir(tc1.pathToFile), fs.ModePerm)
		GinkgoWriter.Printf("tc1 err: %+v\n", err)
		Expect(err).To(BeNil())

		f, err := os.OpenFile(tc1.pathToFile, os.O_CREATE|os.O_RDWR, fs.ModePerm)
		Expect(err).To(BeNil())
		Expect(f).NotTo(BeNil())
		DeferCleanup(func() {
			Expect(f).ToNot(BeNil())
			Expect(f.Close()).To(BeNil())
			Expect(os.RemoveAll(filepath.Dir(tc1.pathToFile))).To(BeNil())
		})

		err = FastZeroFile(f, tc1.size)
		Expect(err).To(BeNil())

		info, err := os.Stat(tc1.pathToFile)
		Expect(err).To(BeNil())
		Expect(info.Size()).To(Equal(int64(tc1.size)))
	}, SpecTimeout(3*time.Second))
	It("should create a zero file with size 100MB by ZeroFileByTruncate", func(ctx SpecContext) {
		By("Initialize the test case 2")
		tc2 := &zeroFileTestCase{
			pathToFile: filepath.Join(os.TempDir(), "zero", "tc2", "2.zero"),
			size:       100 * MB,
		}
		tc2Info, err := os.Stat(tc2.pathToFile)
		Expect(err).To(MatchError(os.ErrNotExist))
		Expect(tc2Info).To(BeNil())

		err = os.MkdirAll(filepath.Dir(tc2.pathToFile), fs.ModePerm)
		GinkgoWriter.Printf("tc2 err: %+v\n", err)

		f, err := os.OpenFile(tc2.pathToFile, os.O_CREATE|os.O_RDWR, fs.ModePerm)
		Expect(err).To(BeNil())
		Expect(f).NotTo(BeNil())

		DeferCleanup(func() {
			Expect(f).ToNot(BeNil())
			Expect(f.Close()).To(BeNil())
			Expect(os.RemoveAll(filepath.Dir(tc2.pathToFile))).To(BeNil())
		})

		err = FastZeroFileByTruncate(f, tc2.size)
		Expect(err).To(BeNil())

		info, err := os.Stat(tc2.pathToFile)
		Expect(err).To(BeNil())
		Expect(info.Size()).To(Equal(int64(tc2.size)))
	}, SpecTimeout(3*time.Second))
})

var _ = Describe("A big file cut into pieces", func() {
	type testcase struct {
		spec struct {
			pathToFile   string
			size         PieceSize
			shardingSize PieceSize
		}
		want struct {
			piecesCount int
		}
	}
	It("A 150MB file will be cut into 4 pieces, each pieces is le to 40MB", func(ctx SpecContext) {
		By("Initialize the test case")
		tc := &testcase{
			spec: struct {
				pathToFile   string
				size         PieceSize
				shardingSize PieceSize
			}{
				pathToFile:   filepath.Join(os.TempDir(), "zero", "tc", "1.zero"),
				size:         150 * MB,
				shardingSize: 40 * MB,
			},
			want: struct {
				piecesCount int
			}{
				piecesCount: 4,
			},
		}
		tc1Info, err := os.Stat(tc.spec.pathToFile)
		Expect(err).To(MatchError(os.ErrNotExist))
		Expect(tc1Info).To(BeNil())

		By("Mkdir all")
		err = os.MkdirAll(filepath.Dir(tc.spec.pathToFile), fs.ModePerm)
		GinkgoWriter.Printf("tc err: %+v\n", err)
		Expect(err).To(BeNil())

		By("Create big zero file")
		f, err := os.OpenFile(tc.spec.pathToFile, os.O_CREATE|os.O_RDWR, fs.ModePerm)
		Expect(err).To(BeNil())
		Expect(f).NotTo(BeNil())
		DeferCleanup(func() {
			By("Defer cleanup big zero file")
			Expect(f).ToNot(BeNil())
			Expect(f.Close()).To(BeNil())
			Expect(os.RemoveAll(filepath.Dir(tc.spec.pathToFile))).To(BeNil())
		})
		err = FastZeroFile(f, tc.spec.size)
		Expect(err).To(BeNil())
		info, err := os.Stat(tc.spec.pathToFile)
		Expect(err).To(BeNil())
		Expect(info.Size()).To(Equal(int64(tc.spec.size)))

		By("Cut big zero file into pieces")
		pieces, err := Cutting(tc.spec.pathToFile, tc.spec.shardingSize)
		Expect(err).To(BeNil())
		Expect(pieces).ToNot(BeNil())
		Expect(len(pieces)).To(Equal(tc.want.piecesCount))
		GinkgoWriter.Printf("pieces: %s\n", spew.Sdump(pieces))
	})
})
