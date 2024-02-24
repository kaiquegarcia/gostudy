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
	file     *os.File
	previous int64
	readed   int64
	next     int64
	size     int64
	buffer   []byte
	ended    bool
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

	if len(columns) == 1 && columns[0] == "" {
		return nil, ErrEOF
	}

	return columns, nil
}

func (cds *csvDataStream) readLine() (string, error) {
	if cds.ended {
		return "", ErrEOF
	}

	phrase := ""
	cds.readed = 0
	for {
		for index, b := range cds.buffer {
			if b == '\n' && phrase != "" {
				cds.previous = cds.next
				cds.next += cds.readed
				cds.buffer = cds.buffer[index:]
				if cds.previous == cds.next {
					cds.ended = true
				}

				return phrase, nil
			}

			phrase += string(b)
			cds.readed++
		}

		if cds.next+cds.readed >= cds.size {
			cds.ended = true
			return phrase, nil
		}

		buffer := make([]byte, bufferSize)
		readed, err := cds.file.ReadAt(buffer, cds.next+cds.readed)
		if errors.Is(io.EOF, err) && readed == 0 {
			cds.ended = true
			return phrase, nil
		}

		if err != nil && !errors.Is(io.EOF, err) {
			return "", err
		}

		cds.buffer = buffer[:readed]
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
