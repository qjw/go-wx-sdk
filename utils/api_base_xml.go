package utils

import (
	"reflect"
	"fmt"
	"log"
	"errors"
	"encoding/xml"
	"net/http"
)

type ApiBaseXml struct {
}

func (this ApiBaseXml) DoPost(uri_pattern string, body string, res interface{}, a ...interface{}) error {
	return this.DoPostRaw(uri_pattern, []byte(body), res, a...)
}

func (this ApiBaseXml) DoPostObject(uri_pattern string, body interface{}, res interface{}, a ...interface{}) error {
	tp := reflect.TypeOf(body)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Struct{
		panic("invalid body object type")
	}

	raw,err := xml.Marshal(body)
	if err != nil{
		return err
	}
	return this.DoPostRaw(uri_pattern, raw, res, a...)
}

func (this ApiBaseXml) DoPostRaw(uri_pattern string, body []byte, res interface{}, a ...interface{}) (err error) {
	var uri string
	if len(a) == 0 {
		uri = uri_pattern
	} else {
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := PostJSONRaw(uri, body)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiBaseXml) DoGet(uri_pattern string, res interface{}, a ...interface{}) (err error) {
	var uri string
	if len(a) == 0 {
		uri = uri_pattern
	} else {
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := HTTPGet(uri)
	return this.doProcessResponse(response,code,err,res)
}


func (this ApiBaseXml) doProcessResponse(response []byte, code int,  err error, res interface{}) error{
	if err != nil {
		return err
	}

	tmp := string(response)
	log.Print(tmp)

	if code == http.StatusNoContent && len(response) == 0{
		response = []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	}

	err = xml.Unmarshal(response, res)
	if err != nil {
		return errors.New(string(response))
	}
	return nil
}