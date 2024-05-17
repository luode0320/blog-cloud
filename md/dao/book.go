package dao

import (
	"md/model/entity"
	"md/util"
	"sort"

	"github.com/jmoiron/sqlx"
)

// 添加一级目录
func BookAdd(tx *sqlx.Tx, book entity.Book) error {
	sql := `insert into t_book (id,name,create_time,user_id) values (:id,:name,:create_time,:user_id)`
	_, err := tx.NamedExec(sql, book)
	return err
}

// 修改一级目录
func BookUpdate(tx *sqlx.Tx, book entity.Book) error {
	sql := `update t_book set name=:name where id=:id and user_id=:user_id`
	_, err := tx.NamedExec(sql, book)
	return err
}

// 根据id删除一级目录
func BookDeleteById(tx *sqlx.Tx, id, userId string) error {
	sql := `delete from t_book where id=$1 and user_id=$2`
	_, err := tx.Exec(sql, id, userId)
	return err
}

// 查询一级目录列表
func BookList(db *sqlx.DB, userId string) ([]entity.Book, error) {
	sql := `select id,name,create_time from t_book where user_id=$1`
	result := []entity.Book{}
	err := db.Select(&result, sql, userId)
	// 按名称升序
	sort.Slice(result, func(i, j int) bool {
		return util.StringSort(result[i].Name, result[j].Name)
	})
	return result, err
}

// 根据名称查询一级目录列表
func BookListByName(tx *sqlx.Tx, name, userId string) ([]entity.Book, error) {
	sql := `select * from t_book where user_id=$1 and name=$2`
	result := []entity.Book{}
	err := tx.Select(&result, sql, userId, name)
	return result, err
}

// 查询一级目录
func Book(db *sqlx.DB, id string) (entity.Book, error) {
	sql := `select * from t_book where id=$1`
	result := entity.Book{}
	err := db.Get(&result, sql, id)
	return result, err
}
