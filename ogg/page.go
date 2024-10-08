package ogg

import (
	"errors"
	"io"

	"github.com/steabert/gopus/binary"
)

const (
	ogg_page_header_size      = 27
	ogg_page_header_magic_sig = 0x5367674f // "OggS"
)

type Page struct {
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
func ParsePage(r io.Reader, page *Page) error {
	var err error

	br := binary.NewReader(r)

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
	capture_pattern := br.ReadUint32()
	version := br.ReadUint8()
	header_type := br.ReadUint8()
	granule_position := br.ReadUint64()
	serial_number := br.ReadUint32()
	sequence_number := br.ReadUint32()
	crc_checksum := br.ReadUint32()
	page_segments := br.ReadUint8()

	if br.Err() != nil {
		return br.Err()
	}

	if capture_pattern != ogg_page_header_magic_sig {
		return errors.New("expected magic string OggS")
	}

	if version != 0 {
		return errors.New("expected version to be 0")
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
		return err
	}

	page.Complete = segment_table[len(segment_table)-1] < 255

	page_size := 0
	for _, lacing_value := range segment_table {
		page_size += int(lacing_value)
	}

	page.Body = make([]byte, page_size)
	_, err = io.ReadFull(r, page.Body)
	if err != nil {
		return err
	}

	return nil
}
