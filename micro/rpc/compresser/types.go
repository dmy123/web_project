package compresser

type Compresser interface {
	Code() byte
	Compress([]byte) ([]byte, error)
	Uncompress([]byte) ([]byte, error)
}

type DoNothingCompresser struct {
}

func (d DoNothingCompresser) Code() byte {
	return 0
}

func (d DoNothingCompresser) Compress(bytes []byte) ([]byte, error) {
	return bytes, nil
}

func (d DoNothingCompresser) Uncompress(bytes []byte) ([]byte, error) {
	return bytes, nil
}
