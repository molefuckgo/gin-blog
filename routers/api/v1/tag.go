package v1

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/molefuckgo/gin-blog/models"
	"github.com/molefuckgo/gin-blog/pkg/e"
	"github.com/molefuckgo/gin-blog/pkg/logging"
	"github.com/molefuckgo/gin-blog/pkg/setting"
	"github.com/molefuckgo/gin-blog/pkg/util"
	"net/http"
)

// @Summary 获取多个文章标签
// @Produce json
// @Param page query int false "Page"
// @Success 200 {string} string "{"code":200,"data":{"lists":[{"id":1,"created_on":0,"modified_on":1573561866,"name":"1113111113","created_by":"","modified_by":"guowei","state":1}], "total":1}, "msg":"ok"}
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	//c *gin.Context是Gin很重要的组成部分，可以理解为上下文，它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		//fmt.Println("state:", state)
		maps["state"] = state
	}

	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @Summary 新增文章标签
// @Produce json
// Param name in type required description
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param create_by query string true "CreateBy"
// @Success 200 {string} string "{"code":200, "data":{}, "msg":"ok"}"
// @Failure 10001 {string} string "{"code":10001, "data":{}, "msg":"已存在该标签名称"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")                                   // Query获取请求url中参数值?name=123, name就是123了
	state, _ := com.StrTo(c.DefaultQuery("state", "0")).Int() // str转int,失败返回默认值0
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名字最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 修改文章标签
// @Produce json
// @Param name query string true "Name"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} string "{"code": 200, "data":{}, "msg":"ok"}"
// @Success 10002 {string} string "{"code":10002, "data":{}, "msg":"该标签不存在"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()
	//http://127.0.0.1:8080/api/v1/tags/7?name=edit1&state=0&modified_by=edit1
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state, _ := com.StrTo(arg).Int()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(name, "name").Message("名称不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(name, 100, "name").Message("名字最长为100字符")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人名字最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			code = e.SUCCESS
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

// @Summary 删除文章标签
// @Pruduce json
// @Success 200 {string} string "{"code":200, "data":{}, "msg":"ok"}"
// @Success 10002 {string} string "{"code":10002,"data":{},"msg":"该标签不存在"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			code = e.SUCCESS
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}
