package utils

import (
	"fmt"
	"log"
	"io"
	"encoding/json"
	"errors"
	"reflect"
	"net/http"
)

type ContextToken interface{
	GetAccessToken() (accessToken string, err error);
	GetJsTicket() (jsTicket string, err error)
}

type ApiTokenBase struct {
	ContextToken ContextToken
}

func (this ApiTokenBase) DoPost(uri_pattern string, body string, res interface{}, a ...interface{}) error {
	return this.DoPostRaw(uri_pattern, []byte(body), res, a...)
}

func (this ApiTokenBase) DoPostObject(uri_pattern string, body interface{}, res interface{}, a ...interface{}) error {
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

func (this ApiTokenBase) DoPostRaw(uri_pattern string, body []byte, res interface{}, a ...interface{}) error {
	accessToken, err := this.ContextToken.GetAccessToken()
	if err != nil {
		return err
	}

	var uri string
	if len(a) == 0 {
		uri = fmt.Sprintf(uri_pattern, accessToken)
	} else {
		//todo 性能
		a = append([]interface{}{accessToken}, a...)
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := PostJSONRaw(uri, body)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiTokenBase) DoGet(uri_pattern string, res interface{}, a ...interface{}) error {
	accessToken, err := this.ContextToken.GetAccessToken()
	if err != nil {
		return err
	}

	var uri string
	if len(a) == 0 {
		uri = fmt.Sprintf(uri_pattern, accessToken)
	} else {
		//todo 性能
		a = append([]interface{}{accessToken}, a...)
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := HTTPGet(uri)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiTokenBase) DoPostFile(reader io.Reader, fieldname, filename string, res interface{},
	uri_pattern string, a ...interface{}) error {
	accessToken, err := this.ContextToken.GetAccessToken()
	if err != nil {
		return err
	}

	var uri string
	if len(a) == 0 {
		uri = fmt.Sprintf(uri_pattern, accessToken)
	} else {
		//todo 性能
		a = append([]interface{}{accessToken}, a...)
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	response, code, err := PostReader(fieldname, filename, uri, reader)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiTokenBase) DoPostFileExtra(reader io.Reader, fieldname, extraFieldname, filename string,
	extra interface{}, res interface{},
	uri_pattern string, a ...interface{}) error {
	accessToken, err := this.ContextToken.GetAccessToken()
	if err != nil {
		return err
	}

	descBytes, err := json.Marshal(extra)
	if err != nil {
		return err
	}

	var uri string
	if len(a) == 0 {
		uri = fmt.Sprintf(uri_pattern, accessToken)
	} else {
		//todo 性能
		a = append([]interface{}{accessToken}, a...)
		uri = fmt.Sprintf(uri_pattern, a...)
	}

	fields := []MultipartFormField{
		{
			IsFile:    false,
			Fieldname: fieldname,
			Filename:  filename,
			Reader:    reader,
		},
		{
			IsFile:    false,
			Fieldname: extraFieldname,
			Filename:  filename,
			Value:     descBytes,
		},
	}

	response, code, err := PostMultipartForm(fields, uri)
	return this.doProcessResponse(response,code,err,res)
}

func (this ApiTokenBase) doProcessResponse(response []byte, code int,  err error, res interface{}) error{
	if err != nil {
		return err
	}

	tmp := string(response)
	log.Print(tmp)

	if code == http.StatusNoContent && len(response) == 0{
		response = []byte("{}")
	}

	err = json.Unmarshal(response, res)
	if err != nil {
		return errors.New(string(response))
	}
	return nil
}
//---------------------------------------------------------------------------------------------------------------------

