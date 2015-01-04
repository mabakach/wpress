/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Yani Iliev <yani@iliev.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package wpress

import (
	"bytes"
	"errors"
	"os"
)

type Reader struct {
	Filename      string
	File          *os.File
	NumberOfFiles int
}

func NewReader(filename string) (*Reader, error) {
	// create a new instance of Reader
	r := &Reader{filename, nil, 0}

	// call the constructor
	err := r.Init()
	if err != nil {
		return nil, err
	}

	// return Reader instance
	return r, nil
}

func (r *Reader) Init() error {
	// try to open the file
	file, err := os.Open(r.Filename)
	if err != nil {
		return err
	}

	// file was openned, assign the handle to the holding variable
	r.File = file

	return nil
}

func (r Reader) ExtractFile(filename string, path string) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

func (r Reader) Extract(s string) error {
	// TODO: implement
	return nil
}

func (r Reader) GetHeaderBlock() ([]byte, error) {
	// create buffer to keep the header block
	block := make([]byte, headerSize)

	// read the header block
	bytesRead, err := r.File.Read(block)
	if err != nil {
		return nil, err
	}

	if bytesRead != headerSize {
		return nil, errors.New("Unable to read header block size")
	}

	return block, nil
}

func (r Reader) GetFilesCount() (int, error) {
	// test if we have enumerated the archive already
	if r.NumberOfFiles != 0 {
		return r.NumberOfFiles, nil
	}

	// put pointer at the beginning of the file
	r.File.Seek(0, 0)

	// loop until end of file was reached
	for {
		// read header block
		block, err := r.GetHeaderBlock()
		if err != nil {
			return 0, err
		}

		// initialize new header
		h := &Header{}

		// check if block equals EOF sequence
		if bytes.Compare(block, h.GetEofBlock()) == 0 {
			// EOF reached, stop the loop
			break
		}

		// populate header from our block bytes
		h.PopulateFromBytes(block)

		// set pointer after file content, to the next header block
		size, err := h.GetSize()
		if err != nil {
			return 0, err
		}
		r.File.Seek(int64(size), 1)

		// increment file counter
		r.NumberOfFiles++
	}

	return r.NumberOfFiles, nil
}