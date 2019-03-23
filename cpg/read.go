package cpg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Read a .cpg file containing a character encoding.
func Read(r io.Reader) (CharacterEncoding, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 {
			continue
		}

		s = strings.ToUpper(s)
		switch s {
		case "ASCII":
			return EncodingASCII, nil
		case "UTF8":
			fallthrough
		case "UTF-8":
			return EncodingUTF8, nil
		default:
			return EncodingUnknown, nil
		}
	}
	return EncodingUnknown, fmt.Errorf("invalid format")
}
