package booster

import (
	"encoding/base64"
	"errors"
	"go.uber.org/zap/buffer"
	"io/ioutil"
	"os"
	"strings"
)

var errIsDirectory = errors.New("destination path is a directory")

func readFile(path string) ([]byte, error) {
	openFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := openFile.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		data, err := ioutil.ReadAll(openFile)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errIsDirectory
}

func toB64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded

}
func BytesToString(b []byte) string {
	builder := strings.Builder{}
	builder.Write(b)
	return builder.String()
}
func StrToBytes(s string) []byte {
	var buff buffer.Buffer
	buff.AppendString(s)
	return buff.Bytes()
}

func DecodeBytesB64(b string) (string, error) {
	//var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	l, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return "", err
	}
	return BytesToString(l), nil
}
