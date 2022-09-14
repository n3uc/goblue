package bluefile

import (
	"encoding/json"
	"goblue/pkg/blueheaders"
	"os"
)

type BlueFile struct {
	FileName string                  `json:"filename"`
	Header   *blueheaders.FullHeader `json:"header"`
}

// ----------------------------------------------------------------------------
// Public Functions
// ----------------------------------------------------------------------------
func (h *BlueFile) String() string {
	s, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(s)
}

// New Returns a Blue File
func New(filename string) (*BlueFile, error) {
	bf := BlueFile{FileName: filename}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h, err := blueheaders.New(file)
	if err != nil {
		return nil, err
	}

	bf.Header = h

	return &bf, nil
}
