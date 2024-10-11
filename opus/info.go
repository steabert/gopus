package opus

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/steabert/gopus/binary"
	"github.com/steabert/gopus/ogg"
)

const (
	opus_id_header_size           = 19
	opus_id_header_magic_sig      = 0x646165487375704f // "OpusHead"
	opus_comment_header_size      = 16
	opus_comment_header_magic_sig = 0x736761547375704f // "OpusTags"
)

type OpusInfo struct {
	Vendor        string
	Comments      map[string]string
	SampleRate    uint32
	PreSkip       uint16
	OutputGain    float64
	Channels      uint8
	MappingFamily uint8
}

func ParseInfo(path string) (OpusInfo, error) {
	var info OpusInfo

	f, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	// Parse the identification header. This is a single-page header
	// at the "Beginning Of Stream" that has to be complete.

	var page ogg.Page
	err = ogg.ParsePage(r, &page)
	if err != nil {
		return info, fmt.Errorf("invalid OGG stream, %v", err)
	}

	if !page.FirstPage || !page.Complete {
		return info, fmt.Errorf("invalid identification header page, %v", err)
	}

	err = parseIDHeader(bytes.NewReader(page.Body), &info)
	if err != nil {
		return info, fmt.Errorf("invalid identification header, %v", err)
	}

	// Parse the comment header. This can span multiple pages.

	var commentHeaderPages []io.Reader
	for {
		var page ogg.Page
		err := ogg.ParsePage(r, &page)
		if err != nil {
			return info, fmt.Errorf("invalid OGG stream, %v", err)
		}

		commentHeaderPages = append(commentHeaderPages, bytes.NewReader(page.Body))

		if page.Complete {
			break
		}
	}

	err = parseCommentHeader(io.MultiReader(commentHeaderPages...), &info)
	if err != nil {
		return info, fmt.Errorf("invalid identification header, %v", err)
	}

	return info, nil
}

// parseHeader parses an Opus identification (ID) header.
// The reader's underlying data _must_ be able to support
// the mininum size of a header, otherwise the parser will
// panic.
func parseIDHeader(r io.Reader, info *OpusInfo) (err error) {
	br := binary.NewReader(r)

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
	capture_pattern := br.ReadUint64()
	version := br.ReadUint8()
	channel_count := br.ReadUint8()
	pre_skip := br.ReadUint16()
	sample_rate := br.ReadUint32()
	output_gain := br.ReadUint16()
	mapping_family := br.ReadUint8()

	if br.Err() != nil {
		return br.Err()
	}

	if capture_pattern != opus_id_header_magic_sig {
		return errors.New("expected magic signature 'OpusHead'")
	}

	if version != 1 {
		return errors.New("expected version=1")
	}

	info.Channels = channel_count
	info.PreSkip = pre_skip
	info.SampleRate = sample_rate
	info.OutputGain = float64(int16(output_gain)) / float64(256.0)
	info.MappingFamily = mapping_family

	return nil
}

func parseCommentHeader(r io.Reader, info *OpusInfo) error {
	br := binary.NewReader(r)

	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |      'O'      |      'p'      |      'u'      |      's'      |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |      'T'      |      'a'      |      'g'      |      's'      |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                     Vendor String Length                      |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                                                               |
	// :                        Vendor String...                       :
	// |                                                               |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                   User Comment List Length                    |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                 User Comment #0 String Length                 |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                                                               |
	// :                   User Comment #0 String...                   :
	// |                                                               |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                 User Comment #1 String Length                 |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// :                                                               :

	capture_pattern := br.ReadUint64()
	if br.Err() != nil {
		return br.Err()
	}

	if capture_pattern != opus_comment_header_magic_sig {
		return errors.New("expected magic signature 'OpusHead'")
	}

	vendor_string_length := br.ReadUint32()
	if br.Err() != nil {
		return br.Err()
	}

	vendor_string := make([]byte, vendor_string_length)
	_, err := io.ReadFull(r, vendor_string)
	if err != nil {
		return err
	}

	user_comment_list_length := br.ReadUint32()
	if br.Err() != nil {
		return br.Err()
	}

	user_comments := make(map[string]string, user_comment_list_length)
	for range user_comment_list_length {
		user_comment_string_length := br.ReadUint32()
		if br.Err() != nil {
			return br.Err()
		}

		user_comment_string := make([]byte, user_comment_string_length)
		_, err := io.ReadFull(r, user_comment_string)
		if err != nil {
			return err
		}
		key, value, found := bytes.Cut(user_comment_string, []byte("="))
		if !found {
			continue
		}
		user_comments[strings.ToUpper(string(key))] = string(value)
	}

	info.Vendor = string(vendor_string)
	info.Comments = user_comments

	return nil
}
