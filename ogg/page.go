package ogg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	ogg_page_header_size = 27
	ogg_capture_pattern  = "OggS"
)

type OggPage struct {
	Body            []byte
	GranulePosition int64
	SerialNumber    uint32
	SequenceNumber  uint32
	Checksum        uint32
	Continued       bool
	FirstPage       bool
	LastPage        bool
	Complete        bool
}

// ParsePage parses a single page of an OGG stream.
func ParsePage(r io.Reader) (OggPage, error) {
	var err error
	var page OggPage

	ogg_page_header := make([]byte, ogg_page_header_size)
	_, err = io.ReadFull(r, ogg_page_header)
	if err != nil {
		return page, err
	}

	//	0                   1                   2                   3
	//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1| Byte
	//
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// | capture_pattern: Magic number for page start "OggS"           | 0-3
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// | version       | header_type   | granule_position              | 4-7
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                                                               | 8-11
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                               | bitstream_serial_number       | 12-15
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                               | page_sequence_number          | 16-19
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                               | CRC_checksum                  | 20-23
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                               |page_segments  | segment_table | 24-27
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// | ...                                                           | 28-
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	capture_pattern := ogg_page_header[0 : 0+4]
	version := ogg_page_header[4]
	header_type := ogg_page_header[5]
	granule_position := binary.LittleEndian.Uint64(ogg_page_header[6 : 6+8])
	serial_number := binary.LittleEndian.Uint32(ogg_page_header[14 : 14+4])
	sequence_number := binary.LittleEndian.Uint32(ogg_page_header[18 : 18+4])
	crc_checksum := binary.LittleEndian.Uint32(ogg_page_header[22 : 22+4])
	page_segments := ogg_page_header[26]

	if !bytes.Equal(capture_pattern, []byte(ogg_capture_pattern)) {
		return page, errors.New("expected magic string OggS")
	}

	if version != 0 {
		return page, errors.New("expected version to be 0")
	}

	page.Continued = (header_type & 0x01) == 0x01
	page.FirstPage = (header_type & 0x02) == 0x02
	page.LastPage = (header_type & 0x04) == 0x04
	page.GranulePosition = int64(granule_position)
	page.SerialNumber = serial_number
	page.SequenceNumber = sequence_number
	page.Checksum = crc_checksum

	segment_table := make([]byte, page_segments)
	_, err = io.ReadFull(r, segment_table)
	if err != nil {
		return page, err
	}

	page.Complete = segment_table[len(segment_table)-1] < 255

	page_size := 0
	for _, lacing_value := range segment_table {
		page_size += int(lacing_value)
	}

	page.Body = make([]byte, page_size)
	_, err = io.ReadFull(r, page.Body)
	if err != nil {
		return page, err
	}

	return page, nil
}
