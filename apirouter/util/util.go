package util

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"reflect"
	"unsafe"
)

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
func UnmarshalFromFile(dst interface{}, filename, objPath string) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file '%s', error: %s", filename, err.Error())
	}

	if objPath != "" {
		err = jsoniter.UnmarshalFromString(gjson.GetBytes(data, objPath).String(), dst)
	} else {
		err = jsoniter.Unmarshal(data, dst)
	}
	if err != nil {
		return err
	}

	return nil
}
