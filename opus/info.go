package opus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
	OutputGain    int16
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

	page, err := ogg.ParsePage(r)
	if err != nil {
		return info, fmt.Errorf("invalid OGG stream, %v", err)
	}

	if !page.FirstPage || !page.Complete {
		return info, fmt.Errorf("invalid identification header page, %v", err)
	}

	fmt.Printf("%d, %+v\n", len(page.Body), page.Body)
	err = parseIDHeader(bytes.NewReader(page.Body), &info)
	if err != nil {
		return info, fmt.Errorf("invalid identification header, %v", err)
	}

	// Parse the comment header. This can span multiple pages.

	var commentHeader []io.Reader
	for {
		page, err := ogg.ParsePage(r)
		if err != nil {
			return info, fmt.Errorf("invalid OGG stream, %v", err)
		}

		commentHeader = append(commentHeader, bytes.NewReader(page.Body))

		if page.Complete {
			break
		}
	}

	err = parseCommentHeader(io.MultiReader(commentHeader...), &info)
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
	capture_pattern := mustReadUint64(r)
	version := mustReadUint8(r)
	channel_count := mustReadUint8(r)
	pre_skip := mustReadUint16(r)
	sample_rate := mustReadUint32(r)
	output_gain := mustReadUint16(r)
	mapping_family := mustReadUint8(r)

	if capture_pattern != opus_id_header_magic_sig {
		return errors.New("expected magic signature 'OpusHead'")
	}

	if version != 1 {
		return errors.New("expected version=1")
	}

	info.Channels = channel_count
	info.PreSkip = pre_skip
	info.SampleRate = sample_rate
	info.OutputGain = int16(output_gain)
	info.MappingFamily = mapping_family

	return nil
}

func parseCommentHeader(r io.Reader, info *OpusInfo) error {
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

	capture_pattern := mustReadUint64(r)

	if capture_pattern != opus_comment_header_magic_sig {
		return errors.New("expected magic signature 'OpusHead'")
	}

	vendor_string_length := mustReadUint32(r)
	fmt.Println("vendor", vendor_string_length)
	vendor_string := make([]byte, vendor_string_length)
	_, err := io.ReadFull(r, vendor_string)
	if err != nil {
		return err
	}

	user_comment_list_length := mustReadUint32(r)
	user_comments := make(map[string]string, user_comment_list_length)
	for range user_comment_list_length {
		user_comment_string_length := mustReadUint32(r)
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
	fmt.Println(user_comments)
	info.Comments = user_comments

	return nil
}

func mustReadUint8(r io.Reader) uint8 {
	b := [1]byte{}
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	return uint8(b[0])
}

func mustReadUint16(r io.Reader) uint16 {
	b := [2]byte{}
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint16(b[:])
}

func mustReadUint32(r io.Reader) uint32 {
	b := [4]byte{}
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint32(b[:])
}

func mustReadUint64(r io.Reader) uint64 {
	b := [8]byte{}
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(b[:])
}
