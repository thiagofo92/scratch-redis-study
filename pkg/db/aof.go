package db

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"thiagofo92/scratch-redis/pkg/adapter"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   *sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("error to create aof %w", err)
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(data adapter.RespDataOutPut) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(data.Marshal())
	if err != nil {
		return err
	}

	return nil
}
