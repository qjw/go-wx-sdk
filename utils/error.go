package utils

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (this CommonError) IsOK() bool {
	return this.ErrCode == 0
}

type Pagination struct {
	Offset int `json:"offset"`
	Count int `json:"count"`
}