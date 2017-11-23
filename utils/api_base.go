package utils

import (
	"reflect"
	"fmt"
	"encoding/json"
	"errors"
	"net/http"
)

type ApiBase struct {
}

func (this ApiBase) DoPost(uri_pattern string, body string, res interface{}, a ...interface{}) error {
	return this.DoPostRaw(uri_pattern, []byte(body), res, a...)
}

func (this ApiBase) DoPostObject(uri_pattern string, body interface{}, res interface{}, a ...interface{}) error {
	tp := reflect.TypeOf(body)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Struct{
		panic("invalid body object type")
	}

	raw,err := json.Marshal(body)
	if err != nil{
		return err
	}
	return this.DoPostRaw(uri_pattern, raw, res, a...)
}

func (this ApiBase) DoPostRaw(uri_pattern string, body []byte, res interface{}, a ...interface{}) (err error) {
	var uri string
	if len(a) == 0 {
		uri = uri_pattern
	} else {
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := PostJSONRaw(uri, body)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiBase) DoGet(uri_pattern string, res interface{}, a ...interface{}) (err error) {
	return this.doHttpByUri(HTTPGet,uri_pattern,res,a...)
}

func (this ApiBase) DoDelete(uri_pattern string, res interface{}, a ...interface{}) (err error) {
	return this.doHttpByUri(HTTPDelete,uri_pattern,res,a...)
}

func (this ApiBase) doHttpByUri(httpProcessor func (string)([]byte, int, error),
		uri_pattern string, res interface{}, a ...interface{}) (err error) {
	var uri string
	if len(a) == 0 {
		uri = uri_pattern
	} else {
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code,  err := httpProcessor(uri)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiBase) doProcessResponse(response []byte, code int,  err error, res interface{}) error{
	if err != nil {
		return err
	}

	if code == http.StatusNoContent && len(response) == 0{
		response = []byte("{}")
	}

	err = json.Unmarshal(response, res)
	if err != nil {
		return errors.New(string(response))
	}
	return nil
}