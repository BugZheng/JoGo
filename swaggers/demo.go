/**
 * @Author: BugZheng
 * @Description:
 * @File:  user
 * @Version: 1.0.0
 * @Date: 2021/02/01 7:32 下午
 */
package swaggers

type UserListData struct {
	List []struct {
		ID         int    `json:"id"`         //主键
		UserName   string `json:"userName"`   //顾问账号
		CreateDate string `json:"createDate"` //创建时间
		LastEditBy string `json:"lastEditBy"` //最后编辑人
		Tel        string `json:"tel"`        //手机号
		Createby   string `json:"createby"`   //创建人
	} `json:"list"`
	Total int `json:"total"` //总记录数
}

type DemoData struct {
	Avatar         string      `json:"Avatar"`
	CreatedAt      string      `json:"CreatedAt"`
	DeletedAt      interface{} `json:"DeletedAt"`
	ID             int64       `json:"ID"`
	Nickname       string      `json:"Nickname"`
	PasswordDigest string      `json:"PasswordDigest"`
	Status         string      `json:"Status"`
	UpdatedAt      string      `json:"UpdatedAt"`
	UserName       string      `json:"UserName"`
}
