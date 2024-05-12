package dao

import (
	"errors"
	"md/model/common"
	"md/model/entity"

	"github.com/jmoiron/sqlx"
)

// 添加用户
func UserAdd(tx *sqlx.Tx, user entity.User) error {
	sql := `insert into t_user (id,name,password,create_time) values (:id,:name,:password,:create_time)`
	_, err := tx.NamedExec(sql, user)
	return err
}

// 修改用户密码
func UserResetPassword(tx *sqlx.Tx, user entity.User) error {
	sql := `update t_user set password=:password where id=:id`
	_, err := tx.NamedExec(sql, user)
	return err
}

// 根据id查询用户
func UserGetById(tx interface{}, id string) (entity.User, error) {
	sql := `select * from t_user where id=$1`
	result := entity.User{}
	var err error
	switch tx := tx.(type) {
	case *sqlx.Tx:
		err = tx.Get(&result, sql, id)
	case *sqlx.DB:
		err = tx.Get(&result, sql, id)
	default:
		err = errors.New("数据库事务异常")
	}
	return result, err
}

// 根据用户名查询用户
func UserGetByName(db *sqlx.DB, name string) (entity.User, error) {
	sql := `select * from t_user where name=$1`
	result := entity.User{}
	err := db.Get(&result, sql, name)
	return result, err
}

// 根据用户名查询用户数量
func UserCountByName(tx *sqlx.Tx, name string) (common.CountResult, error) {
	sql := `select count(*) as count from t_user where name=$1`
	result := common.CountResult{}
	err := tx.Get(&result, sql, name)
	return result, err
}

// 查询用户数量
func UserCount(tx *sqlx.Tx) (common.CountResult, error) {
	sql := `select count(*) as count from t_user`
	result := common.CountResult{}
	err := tx.Get(&result, sql)
	return result, err
}
