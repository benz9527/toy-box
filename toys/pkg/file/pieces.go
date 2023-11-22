package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PieceSize = int64

const (
	Byte PieceSize = 1
	KB   PieceSize = 1024 * Byte
	MB   PieceSize = 1024 * KB
)

type Piece struct {
	hasBeenCut bool
	shardingCount,
	pieceIndex int
	startSeek,
	offset,
	shardingSize,
	realSize int64
	originalFilename,
	originalPath string
}

func (p *Piece) HasBeenCut() bool {
	return p.hasBeenCut
}

func (p *Piece) GetFilename() string {
	if p.hasBeenCut {
		return fmt.Sprintf("%s_%d.piece", p.originalFilename, p.pieceIndex)
	}
	return p.originalFilename
}

func (p *Piece) GetPath() string {
	return filepath.Join(p.originalPath, p.originalFilename)
}

func (p *Piece) GetSeekInfo() (int64, int64, int64) {
	return p.startSeek, p.offset, p.shardingSize
}

func (p *Piece) GetShardingInfo() (int, int) {
	return p.pieceIndex, p.shardingCount
}

func Cutting(pathToFile string, maxSize PieceSize) ([]*Piece, error) {
	var (
		info os.FileInfo
		err  error
	)
	if info, err = os.Stat(pathToFile); os.IsPermission(err) {
		return nil, fmt.Errorf("permission denied, unable to access file <%s>", pathToFile)
	} else if os.IsNotExist(err) {
		return nil, fmt.Errorf("file <%s> not exist %w", pathToFile, err)
	} else if err != nil {
		return nil, fmt.Errorf("unknown error %w", err)
	} else if info != nil && info.IsDir() {
		return nil, fmt.Errorf("<%s> is a directory", pathToFile)
	}

	dir, _ := filepath.Split(pathToFile)
	if info.Size() <= maxSize {
		return []*Piece{
			{
				hasBeenCut:       false,
				shardingCount:    1,
				pieceIndex:       0,
				startSeek:        0,
				offset:           info.Size(),
				shardingSize:     maxSize,
				realSize:         info.Size(),
				originalFilename: info.Name(),
				originalPath:     dir,
			},
		}, nil
	}

	quotient := int(info.Size() / maxSize)
	remainder := info.Size() % maxSize
	shardingCount := quotient
	if remainder > 0 {
		shardingCount++
	}
	pieces := make([]*Piece, 0, shardingCount)
	for i := 0; i < quotient; i++ {
		pieces = append(pieces, &Piece{
			hasBeenCut:       true,
			shardingCount:    shardingCount,
			pieceIndex:       i,
			startSeek:        int64(i) * maxSize,
			offset:           maxSize,
			shardingSize:     maxSize,
			realSize:         info.Size(),
			originalFilename: info.Name(),
			originalPath:     dir,
		})
	}
	if remainder > 0 {
		seek := int64(shardingCount-1) * maxSize
		pieces = append(pieces, &Piece{
			hasBeenCut:       true,
			shardingCount:    shardingCount,
			pieceIndex:       shardingCount - 1,
			startSeek:        seek,
			offset:           info.Size() - seek,
			shardingSize:     maxSize,
			realSize:         info.Size(),
			originalFilename: info.Name(),
			originalPath:     dir,
		})
	}
	return pieces, nil
}

type TruncateWriter interface {
	io.WriteSeeker
	Truncate(size int64) error
}

// FastZeroFile Create an empty big file quickly.
func FastZeroFile(w TruncateWriter, fileSize int64) (err error) {
	if _, err = w.Seek(fileSize-1, io.SeekStart); err != nil {
		return err
	}
	if _, err = w.Write([]byte{0}); err != nil {
		return err
	}
	return nil
}

func FastZeroFileByTruncate(w TruncateWriter, fileSize int64) (err error) {
	if err = w.Truncate(fileSize); err != nil {
		return err
	}
	return nil
}
