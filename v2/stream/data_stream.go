package stream

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	bufferSize   = 64
	csvSeparator = ","
)

var (
	ErrCannotUnread = fmt.Errorf("it's not possible to unread something not readed")
	ErrEOF          = fmt.Errorf("reached end of file")
)

type DataStream interface {
	Read() ([]string, error)
	Unread() error
	Close() error
}

func NewCSVDataStream(file *os.File) (DataStream, error) {
	finfo, err := file.Stat()
	if err != nil {
		return nil, ErrEOF
	}

	return &csvDataStream{
		file:     file,
		previous: 0,
		readed:   0,
		next:     0,
		size:     finfo.Size(),
		ended:    false,
	}, nil
}

type csvDataStream struct {
	file        *os.File
	previous    int64
	readed      int64
	next        int64
	size        int64
	buffer      []byte
	bufferIndex int
	bufferSize  int
	ended       bool
}

func (cds *csvDataStream) Read() ([]string, error) {
	line, err := cds.readLine()
	if err != nil {
		return nil, err
	}

	rawColumns := strings.Split(line, csvSeparator)
	columns := make([]string, len(rawColumns))
	for i, c := range rawColumns {
		columns[i] = strings.TrimSpace(c)
	}

	return columns, nil
}

func (cds *csvDataStream) readLine() (string, error) {
	if cds.ended {
		return "", ErrEOF
	}

	currentRow := ""
	for {
		cds.readed = 0
		for index := cds.bufferIndex; index < cds.bufferSize && index < len(cds.buffer); index++ {
			cds.bufferIndex++
			b := cds.buffer[index]
			isBreakLine := b == '\n'
			if isBreakLine && currentRow == "" {
				cds.readed++
				continue
			}

			if isBreakLine {
				cds.previous = cds.next
				cds.next += int64(index)
				fmt.Println("row = '", currentRow, "', previous=", cds.previous, ", next=", cds.next, ", index=", index)
				cds.bufferSize -= index
				cds.bufferIndex = 0
				cds.buffer = cds.buffer[index:]
				if cds.previous == cds.next {
					cds.ended = true
				}
				return currentRow, nil
			}

			currentRow += string(b)
			cds.readed++
		}

		buffer := make([]byte, bufferSize)
		readed, err := cds.file.ReadAt(buffer, cds.next+cds.readed)
		if errors.Is(io.EOF, err) && readed == 0 {
			cds.ended = true
			return currentRow, nil
		}

		if err != nil && !errors.Is(io.EOF, err) {
			return "", err
		}

		cds.bufferIndex = 0
		cds.bufferSize = readed
		cds.buffer = buffer
	}
}

func (cds *csvDataStream) Unread() error {
	if cds.previous >= cds.next {
		return ErrCannotUnread
	}

	cds.next = cds.previous
	cds.buffer = make([]byte, 0)
	cds.ended = false
	return nil
}

func (cds *csvDataStream) Close() error {
	cds.buffer = make([]byte, 0)
	return cds.file.Close()
}
