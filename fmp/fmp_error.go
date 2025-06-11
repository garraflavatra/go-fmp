package fmp

type FmpError string

func (e FmpError) Error() string { return string(e) }

var (
	ErrRead               = FmpError("read error")
	ErrBadMagic           = FmpError("bad magic number")
	ErrBadHeader          = FmpError("bad header")
	ErrUnsupportedCharset = FmpError("unsupported character set")
	ErrBadSectorCount     = FmpError("bad sector count")
	ErrNoFmemopen         = FmpError("no fmemopen")
	ErrOpen               = FmpError("could not open file")
	ErrSeek               = FmpError("seek failed")
	ErrMalloc             = FmpError("malloc failed")
	ErrUserAborted        = FmpError("user aborted")
)
