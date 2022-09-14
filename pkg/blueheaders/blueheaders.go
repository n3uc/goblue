package blueheaders

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	MIDAS_BLUE          = []byte{0x42, 0x4C, 0x55, 0x45} // BLUE
	MIDAS_BIG_ENDIAN    = []byte{0x49, 0x45, 0x45, 0x45} // IEEE
	MIDAS_LITTLE_ENDIAN = []byte{0x45, 0x45, 0x45, 0x49} // EEEI
)

type header struct {
	Version    [4]byte    `json:"-"`
	Head_rep   [4]byte    `json:"-"`
	Data_rep   [4]byte    `json:"-"`
	Detached   uint32     `json:"detached"`
	Protected  uint32     `json:"protected"`
	Pipe       uint32     `json:"pipe"`
	Ext_start  uint32     `json:"ext_start"`
	Ext_size   uint32     `json:"ext_size"`
	Data_start float64    `json:"data_start"`
	Data_size  float64    `json:"data_size"`
	Ftype      uint32     `json:"type"`
	Format     [2]byte    `json:"-"`
	Flagmask   uint16     `json:"flagmask"`
	Timecode   float64    `json:"timecode"`
	Inlet      [2]byte    `json:"-"`
	Outlets    [2]byte    `json:"-"`
	Outmask    uint32     `json:"outmask"`
	Pipeloc    uint32     `json:"pipeloc"`
	Pipesize   uint32     `json:"pipesize"`
	In_byte    float64    `json:"in_byte"`
	Out_byte   float64    `json:"out_byte"`
	Outbytes   [8]float64 `json:"out_bytes"`
	Keylength  uint32     `json:"keylength"`
	Keywords   [92]byte   `json:"-"`
	Adjunct    [256]byte  `json:"adjunct"`
}

// FullHeader holds the read in entire header
type FullHeader struct {
	Versions  string `json:"version"`
	Head_reps string `json:"head_rep"`
	Data_reps string `json:"data_rep"`
	header
	Formats         string            `json:"format"`
	Inlete          uint16            `json:"inlet"`
	Outletse        uint16            `json:"outlets"`
	Keywordsp       map[string]string `json:"keywords"`
	ExtendedHeaders map[string]string `json:"extendedHeaders"`
}

type keywordLoc struct {
	Lkey  int32
	Lext  int16
	Ltag  int8
	Dtype byte
}

var (
	errBlue = errors.New("file is not a blue file")
)

// CheckEndian Returns Midas endianness of a byte slice
func checkEndian(head_rep []byte) (binary.ByteOrder, error) {
	if bytes.Equal(head_rep, MIDAS_BIG_ENDIAN) {
		return binary.BigEndian, nil
	}
	if bytes.Equal(head_rep, MIDAS_LITTLE_ENDIAN) {
		return binary.LittleEndian, nil
	}
	return binary.LittleEndian, errors.New("invalid endianness")
}

func loadHeader(file *os.File) (header, error) {
	// Check for endianness and load headers accordingly
	header := header{}
	lookbuf := make([]byte, 4)

	_, err := file.ReadAt(lookbuf, 4)
	if err != nil {
		return header, err
	}

	end, err := checkEndian(lookbuf)
	if err != nil {
		return header, err
	}

	return header, binary.Read(file, end, &header)
}

// Creates a map containing a bluefile's extended headers
func loadExtHeader(file *os.File, header header) (map[string]string, error) {
	// read in extended header
	extbuf := make([]byte, header.Ext_size)
	rMap := make(map[string]string)

	// ext_start determines how many 512 byte sized blocks from the start the extended headers are
	_, err := file.ReadAt(extbuf, int64(header.Ext_start)*512)

	if err != nil {
		return rMap, err
	}

	end, err := checkEndian(header.Head_rep[:])
	if err != nil {
		return rMap, err
	}

	var offset int64 = 0
	for int(offset) < int(header.Ext_size) {
		var key string
		var value string
		offset, key, value, err = parseKeyword(extbuf, offset, end)
		if err != nil {
			return rMap, err
		}
		rMap[key] = value
	}

	return rMap, err
}

