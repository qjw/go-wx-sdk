package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"encoding/json"
	"fmt"
	"io"
)

//-----------------------------------素材--------------------------------------------------------------------------------

const (
	get_material_count = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=%s"
	batch_get_material = "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s"
	media_upload       = "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	media_get          = "https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
	mediaVideoGet      = "http://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
	add_material       = "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=%s"
	get_material       = "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=%s"
	del_material       = "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token=%s"
	media_uploadimg    = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s"
	add_news           = "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=%s"
)

const (
	get_material_temp = `{
		"media_id":"%s"
	}`
	batchgetMaterialTemp = `{"type":"%s","offset":%d,"count":%d}`
)


type MaterialCount struct {
	VoiceCount int `json:"voice_count"`
	VideoCount int `json:"video_count"`
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
}

func (this WechatApi) GetMaterialCount() (*MaterialCount, error) {
	var res MaterialCount
	if err := this.DoGet(get_material_count, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type Material struct {
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Items      []struct {
		MediaID    string `json:"media_id"`
		Name       string `json:"name"`
		UpdateTime int64  `json:"update_time"`
		Url        string `json:"url"`
	} `json:"item"`
}

func (this WechatApi) GetMaterials(ty string, offset int, count int) (*Material, error) {
	var res Material
	body := fmt.Sprintf(batchgetMaterialTemp, ty, offset, count)
	if err := this.DoPost(batch_get_material, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type NewsMaterial struct {
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Items      []struct {
		MediaID    string `json:"media_id"`
		UpdateTime int64  `json:"update_time"`
		Content    struct {
			NewsItem []struct {
				Title              string `json:"title"`
				ThumbMediaID       string `json:"thumb_media_id"`
				ShowCoverPic       int    `json:"show_cover_pic"`
				Author             string `json:"author"`
				Digest             string `json:"digest"`
				Content            string `json:"content"`
				Url                string `json:"url"`
				ContentSourceUrl   string `json:"content_source_url"`
				NeedOpenComment    int    `json:"need_open_comment"`
				OnlyFansCanComment int    `json:"only_fans_can_comment"`
				ThumbUrl           string `json:"thumb_url"`
			} `json:"news_item"`
		} `json:"content"`
	} `json:"item"`
}

func (this WechatApi) GetNewsMaterials(offset int, count int) (*NewsMaterial, error) {
	var res NewsMaterial
	body := fmt.Sprintf(batchgetMaterialTemp, "news", offset, count)
	if err := this.DoPost(batch_get_material, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UploadRes struct {
	utils.CommonError
	Type         string `json:"type"`
	MediaID      string `json:"media_id"`
	ThumbMediaID string `json:"thumb_media_id"`
	CreatedAt    int64  `json:"created_at"`
}

func (this WechatApi) UploadTmpMaterial(reader io.Reader, filename, tp string) (*UploadRes, error) {
	var res UploadRes
	if err := this.DoPostFile(reader, "media", filename, &res, media_upload, tp); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetTmpMaterial(media_id string) (string, error) {
	accessToken, err := this.Context.GetAccessToken()
	if err != nil {
		return "", err
	}
	return string(fmt.Sprintf(media_get, accessToken, media_id)), nil
}

type VideoTmpMaterialRes struct {
	VideoUrl string `json:"video_url"`
}

func (this WechatApi) GetVideoTmpMaterial(media_id string) (*VideoTmpMaterialRes, error) {
	var res VideoTmpMaterialRes
	if err := this.DoGet(mediaVideoGet, &res, media_id); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) UploadMaterial(reader io.Reader, filename, tp string) (*UploadRes, error) {
	var res UploadRes
	if err := this.DoPostFile(reader, "media", filename, &res, add_material, tp); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type VideoMaterialInfo struct {
	Description string `json:"introduction"`
	Title       string `json:"title"`
}

func (this WechatApi) UploadVideoMaterial(reader io.Reader,
	filename string, info *VideoMaterialInfo) (*UploadRes, error) {
	var res UploadRes
	if err := this.DoPostFileExtra(reader, "media", "description", filename, info,
		&res, add_material, "video"); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetMaterial() (string, error) {
	accessToken, err := this.Context.GetAccessToken()
	if err != nil {
		return "", err
	}
	return string(fmt.Sprintf(get_material, accessToken)), nil
}

type VideoMaterialRes struct {
	DownUrl     string `json:"down_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (this WechatApi) GetVideoMaterial(media_id string) (*VideoMaterialRes, error) {
	var res VideoMaterialRes
	body := fmt.Sprintf(get_material_temp, media_id)
	if err := this.DoPost(get_material, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) DeleteMaterial(media_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(get_material_temp, media_id)
	if err := this.DoPost(del_material, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type ArticleImageRes struct {
	Url string `json:"url"`
}

func (this WechatApi) UploadArticleImage(reader io.Reader, filename string) (*ArticleImageRes, error) {
	var res ArticleImageRes
	if err := this.DoPostFile(reader, "media", filename, &res, media_uploadimg); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type ArticleCreateEntry struct {
	Title            string `json:"title" doc:"图文消息的标题"`
	ThumbMediaID     string `json:"thumb_media_id" doc:"图文消息的封面图片素材id（必须是永久mediaID）"`
	Author           string `json:"author" doc:"作者"`
	Digest           string `json:"digest" doc:"图文消息的摘要，仅有单图文消息才有摘要，多图文此处为空"`
	ShowCoverPic     int8   `json:"show_cover_pic" doc:"是否显示封面，0为false，即不显示，1为true，即显示"`
	Content          string `json:"content"`
	ContentSourceUrl string `json:"content_source_url" doc:"图文消息的原文地址，即点击“阅读原文”后的URL"`
	Url              string `json:"url"`
}

type ArticleCreate struct {
	Articles []ArticleCreateEntry `json:"articles"`
}

type ArticleCreateRes struct {
	MediaID string `json:"media_id"`
}

func (this WechatApi) CreateArticle(articles []ArticleCreateEntry) (*ArticleCreateRes, error) {
	var res ArticleCreateRes
	bodyObj := ArticleCreate{
		Articles: articles,
	}
	body, _ := json.Marshal(&bodyObj)
	if err := this.DoPostRaw(add_news, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}


type ArticleDetail struct {
	NewsItem []ArticleCreateEntry `json:"news_item"`
}

func (this WechatApi) GetArticleMaterial(media_id string) (*ArticleDetail, error) {
	var res ArticleDetail
	body := fmt.Sprintf(get_material_temp, media_id)
	if err := this.DoPost(get_material, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}