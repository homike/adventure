package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"strings"
)

const (
	BZCompressTag = "bz"
)

// zlib+base64,
func CompressStringBZ(is string) (string, error) {

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(is)); err != nil {
		return "", err
	}
	gz.Close()

	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	s += BZCompressTag
	return s, nil
}

// zlib+base64, the is must have bz as suffix
func DecompressStringBZ(is string) (string, error) {

	if strings.HasSuffix(is, BZCompressTag) {
		//log.Printf("DecompressStringBZ, is %v \n ", is)
		origins := strings.TrimSuffix(is, BZCompressTag)
		bs, err := base64.StdEncoding.DecodeString(origins)
		if err != nil {
			return "", err
		}
		//log.Printf("DecompressStringBZ, bs %v \n", bs)
		bsr := bytes.NewReader(bs)
		gz, err := gzip.NewReader(bsr)
		if err != nil {
			return "", err
		}
		jb, err := ioutil.ReadAll(gz)
		gz.Close()
		if err != nil {
			return "", err
		}
		//log.Printf("DecompressStringBZ, jb %v \n", jb)
		return string(jb), nil

	} else {
		return is, nil
	}
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

func Int2Byte(data int) []byte {
	b := []byte{0x00, 0x00, 0x03, 0xe8}
	b_buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	return b_buf.Bytes()
}
