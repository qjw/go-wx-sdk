package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"net/url"
	"strconv"
)

const (
	departmentList   = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%d"
	departmentList2  = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s"
	createDepartment = "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=%s"
	updateDepartment = "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=%s"
	deleteDepartment = "https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=%s&id=%d"
	userSimplelist   = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=%s&%s"
	userList         = "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&%s"
)

type DepartmentObj struct {
	ID       *int64 `json:"id,omitempty"`
	Name     string `json:"name"`
	ParentID int64  `json:"parentid"`
	Order    *int64 `json:"order,omitempty"`
}

type DepartListObj struct {
	utils.CommonError
	Departments []DepartmentObj `json:"department"`
}

func (this CorpApi) GetDepartments(parent *int64) (*DepartListObj, error) {
	var res DepartListObj
	var err error
	if parent != nil {
		err = this.DoGet(departmentList, &res, *parent)
	} else {
		err = this.DoGet(departmentList2, &res)
	}
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CreateDepartmentRes struct {
	utils.CommonError
	id int64 `json:"id"`
}

func (this CorpApi) CreateDepartment(department *DepartmentObj) (*CreateDepartmentRes, error) {
	var res CreateDepartmentRes
	if err := this.DoPostObject(createDepartment, department, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type DepartmentUpdateObj struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name,omitempty"`
	ParentID *int64  `json:"parentid,omitempty"`
	Order    *int64  `json:"order,omitempty"`
}

func (this CorpApi) UpdateDepartment(department *DepartmentUpdateObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(updateDepartment, department, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) DeleteDepartments(id int64) (*utils.CommonError, error) {
	var res utils.CommonError
	err := this.DoGet(deleteDepartment, &res, id)
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type SimpleUserlistObj struct {
	ID         int64  `json:"department_id"`
	FetchChild *int64 `json:"fetch_child,omitempty" doc:"1/0：是否递归获取子部门下面的成员"`
	Status     *int64 `json:"status,omitempty" doc:"0获取全部成员，1获取已关注成员列表，2获取禁用成员列表，4获取未关注成员列表。status可叠加,未填写则默认为4"`
}

type SimpleUserListRes struct {
	utils.CommonError
	UserList []struct {
		UserID     string  `json:"userid"`
		Name       string  `json:"name"`
		Department []int64 `json:"department"`
	} `json:"userlist"`
}

func (this CorpApi) DepartmentSimpleUserlist(param *SimpleUserlistObj) (*SimpleUserListRes, error) {
	v := url.Values{}
	v.Add("department_id", strconv.FormatInt(param.ID, 10))
	if param.FetchChild != nil {
		v.Add("fetch_child", strconv.FormatInt(*param.FetchChild, 10))
	}
	if param.Status != nil {
		v.Add("status", strconv.FormatInt(*param.Status, 10))
	}

	var res SimpleUserListRes
	err := this.DoGet(userSimplelist, &res, v.Encode())
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UserCreateObj struct {
	UserID     string  `json:"userid" doc:"成员UserID。对应管理端的帐号，企业内必须唯一。长度为1~64个字节"`
	Name       string  `json:"name" doc:"成员名称。长度为0~64个字节"`
	Department []int64 `json:"department" doc:"成员所属部门id列表，不超过20个"`
	Position   *string `json:"position,omitempty" doc:"职位信息。长度为0~64个字节"`
	Mobile     *string `json:"mobile,omitempty" doc:"手机号码。企业内必须唯一，mobile/weixinid/email三者不能同时为空"`
	Gender     *string `json:"gender,omitempty" doc:"性别。1表示男性，2表示女性"`
	Email      *string `json:"email,omitempty" doc:"邮箱。长度为0~64个字节。企业内必须唯一"`
	Weixinid   *string `json:"weixinid,omitempty" doc:"微信号。企业内必须唯一。（注意：是微信号，不是微信的名字）"`
	Avatar     *string `json:"avatar,omitempty" doc:"成员头像的mediaid，通过多媒体接口上传图片获得的mediaid"`
	Extattr    *struct {
		Attrs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attrs,omitempty"`
	} `json:"extattr,omitempty" doc:"扩展属性。扩展属性需要在WEB管理端创建后才生效，否则忽略未知属性的赋值userid"`
}

type UserUpdateObj struct {
	UserID     string   `json:"userid" doc:"成员UserID。对应管理端的帐号，企业内必须唯一。长度为1~64个字节"`
	Name       *string  `json:"name,omitempty" doc:"成员名称。长度为0~64个字节"`
	Department *[]int64 `json:"department,omitempty" doc:"成员所属部门id列表，不超过20个"`
	Position   *string  `json:"position,omitempty" doc:"职位信息。长度为0~64个字节"`
	Mobile     *string  `json:"mobile,omitempty" doc:"手机号码。企业内必须唯一，mobile/weixinid/email三者不能同时为空"`
	Gender     *string  `json:"gender,omitempty" doc:"性别。1表示男性，2表示女性"`
	Email      *string  `json:"email,omitempty" doc:"邮箱。长度为0~64个字节。企业内必须唯一"`
	Weixinid   *string  `json:"weixinid,omitempty" doc:"微信号。企业内必须唯一。（注意：是微信号，不是微信的名字）"`
	Avatar     *string  `json:"avatar,omitempty" doc:"成员头像的mediaid，通过多媒体接口上传图片获得的mediaid"`
	Extattr    *struct {
		Attrs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attrs,omitempty"`
	} `json:"extattr,omitempty" doc:"扩展属性。扩展属性需要在WEB管理端创建后才生效，否则忽略未知属性的赋值userid"`
	Enable int `json:"enable" doc:"启用/禁用成员。1表示启用成员，0表示禁用成员"`
}

type UserObj struct {
	UserCreateObj
	Status int `json:"status"`
}

type UserListRes struct {
	utils.CommonError
	UserList []UserObj `json:"userlist"`
}

func (this CorpApi) DepartmentUserlist(param *SimpleUserlistObj) (*UserListRes, error) {
	v := url.Values{}
	v.Add("department_id", strconv.FormatInt(param.ID, 10))
	if param.FetchChild != nil {
		v.Add("fetch_child", strconv.FormatInt(*param.FetchChild, 10))
	}
	if param.Status != nil {
		v.Add("status", strconv.FormatInt(*param.Status, 10))
	}

	var res UserListRes
	err := this.DoGet(userList, &res, v.Encode())
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

const (
	userGet         = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
	userDelete      = "https://qyapi.weixin.qq.com/cgi-bin/user/delete?access_token=%s&userid=%s"
	userBatchDelete = "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete?access_token=%s"
	userCreate      = "https://qyapi.weixin.qq.com/cgi-bin/user/create?access_token=%s"
	userUpdate      = "https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token=%s"
)

type UserInfoRes struct {
	utils.CommonError
	UserObj
}

func (this CorpApi) GetUser(userid string) (*UserInfoRes, error) {
	var res UserInfoRes
	if err := this.DoGet(userGet, &res, userid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) DeleteUser(userid string) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoGet(userDelete, &res, userid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type BatchDeleteUserObj struct {
	UserIDList []string `json:"useridlist"`
}

func (this CorpApi) BatchDeleteUser(userids []string) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(userBatchDelete, &BatchDeleteUserObj{
		UserIDList: userids,
	}, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) CreateUser(param *UserCreateObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(userCreate, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) UpdateUser(param *UserUpdateObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(userUpdate, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

const (
	tagCreate   = "https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token=%s"
	tagUpdate   = "https://qyapi.weixin.qq.com/cgi-bin/tag/update?access_token=%s"
	tagDelete   = "https://qyapi.weixin.qq.com/cgi-bin/tag/delete?access_token=%s&tagid=%d"
	tagGetUsers = "https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=%s&tagid=%d"
	tagAddUsers = "https://qyapi.weixin.qq.com/cgi-bin/tag/addtagusers?access_token=%s"
	tagDelUsers = "https://qyapi.weixin.qq.com/cgi-bin/tag/deltagusers?access_token=%s"
	tagList     = "https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=%s"
)

type TagUserListRes struct {
	utils.CommonError
	Userlist []struct {
		Userid string `json:"userid"`
		Name   string `json:"name"`
	} `json:"userlist"`
	Partylist []int64 `json:"partylist"`
}

func (this CorpApi) GetTagUsers(tagid int64) (*TagUserListRes, error) {
	var res TagUserListRes
	if err := this.DoGet(tagGetUsers, &res, tagid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TagUserUpdateObj struct {
	TagID     int64    `json:"tagid"`
	Userlist  []string `json:"userlist"`
	Partylist []int64  `json:"partylist"`
}

type TagUpdateRes struct {
	utils.CommonError
	Invalidlist  string  `json:"invalidlist"`
	Invalidparty []int64 `json:"invalidparty"`
}

func (this CorpApi) AddTagUsers(tag *TagUserUpdateObj) (*TagUpdateRes, error) {
	var res TagUpdateRes
	if err := this.DoPostObject(tagAddUsers, tag, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) DelTagUsers(tag *TagUserUpdateObj) (*TagUpdateRes, error) {
	var res TagUpdateRes
	if err := this.DoPostObject(tagDelUsers, tag, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TagObj struct {
	Tagid   int64  `json:"tagid"`
	Tagname string `json:"tagname"`
}

type TagCreateObj struct {
	Tagid   *int64 `json:"tagid,omitempty"`
	Tagname string `json:"tagname"`
}

type TagListRes struct {
	utils.CommonError
	Taglist []TagObj `json:"taglist"`
}

func (this CorpApi) UpdateTag(tag *TagObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(tagUpdate, tag, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TagCreateRes struct {
	utils.CommonError
	Tagid int64 `json:"tagid"`
}

func (this CorpApi) CreateTag(tag *TagCreateObj) (*TagCreateRes, error) {
	var res TagCreateRes
	if err := this.DoPostObject(tagCreate, tag, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) DeleteTag(tagid int64) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoGet(tagDelete, &res, tagid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) GetTags() (*TagListRes, error) {
	var res TagListRes
	if err := this.DoGet(tagList, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}
