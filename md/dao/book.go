package dao

import (
	"github.com/jmoiron/sqlx"
	"md/model/entity"
	"sort"
	"strconv"
	"strings"
)

// 添加一级目录
func BookAdd(tx *sqlx.Tx, book entity.Book) error {
	sql := `insert into t_book (id,parent_id,name,create_time,user_id) values (:id,:parent_id,:name,:create_time,:user_id)`
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
	sql := `select * from t_book where user_id=$1 and parent_id=''`
	books := []entity.Book{}
	err := db.Select(&books, sql, userId)
	sortBooks(books)

	result := []entity.Book{}
	// 查询每个根book的所有子集， 并且按顺序排序在根book的后面
	for _, book := range books {
		result = append(result, book)
		subBooks, err := BookByParentId(db, userId, book.Id)
		if err != nil {
			return result, err
		}
		sortBooks(subBooks)
		result = append(result, subBooks...)
	}

	return result, err
}

// 根据名称查询一级目录列表
func BookListByName(tx *sqlx.Tx, name, userId string) ([]entity.Book, error) {
	sql := `select * from t_book where user_id=$1 and name=$2`
	result := []entity.Book{}
	err := tx.Select(&result, sql, userId, name)
	return result, err
}

// 根据id查询二级目录
func BookByParentId(db *sqlx.DB, userId string, parentId string) ([]entity.Book, error) {
	sql := `select * from t_book where user_id=$1 and parent_id=$2`
	result := []entity.Book{}
	err := db.Select(&result, sql, userId, parentId)
	return result, err
}

// 查询一级目录
func Book(db *sqlx.DB, id string) (entity.Book, error) {
	sql := `select * from t_book where id=$1`
	result := entity.Book{}
	err := db.Get(&result, sql, id)
	return result, err
}

// 自定义目录排序逻辑
func sortBooks(books []entity.Book) {
	// 自定义排序逻辑
	sort.Slice(books, func(i, j int) bool {
		numA, errA := strconv.Atoi(strings.TrimSpace(books[i].Name)[:2])
		numB, errB := strconv.Atoi(strings.TrimSpace(books[j].Name)[:2])

		switch {
		case errA != nil && errB != nil: // 如果两者都不是数字
			return strings.ToLower(books[i].Name) < strings.ToLower(books[j].Name) // 按照字符串比较
		case errA != nil: // A不是数字
			return true // 非数字的元素排在数字元素后面
		case errB != nil: // B不是数字
			return false
		default: // 都是数字
			return numA < numB
		}
	})
}
