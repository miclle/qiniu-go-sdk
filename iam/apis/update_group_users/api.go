// THIS FILE IS GENERATED BY api-generator, DO NOT EDIT DIRECTLY!

// 为用户分组中重新分配 IAM 子账号
package update_group_users

import (
	"encoding/json"
	credentials "github.com/qiniu/go-sdk/v7/storagev2/credentials"
)

// 调用 API 所用的请求
type Request struct {
	Alias       string                          // 用户分组别名
	Credentials credentials.CredentialsProvider // 鉴权参数，用于生成鉴权凭证，如果为空，则使用 HTTPClientOptions 中的 CredentialsProvider
	UserAliases UserAliases                     // IAM 子账号别名集合
}

// 为用户分组重新分配 IAM 子账号别名集合
type UserAliases = []string

// 为用户分组重新分配 IAM 子账号参数
type UpdatedGroupIamUsersParam = Request
type jsonRequest struct {
	UserAliases UserAliases `json:"user_aliases,omitempty"` // IAM 子账号别名集合
}

func (j *Request) MarshalJSON() ([]byte, error) {
	if err := j.validate(); err != nil {
		return nil, err
	}
	return json.Marshal(&jsonRequest{UserAliases: j.UserAliases})
}
func (j *Request) UnmarshalJSON(data []byte) error {
	var nj jsonRequest
	if err := json.Unmarshal(data, &nj); err != nil {
		return err
	}
	j.UserAliases = nj.UserAliases
	return nil
}
func (j *Request) validate() error {
	return nil
}

// 获取 API 所用的响应
type Response struct{}
