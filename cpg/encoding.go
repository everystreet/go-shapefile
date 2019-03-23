package cpg

type CharacterEncoding uint

const (
	// EncodingUnknown represents any unknown or unsupported character encodings.
	EncodingUnknown CharacterEncoding = iota
	// EncodingASCII is ASCII encoding.
	EncodingASCII
	// EncodingUTF8 is UTF-8 encoding.
	EncodingUTF8
)

// String name of encoding.
func (e CharacterEncoding) String() string {
	switch e {
	case EncodingASCII:
		return "ASCII"
	case EncodingUTF8:
		return "UTF-8"
	case EncodingUnknown:
		fallthrough
	default:
		return "unknown"
	}
}
