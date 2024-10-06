package opus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/steabert/gopus/ogg"
)

type OpusInfo struct {
	Channels      uint8
	PreSkip       uint16
	SampleRate    uint32
	OutputGain    int16
	MappingFamily uint8
	Title         string
}

func ParseInfo(path string) (OpusInfo, error) {
	var metadata OpusInfo

	f, err := os.Open(path)
	if err != nil {
		return metadata, err
	}
	defer f.Close()

	page, err := ogg.ParsePage(f)
	if err != nil {
		return metadata, fmt.Errorf("invalid OGG stream, %v", err)
	}
	fmt.Printf("page: %+v\n", page)

	err = parseHeader(page.Body, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("invalid OPUS file, %v", err)
	}

	return metadata, nil
}

const (
	opus_header_size       = 19
	opus_header_magic_sig  = "OpusHead"
	opus_comment_magic_sig = "OpusTags"
)

// parseHeader parses an Opus identification header.
//
//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      'O'      |      'p'      |      'u'      |      's'      | 0-3
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      'H'      |      'e'      |      'a'      |      'd'      | 4-7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  Version = 1  | Channel Count |           Pre-skip            | 8-11
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     Input Sample Rate (Hz)                    | 12-15
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Output Gain (Q7.8 in dB)    | Mapping Family|               | 16-19
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               :
// |                                                               |
// :               Optional Channel Mapping Table...               :
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func parseHeader(buf []byte, metadata *OpusInfo) error {
	if len(buf) < opus_header_size {
		return errors.New("wrong header size")
	}

	capture_pattern := buf[0 : 0+8]
	version := buf[8]
	channel_count := buf[9]
	pre_skip := binary.LittleEndian.Uint16(buf[10 : 10+2])
	sample_rate := binary.LittleEndian.Uint32(buf[12 : 12+4])
	output_gain := int16(binary.LittleEndian.Uint16(buf[16 : 16+2]))
	mapping_family := buf[18]

	if !bytes.Equal(capture_pattern, []byte(opus_header_magic_sig)) {
		return errors.New("expected magic signature 'OpusHead'")
	}

	if version != 1 {
		return errors.New("expected version=1")
	}

	metadata.Channels = channel_count
	metadata.PreSkip = pre_skip
	metadata.SampleRate = sample_rate
	metadata.OutputGain = output_gain
	metadata.MappingFamily = mapping_family

	return nil
}
