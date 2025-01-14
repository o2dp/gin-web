package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	// 定期缓存, 避免每次频繁查询数据库
	menuTreeCache = cache.New(24*time.Hour, 48*time.Hour)
)

// 查询当前用户菜单树
func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)
	oldCache, ok := menuTreeCache.Get(fmt.Sprintf("%d", user.Id))
	if ok {
		resp, _ := oldCache.([]response.MenuTreeResponseStruct)
		response.SuccessWithData(resp)
		return
	}

	// 创建服务
	s := service.New(c)
	menus, err := s.GetMenuTree(user.RoleId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	// 写入缓存
	menuTreeCache.Set(fmt.Sprintf("%d", user.Id), resp, cache.DefaultExpiration)
	response.SuccessWithData(resp)
}

// 查询指定角色的菜单树
func GetAllMenuByRoleId(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	// 创建服务
	s := cache_service.New(c)
	menus, ids, err := s.GetAllMenuByRoleId(user.Role, utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	var resp response.MenuTreeWithAccessResponseStruct
	resp.AccessIds = ids
	utils.Struct2StructByJson(menus, &resp.List)
	response.SuccessWithData(resp)
}

// 查询所有菜单
func GetMenus(c *gin.Context) {
	// 绑定当前用户角色排序(隐藏特定用户)
	user := GetCurrentUser(c)
	// 创建服务
	s := cache_service.New(c)
	menus := s.GetMenus(user.Role)
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(resp)
}

// 创建菜单
func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	// 绑定参数
	var req request.CreateMenuRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 参数校验
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.Nickname + user.Username
	// 创建服务
	s := service.New(c)
	err = s.CreateMenu(user.Role, &req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 更新菜单
func UpdateMenuById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateMenuRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 获取path中的menuId
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		response.FailWithMsg("菜单编号不正确")
		return
	}
	// 创建服务
	s := service.New(c)
	// 更新数据
	err = s.UpdateById(menuId, req, new(models.SysMenu))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批量删除菜单
func BatchDeleteMenuByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// 创建服务
	s := service.New(c)
	// 删除数据
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysMenu))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
