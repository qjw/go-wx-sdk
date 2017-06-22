package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, int, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	defer response.Body.Close()
	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, response.StatusCode,
			fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	body,err := ioutil.ReadAll(response.Body)
	return body,response.StatusCode,err
}

func HTTPDelete(uri string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, response.StatusCode, err
	}

	defer response.Body.Close()
	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, response.StatusCode,
			fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	return body,response.StatusCode,err
}


//PostJSON post json 数据请求
func PostJSON(uri string, obj interface{}) ([]byte, int, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, response.StatusCode, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, http.StatusBadRequest, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	body2, err := ioutil.ReadAll(response.Body)
	return body2,response.StatusCode,err
}

//PostJSON post json 数据请求
func PostJSONRaw(uri string, obj []byte) ([]byte, int, error) {
	body := bytes.NewBuffer(obj)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, response.StatusCode, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	body2, err := ioutil.ReadAll(response.Body)
	return body2,response.StatusCode,err
}

//PostFile 上传文件
func PostFile(fieldname, filename, uri string) ([]byte, int, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(fields, uri)
}

//PostFile 上传文件
func PostReader(fieldname, filename, uri string, reader io.Reader) ([]byte, int, error) {
	fields := []MultipartFormField{
		{
			IsFile:    false,
			Fieldname: fieldname,
			Filename:  filename,
			Reader:    reader,
		},
	}
	return PostMultipartForm(fields, uri)
}

//MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
	Reader    io.Reader
}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, code int, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	code = http.StatusBadRequest

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			var partWriter io.Writer = nil
			var e error = nil
			//
			if field.Reader == nil {
				field.Reader = bytes.NewReader(field.Value)
				partWriter, e = bodyWriter.CreateFormField(field.Fieldname)
			} else {
				partWriter, e = bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			}
			if e != nil {
				err = e
				return
			}

			if _, err = io.Copy(partWriter, field.Reader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	response, e := http.Post(uri, contentType, bodyBuf)
	if e != nil {
		err = e
		return
	}

	code = response.StatusCode
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return
	}
	respBody, err = ioutil.ReadAll(response.Body)
	return
}