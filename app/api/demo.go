package api

import (
	"JoGo/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Demo ...
// @Tags user
// @Summary 接口名
// @Produce json
// @Param   token     header    string     true     "登录token"
// @Param   xxx     query    service.UserLoginService     false        "字段注释"
// @Success 200 {object} swaggers.SwagCommonResponse{data=[]swaggers.UserListData}
// @Router /api/v1/demo [get]
func Demo(c *gin.Context) {
	var service service.UserLoginService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, service)
}

// DemoAPI 这是一个测试接口
// @Tags demo
// @Summary 这是一个测试接口
// @Produce json
// @Param   token     header    string     true     "登录token"
// @Param   xxx     query    service.DemoService    false        "字段注释"
// @Success 200 {object} swaggers.SwagCommonResponse{data=[]swaggers.写你自己的返回结构体}
// @Router /demo/DemoApi [post]
func DemoAPI(c *gin.Context) {
	var service service.DemoService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse(err))
		return
	}
	res := service.DemoAPI()
	c.JSON(http.StatusOK, res)
}
