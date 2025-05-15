package service

import (
	"godocms/model"
	"godocms/pkg/db"
	"godocms/utils"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserAuthAndShare(user *model.User) ([]model.User, []model.User, error) {
	var shareList []model.User
	if err := db.DB.Where("dept_id = ?", user.DeptId).Select("id", "nickname", "username", "email").Find(&shareList).Error; err != nil {
		return nil, nil, err
	}
	var userAuths []model.User
	if user.DeptAuth != "0" {
		if user.DeptAuth == "-1" {
			userAuths = shareList
		} else {
			authIds := strings.Split(user.DeptAuth, ",")
			if err := db.DB.Where("id IN(?)", authIds).Select("id", "nickname", "username", "email").Find(&userAuths).Error; err != nil {
				return nil, nil, err
			}
		}
	}

	return userAuths, shareList, nil
}
func UpdateLoginUser(c *gin.Context, user *model.User) error {
	ip := utils.GetIpAddress(c)
	user.LoginIP = ip
	user.UpTime = int32(time.Now().Unix())
	user.LoginNum++
	// user.UseSpace = userSum
	if err := db.DB.Save(&user).Error; err != nil {
		slog.Error("更新用户信息失败", "err", err.Error())
		return err
	}
	return nil
}
func GetUserDept(deptId int32) (model.UserDept, error) {
	var userDept model.UserDept
	if err := db.DB.First(&userDept, deptId).Error; err != nil {
		return userDept, err
	}
	return userDept, nil
}
func GetUserRole(roleId int32) (model.UserRole, error) {
	var role model.UserRole
	if err := db.DB.First(&role, roleId).Error; err != nil {
		return role, err
	}
	return role, nil
}
