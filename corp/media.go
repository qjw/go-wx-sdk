package corp

import (
	"fmt"
	"io"
	"github.com/qjw/go-wx-sdk/utils"
)

const (
	mediaUpload = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	mediaGet    = "https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)

type MediaUploadRes struct {
	utils.CommonError
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	// 企业微信返回字符串，企业号返回数字，坑嗲
	// CreatedAt int64  `json:"created_at"`
}

func (this CorpApi) UploadTmpMedia(reader io.Reader, filename, tp string) (*MediaUploadRes, error) {
	var res MediaUploadRes
	if err := this.DoPostFile(reader, "media", filename, &res, mediaUpload, tp); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) GetTmpMedia(media_id string) (string, error) {
	accessToken, err := this.Context.GetAccessToken()
	if err != nil {
		return "", err
	}
	return string(fmt.Sprintf(mediaGet, accessToken, media_id)), nil
}
