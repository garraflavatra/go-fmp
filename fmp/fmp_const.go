package fmp

type FmpError string
type FmpChunkType uint8

func (e FmpError) Error() string { return string(e) }

var (
	ErrRead               = FmpError("read error")
	ErrBadMagic           = FmpError("bad magic number")
	ErrBadHeader          = FmpError("bad header")
	ErrUnsupportedCharset = FmpError("unsupported character set")
	ErrBadSectorCount     = FmpError("bad sector count")
	ErrBadSectorHeader    = FmpError("bad sector header")
	ErrBadChunk           = FmpError("bad chunk")
)

const (
	FMP_CHUNK_SIMPLE_DATA FmpChunkType = iota
	FMP_CHUNK_SEGMENTED_DATA
	FMP_CHUNK_SIMPLE_KEY_VALUE
	FMP_CHUNK_LONG_KEY_VALUE
	FMP_CHUNK_PATH_PUSH
	FMP_CHUNK_PATH_PUSH_LONG
	FMP_CHUNK_PATH_POP
	FMP_CHUNK_NOOP
)

const (
	FMP_COLLATION_ENGLISH     = 0x00
	FMP_COLLATION_FRENCH      = 0x01
	FMP_COLLATION_GERMAN      = 0x03
	FMP_COLLATION_ITALIAN     = 0x04
	FMP_COLLATION_DUTCH       = 0x05
	FMP_COLLATION_SWEDISH     = 0x07
	FMP_COLLATION_SPANISH     = 0x08
	FMP_COLLATION_DANISH      = 0x09
	FMP_COLLATION_PORTUGUESE  = 0x0A
	FMP_COLLATION_NORWEGIAN   = 0x0C
	FMP_COLLATION_FINNISH     = 0x11
	FMP_COLLATION_GREEK       = 0x14
	FMP_COLLATION_ICELANDIC   = 0x15
	FMP_COLLATION_TURKISH     = 0x18
	FMP_COLLATION_ROMANIAN    = 0x27
	FMP_COLLATION_POLISH      = 0x2a
	FMP_COLLATION_HUNGARIAN   = 0x2b
	FMP_COLLATION_RUSSIAN     = 0x31
	FMP_COLLATION_CZECH       = 0x38
	FMP_COLLATION_UKRAINIAN   = 0x3e
	FMP_COLLATION_CROATIAN    = 0x42
	FMP_COLLATION_CATALAN     = 0x49
	FMP_COLLATION_FINNISH_ALT = 0x62
	FMP_COLLATION_SWEDISH_ALT = 0x63
	FMP_COLLATION_GERMAN_ALT  = 0x64
	FMP_COLLATION_SPANISH_ALT = 0x65
	FMP_COLLATION_ASCII       = 0x66
)
