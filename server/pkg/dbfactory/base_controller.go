package dbfactory

import (
	"fmt"
	"godocms/libs"
	"math"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BaseController 是一个通用的控制器基类
type BaseController struct{}

// GetPagination 获取分页参数
func (b *BaseController) GetPagination(ctx *gin.Context) (int, int) {
	pageSize := validateParameter(ctx.DefaultQuery("limit", "10"), 10, 1, 100)
	page := validateParameter(ctx.Query("page"), 1, 1, math.MaxInt32)
	return page, pageSize
}

// validateParameter 解析并校验参数，确保其在指定范围内
func validateParameter(param string, defaultValue, min, max int) int {
	value, err := strconv.Atoi(param)
	if err != nil || value < min || value > max {
		return defaultValue
	}
	return value
}

// List 获取列表
func (b *BaseController) List(ctx *gin.Context, table string, conditions map[string]interface{}, orderby string, result *PageResult) {
	page, pageSize := b.GetPagination(ctx)

	var list *PageResult
	pageParam := PageParams{Page: page, PageSize: pageSize, OrderBy: orderby}
	err := Db.GetPage(table, conditions, pageParam, list)
	if err != nil {
		libs.Error(ctx, "获取"+table+"列表失败: "+err.Error())
		return
	}

	libs.Success(ctx, "获取"+table+"列表成功", list)
}

// AddBefore 添加前的页面展示
func (b *BaseController) AddBefore(ctx *gin.Context) {
	libs.Success(ctx, "success", nil)
}

// AddSave 添加记录
func (b *BaseController) AddSave(ctx *gin.Context, model interface{}, tableName string) {
	if err := ctx.ShouldBindJSON(model); err != nil {
		libs.Error(ctx, "参数错误: "+err.Error())
		return
	}
	res := Db.Create(tableName, model)
	if res != nil {
		libs.Error(ctx, "添加"+tableName+"失败: "+res.Error())
		return
	}
	libs.Success(ctx, "添加成功", nil)
}

// EditBefore 编辑前的页面展示
func (b *BaseController) EditBefore(ctx *gin.Context, model interface{}, tableName string) {
	idStr := ctx.Query("id")
	if idStr == "" {
		libs.Error(ctx, "id为空")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		libs.Error(ctx, "参数不是数字")
		return
	}
	if err := Db.GetByID(tableName, id, model).Error; err != nil {
		libs.Error(ctx, "获取"+tableName+"失败:"+err())
		return
	}
	libs.Success(ctx, "获取"+tableName+"成功", model)
}

// EditSave 编辑记录
func (b *BaseController) EditSave(ctx *gin.Context, model interface{}, tableName string) {
	if err := ctx.ShouldBindJSON(model); err != nil {
		libs.Error(ctx, "参数错误: "+err.Error())
		return
	}
	if err := Db.Update(tableName, model).Error; err != nil {
		libs.Error(ctx, "编辑"+tableName+"失败: "+err())
		return
	}
	libs.Success(ctx, "编辑成功", nil)
}

// Delete 删除记录
func (b *BaseController) Delete(ctx *gin.Context, model interface{}, tableName string) {
	idStr := ctx.Query("id")
	if idStr == "" {
		libs.Error(ctx, "id为空")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		libs.Error(ctx, "参数不是数字")
		return
	}
	if err := Db.Delete(tableName, id).Error; err != nil {
		libs.Error(ctx, "删除"+tableName+"失败: "+err())
		return
	}
	libs.Success(ctx, "删除成功", nil)
}
func (b *BaseController) GetQueryInt(ctx *gin.Context, param string) (int, error) {
	idStr := ctx.Query(param)
	if idStr == "" {
		return 0, fmt.Errorf("id为空")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("参数不是数字")
	}
	return id, nil
}

// UpdateEnable 更新 enable 字段
func (b *BaseController) UpdateEnable(ctx *gin.Context, model interface{}, tableName string) {
	id, err := b.GetQueryInt(ctx, "id")
	if err != nil {
		libs.Error(ctx, err.Error())
		return
	}

	db := Db
	if err := db.GetByID(tableName, id, model).Error; err != nil {
		libs.Error(ctx, "获取"+tableName+"失败: "+err())
		return
	}

	// 使用反射获取 enable 字段的值
	v := reflect.ValueOf(model).Elem()
	enableField := v.FieldByName("Enable")
	if !enableField.IsValid() {
		libs.Error(ctx, "模型中没有 enable 字段")
		return
	}

	// 检查 enable 字段是否可设置
	if !enableField.CanSet() {
		libs.Error(ctx, "enable 字段不可设置")
		return
	}

	// 设置 enable 字段的值
	enable := !enableField.Bool()
	enableField.SetBool(enable)

	if err := db.Update(tableName, model).Error; err != nil {
		libs.Error(ctx, "更新"+tableName+"的enable字段失败: "+err())
		return
	}

	libs.Success(ctx, "更新成功", nil)
}
