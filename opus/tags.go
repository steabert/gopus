package opus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/steabert/gopus/ogg"
)

type OpusMetadata struct {
	Channels   uint8
	SampleRate uint32
	OutputGain int16
	Title      string
}

func Parse(path string) (OpusMetadata, error) {
	var metadata OpusMetadata

	f, err := os.Open(path)
	if err != nil {
		return metadata, err
	}
	defer f.Close()

	bos, err := ogg.ReadBOSPage(f)
	if err != nil {
		return metadata, fmt.Errorf("invalid OGG stream, %v", err)
	}

	err = parseOpusHeader(bos, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("invalid OPUS file, %v", err)
	}

	// bos, err := ogg.Read(f)
	// if err != nil {
	// 	return tags, fmt.Errorf("invalid OGG stream, %v", err)
	// }

	return metadata, nil
}

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      'O'      |      'p'      |      'u'      |      's'      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      'H'      |      'e'      |      'a'      |      'd'      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  Version = 1  | Channel Count |           Pre-skip            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     Input Sample Rate (Hz)                    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Output Gain (Q7.8 in dB)    | Mapping Family|               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               :
// |                                                               |
// :               Optional Channel Mapping Table...               :
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

const (
	opus_required_header_size = 19
	opus_magic_signature_size = 8
	opus_channel_count_offset = 9
	opus_sample_rate_offset   = 12
	opus_output_gain_offset   = 16
)

func parseOpusHeader(bos []byte, metadata *OpusMetadata) error {
	if len(bos) < opus_required_header_size {
		return errors.New("wrong header size")
	}

	if !bytes.Equal(bos[:opus_magic_signature_size], []byte("OpusHead")) {
		return errors.New("wrong magic signature")
	}

	metadata.Channels = uint8(bos[opus_channel_count_offset])
	metadata.SampleRate = binary.LittleEndian.Uint32(bos[opus_sample_rate_offset : opus_sample_rate_offset+4])
	metadata.OutputGain = int16(binary.LittleEndian.Uint16(bos[opus_output_gain_offset : opus_output_gain_offset+2]))

	return nil
}
