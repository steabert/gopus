package ogg

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	ogg_page_header_size = 27
	ogg_page_start_seq   = "OggS"
)

// ReadBOSPage reads the "Beginning Of Stream" page of an OGG file.
// Returns an error if the data does not start with a valid OGG page.
func ReadBOSPage(f io.ReadSeeker) ([]byte, error) {
	var err error

	buf := bufio.NewReader(f)

	ogg_page_header := make([]byte, ogg_page_header_size)
	_, err = io.ReadFull(buf, ogg_page_header)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(ogg_page_header[:len(ogg_page_start_seq)], []byte(ogg_page_start_seq)) {
		return nil, errors.New("expected an Ogg page")
	}

	ogg_number_page_segments := int(ogg_page_header[ogg_page_header_size-1])
	ogg_lacing_values := make([]byte, ogg_number_page_segments)
	_, err = io.ReadFull(buf, ogg_lacing_values)
	if err != nil {
		return nil, err
	}

	ogg_page_size := 0
	for _, lacing_value := range ogg_lacing_values {
		ogg_page_size += int(lacing_value)
	}

	ogg_bos_contents := make([]byte, ogg_page_size)
	_, err = io.ReadFull(buf, ogg_bos_contents)
	if err != nil {
		return nil, err
	}

	return ogg_bos_contents, nil
}
