// @Time : 2020/8/18 11:26 AM
// @Author : acol
// @File : baseResponse
// @Software: GoLand
// @Desc: to do somewhat..

package swaggers

//SwagCommonResponse ...接口数据返回结构
type SwagCommonResponse struct {
	Code string      `json:"code" example:"200"` //200:执行成功，其他:执行出错
	Msg  string      `json:"msg" example:"成功"`   //执行接口的提示信息
	Data interface{} `json:"data"`               //执行接口后返回数据
}

//SwagListData ...接口数据列表返回结构
type SwagListData struct {
	List  interface{} `json:"list"`  //返回列表
	Total int64       `json:"total"` //总数目
}

//SwagRowData ...接口数据列表返回结构
type SwagRowData struct {
	Rows      interface{} `json:"rows"`      //返回列表
	TotalRows int64       `json:"totalRows"` //总数目
}
