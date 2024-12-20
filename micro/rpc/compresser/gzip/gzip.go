package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Compresser struct {
}

func (c Compresser) Code() byte {
	return 1
}

func (c Compresser) Compress(data []byte) ([]byte, error) {
	res := &bytes.Buffer{}
	gw := gzip.NewWriter(res)
	_, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	// 这个地方不能使用 defer，一定自己手动的调用 Close
	// 否则部分数据还没刷新到 res 里面，
	// 这是一个非常容易出错的地方
	if err = gw.Close(); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func (c Compresser) Uncompress(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	return io.ReadAll(gr)
}
