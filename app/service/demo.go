/**
 * @Author: Bugzheng
 * @Description:
 * @File:  user
 * @Version: 1.0.0
 * @Date: 2021/02/01 4:01 下午
 */
package service

import (
	"JoGo/app/model"
	"JoGo/pkg/serializer"
	"net/http"
)

//ttps://www.liwenzhou.com/posts/Go/validator_usages/

// UserRegisterService 管理用户注册服务
type UserRegisterService struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

//UserLoginService  管理用户登录的服务
type UserLoginService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"` //用户姓名
	Password string `form:"password" json:"password" binding:"required,min=8,max=40"`   //密码
	Test     string `form:"test" json:"test" binding:"required,min=8,max=40"`           //测试字段
	Test2    string `form:"test2" json:"test2" binding:"required,min=5,max=30"`         //测试字段2
}

//DemoService demo参数验证结构体
type DemoService struct {
	Test1 string `form:"test1" json:"test1" binding:"required,min=5,max=30"` //测试字段1
	Test2 string `form:"test2" json:"test2" binding:"required"`              //测试字段2
	Test3 string `form:"test3" json:"test3" binding:"-"`                     //测试字段3
}

// DemoAPI 测试
func (d *DemoService) DemoAPI() *serializer.Response {
	user, _ := model.GetUser()
	return &serializer.Response{
		Code: http.StatusOK,
		Data: user,
		Msg:  http.StatusText(200),
	}
}

//// valid 验证表单
//func (service *UserRegisterService) valid() *serializer.Response {
//	if service.PasswordConfirm != service.Password {
//		return &serializer.Response{
//			Code: 40001,
//			Msg:  "两次输入的密码不相同",
//		}
//	}
//
//	count := 0
//	boot.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
//	if count > 0 {
//		return &serializer.Response{
//			Code: 40001,
//			Msg:  "昵称被占用",
//		}
//	}
//
//	count = 0
//	boot.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
//	if count > 0 {
//		return &serializer.Response{
//			Code: 40001,
//			Msg:  "用户名已经注册",
//		}
//	}
//
//	return nil
//}

//// Register 用户注册
//func (service *UserRegisterService) Register() serializer.Response {
//	user := model.User{
//		Nickname: service.Nickname,
//		UserName: service.UserName,
//		Status:   model.Active,
//	}
//
//	// 表单验证
//	if err := service.valid(); err != nil {
//		return *err
//	}
//
//	// 加密密码
//	if err := user.SetPassword(service.Password); err != nil {
//		return serializer.Err(
//			serializer.CodeEncryptError,
//			"密码加密失败",
//			err,
//		)
//	}
//
//	// 创建用户
//	if err := boot.DB.Create(&user).Error; err != nil {
//		return serializer.ParamErr("注册失败", err)
//	}
//	return serializer.Response{
//		Data: user,
//	}
//}
