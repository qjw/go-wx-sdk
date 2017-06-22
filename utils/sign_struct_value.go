package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"reflect"
	"sort"
	"strconv"
	"crypto/sha1"
	"log"
)

type SignStructValue struct {
	Keys     []string          `json:"keys"`
}

func (this SignStructValue) isValidType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8:
		return true
	case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		return true
	case reflect.Ptr:
		return true
	}
	return false
}

func (this *SignStructValue) fetchElem(tp reflect.StructField, value reflect.Value) {
	// 排除空值
	if value.Interface() == reflect.Zero(value.Type()).Interface() {
		return
	}

	var varStr string
	switch value.Type().Kind() {
	case reflect.Int:
		varStr = strconv.FormatInt(value.Int(), 10)
	case reflect.Int8:
		varStr = strconv.FormatInt(value.Int(), 10)
	case reflect.Int16:
		varStr = strconv.FormatInt(value.Int(), 10)
	case reflect.Int32:
		varStr = strconv.FormatInt(value.Int(), 10)
	case reflect.Int64:
		varStr = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		varStr = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint8:
		varStr = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint16:
		varStr = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint32:
		varStr = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint64:
		varStr = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		varStr = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		varStr = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		varStr = value.String()
	case reflect.Bool:
		varStr = strconv.FormatBool(value.Bool())
	default:
		panic("invalid type")
	}
	this.Keys = append(this.Keys, varStr)
}

func (this *SignStructValue) genSignData(variable interface{}) (err error) {
	value := reflect.ValueOf(variable)
	if value.Kind() != reflect.Ptr {
		err = errors.New("invalid varibal to sign (need tobe pointer)")
		return
	}
	value = value.Elem()
	// 分配Map
	this.Keys = make([]string, 0)
	return this.genSignDataImp(value)
}

func (this *SignStructValue) genSignDataImp(value reflect.Value) (err error) {
	if value.Kind() != reflect.Struct {
		err = errors.New("invalid varibal to sign (need tobe struct)")
		return
	}
	tp := value.Type()
	count := tp.NumField()
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		field := tp.Field(i)

		if field.Anonymous {
			this.genSignDataImp(value.Field(i))
			continue
		}

		if "-" == field.Tag.Get("sign") {
			continue
		}

		fieldTp := field.Type
		if !this.isValidType(fieldTp.Kind()) {
			err = fmt.Errorf("unsupport type %s", fieldTp.String())
			return
		}

		if fieldTp.Kind() == reflect.Ptr {
			tmpVar := value.Field(i)
			if tmpVar.IsNil() {
				continue
			}
			this.fetchElem(field, tmpVar.Elem())

		} else {
			this.fetchElem(field, value.Field(i))
		}
	}
	return
}

func (this SignStructValue) Sign(variable interface{}, fn func() hash.Hash) (sign string, err error) {
	if err = this.genSignData(variable); err != nil {
		return
	}
	sort.Strings(this.Keys)
	if fn == nil {
		fn = sha1.New
	}
	h := fn()
	for _, k := range this.Keys {
		log.Printf("key: %s",k)
		io.WriteString(h, k)
	}

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature)), nil

}
