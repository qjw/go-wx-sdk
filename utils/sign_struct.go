package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type SignStruct struct {
	Elements map[string]string `json:"elements"`
	Keys     []string          `json:"keys"`
	ToLower  bool              `doc:"是否自动将key小写"`
	Tag      string            `json:"tag" doc:"使用xml/json或者直接字段名"`
}

func (this SignStruct) isValidType(kind reflect.Kind) bool {
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

func (this *SignStruct) fetchElem(tp reflect.StructField, value reflect.Value, tag string) {
	name := tag
	if len(name) < 1 {
		name = tp.Name
	}
	if this.ToLower {
		name = strings.ToLower(name)
	}

	// 排除空值
	if value.Interface() == reflect.Zero(value.Type()).Interface() {
		return
	}

	switch value.Type().Kind() {
	case reflect.Int:
		this.Elements[name] = strconv.FormatInt(value.Int(), 10)
	case reflect.Int8:
		this.Elements[name] = strconv.FormatInt(value.Int(), 10)
	case reflect.Int16:
		this.Elements[name] = strconv.FormatInt(value.Int(), 10)
	case reflect.Int32:
		this.Elements[name] = strconv.FormatInt(value.Int(), 10)
	case reflect.Int64:
		this.Elements[name] = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		this.Elements[name] = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint8:
		this.Elements[name] = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint16:
		this.Elements[name] = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint32:
		this.Elements[name] = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint64:
		this.Elements[name] = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		this.Elements[name] = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		this.Elements[name] = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		this.Elements[name] = value.String()
	case reflect.Bool:
		this.Elements[name] = strconv.FormatBool(value.Bool())
	default:
		panic("invalid type")
	}
	this.Keys = append(this.Keys, name)
}

func (this *SignStruct) genSignData(variable interface{}) (err error) {
	value := reflect.ValueOf(variable)
	if value.Kind() != reflect.Ptr {
		err = errors.New("invalid varibal to sign (need tobe pointer)")
		return
	}
	value = value.Elem()
	// 分配Map
	this.Elements = make(map[string]string)
	this.Keys = make([]string, 0)
	return this.genSignDataImp(value)
}

func (this SignStruct) getTag(field *reflect.StructField) (tag string, tags []string) {
	if strings.ToLower(this.Tag) == "xml" {
		tag = field.Tag.Get("xml")
	} else if strings.ToLower(this.Tag) == "json" {
		tag = field.Tag.Get("json")
	}

	tags = strings.Split(tag, ",")
	if len(tags) > 0 {
		tag = tags[0]
	}
	return
}

func (this *SignStruct) genSignDataImp(value reflect.Value) (err error) {
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
		tag, _ := this.getTag(&field)
		if tag == "-" {
			tag = ""
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
			this.fetchElem(field, tmpVar.Elem(), tag)

		} else {
			this.fetchElem(field, value.Field(i), tag)
		}
	}
	return
}

func (this SignStruct) Sign(variable interface{}, fn func() hash.Hash, apiKey string) (sign string, err error) {
	if err = this.genSignData(variable); err != nil {
		return
	}
	sort.Strings(this.Keys)
	if fn == nil {
		fn = md5.New
	}
	h := fn()

	firstFlag := true
	for _, k := range this.Keys {
		if !firstFlag {
			io.WriteString(h, "&")
		} else {
			firstFlag = false
		}
		v := this.Elements[k]

		log.Printf("key %s value %s\n", k, v)

		io.WriteString(h, k)
		io.WriteString(h, "=")
		io.WriteString(h, v)
	}

	if len(apiKey) > 0 {
		log.Printf("key %s value %s\n", "key", apiKey)
		io.WriteString(h, "&key=")
		io.WriteString(h, apiKey)
	}

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature)), nil

}
