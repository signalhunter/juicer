package main

import (
	"compress/gzip"
	"net/http"

	"github.com/flier/gohs/hyperscan"
)

var buffersize = 1024 * 1024

type state struct {
	buffer []byte
	iter   int
}

func Scan(url string, db hyperscan.StreamDatabase, ch chan string) error {
	callback := func(id uint, from, to uint64, flags uint, context interface{}) error {
		ctx := context.(*state)
		offset := uint64(ctx.iter * buffersize)
		ch <- string(ctx.buffer[from-offset : to-offset])
		return nil
	}

	st := state{
		buffer: make([]byte, buffersize),
		iter:   0,
	}

	scratch, err := hyperscan.NewScratch(db)
	if err != nil {
		return err
	}
	defer scratch.Free()

	stream, err := db.Open(0, scratch, callback, &st)
	if err != nil {
		return err
	}
	defer stream.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()

	for {
		_, err := gz.Read(st.buffer)
		stream.Scan(st.buffer)
		st.iter += 1

		if err != nil {
			break
		}
	}

	return nil
}
