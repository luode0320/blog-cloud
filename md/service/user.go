package service

import (
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
)

// 更新用户密码
func UserUpdatePassword(userCondition entity.UserCondition) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	// 查询用户
	user, err := dao.UserGetById(tx, userCondition.Id)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	// 判断原密码相同
	if util.EncryptSHA256([]byte(user.Id+userCondition.Password)) != user.Password {
		panic(common.NewError("原密码不正确"))
	}

	// 更新用户
	user.Password = util.EncryptSHA256([]byte(user.Id + userCondition.NewPassword))
	err = dao.UserResetPassword(tx, user)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}
}