// Parses a single extended header keyword and returns the new offsetinto the keyword buffer
func parseKeyword(buf []byte, offset int64, end binary.ByteOrder) (int64, string, string, error) {
	bufRead := bytes.NewReader(buf)
	bufRead.Seek(offset, 0)

	kl := keywordLoc{}
	err := binary.Read(bufRead, end, &kl)
	if err != nil {
		return 0, "", "", err
	}

	var val string
	switch kl.Dtype {
	case 'A': // ASCII
		tval := make([]byte, kl.Lkey-int32(kl.Lext))
		err = binary.Read(bufRead, end, &tval)
		val = string(tval)
	case 'B': // 8-bit integer
		var tval int8
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%d", tval)
	case 'I': // 16-bit integer
		var tval int16
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%d", tval)
	case 'L': // 32-bit integer
		var tval int32
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%d", tval)
	case 'X': // 64-bit integer
		var tval int64
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%d", tval)
	case 'F': // 32-bit float
		var tval float32
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%f", tval)
	case 'D': // 64-bit float
		var tval float64
		err = binary.Read(bufRead, end, &tval)
		val = fmt.Sprintf("%f", tval)
	default:
		return 0, "", "", errors.New("invalid keyword type")
	}
	if err != nil {
		return 0, "", "", err
	}

	tag := make([]byte, kl.Ltag)

	_, err = bufRead.ReadAt(tag, offset+8+int64(kl.Lkey)-int64(kl.Lext))
	if err != nil {
		return 0, "", "", err
	}

	offset = offset + int64(kl.Lkey)
	return offset, string(tag), val, nil
}

// Parses the fixed header keywords and returns it as a map of key value pairs
func parseFixedKeywords(keywords []byte) map[string]string {
	rMap := make(map[string]string)
	pairs := strings.Split(string(keywords), "\x00")
	for _, p := range pairs {
		parts := strings.Split(p, "=")
		if len(parts) == 2 {
			rMap[parts[0]] = parts[1]
		}
	}
	return rMap
}

// LoadFullHeader Creates a fullheader struct from a given file
func loadFullHeader(file *os.File) (FullHeader, error) {
	header, err := loadHeader(file)
	if err != nil {
		return FullHeader{}, err
	}

	extHeader, err := loadExtHeader(file, header)
	if err != nil {
		return FullHeader{}, err
	}

	end, err := checkEndian(header.Head_rep[:])
	if err != nil {
		return FullHeader{}, err
	}

	rheader := FullHeader{
		header:          header,
		Versions:        string(header.Version[:]),
		Head_reps:       string(header.Head_rep[:]),
		Data_reps:       string(header.Data_rep[:]),
		Formats:         string(header.Format[:]),
		Inlete:          end.Uint16(header.Inlet[:]),
		Outletse:        end.Uint16(header.Outlets[:]),
		Keywordsp:       parseFixedKeywords(header.Keywords[:header.Keylength]),
		ExtendedHeaders: extHeader}

	return rheader, nil
}

// CheckBlue Checks the header to validate the version of a given file
func checkBlue(file *os.File) error {
	// Check if header is for a BLUE file
	lookbuf := make([]byte, 4)

	_, err := file.ReadAt(lookbuf, 0)
	if err != nil {
		return err
	}

	if !bytes.Equal(lookbuf, MIDAS_BLUE) {
		return errBlue
	}

	return nil
}

// ----------------------------------------------------------------------------
// Public Functions
// ----------------------------------------------------------------------------
func (h *FullHeader) String() string {
	s, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(s)
}

// New Returns a Full Header
func New(file *os.File) (*FullHeader, error) {
	// Check if file is a bluefile
	err := checkBlue(file)
	if err != nil {
		return nil, err
	}

	h, err := loadFullHeader(file)
	if err != nil {
		return nil, err
	}

	return &h, nil
}
