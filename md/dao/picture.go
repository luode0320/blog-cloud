package dao

import (
	"md/model/common"
	"md/model/entity"
	"md/util"

	"github.com/jmoiron/sqlx"
)

// 分页查询图片记录
func PicturePage(db *sqlx.DB, page common.Page, userId string) ([]entity.Picture, int, error) {
	sqlCompletion := util.SqlCompletion{}
	sqlCompletion.InitSql(`select id,name,path,size,create_time from t_picture`)
	sqlCompletion.Eq("user_id", userId, true)
	sqlCompletion.Order("create_time", false)
	sqlCompletion.Limit(page.Current, page.Size)

	// 查询分页数据
	result := []entity.Picture{}
	err := db.Select(&result, sqlCompletion.GetSql(), sqlCompletion.GetParams()...)
	if err != nil {
		return result, 0, err
	}

	// 查询总记录数
	countResult := common.CountResult{}
	err = db.Get(&countResult, sqlCompletion.GetCountSql(), sqlCompletion.GetCountParams()...)
	if err != nil {
		return result, 0, err
	}

	return result, countResult.Count, nil
}

// 根据id删除图片
func PictureDeleteById(tx *sqlx.Tx, id, userId string) error {
	sql := `delete from t_picture where id=$1 and user_id=$2`
	_, err := tx.Exec(sql, id, userId)
	return err
}

// 根据id查询图片
func PictureGetById(tx *sqlx.Tx, id, userId string) (entity.Picture, error) {
	sql := `select * from t_picture where id=$1 and user_id=$2`
	result := entity.Picture{}
	err := tx.Get(&result, sql, id, userId)
	return result, err
}

// 根据文件大小、hash值查询相同图片的数量
func PictureCountBySizeHash(tx *sqlx.Tx, size int64, hash string) (common.CountResult, error) {
	sql := `select count(*) as count from t_picture where size=$1 and hash=$2`
	result := common.CountResult{}
	err := tx.Get(&result, sql, size, hash)
	return result, err
}

// 根据文件大小、hash值查询相同图片
func PictureBySizeHash(db *sqlx.DB, size int64, hash string) ([]entity.Picture, error) {
	sql := `select * from t_picture where size=$1 and hash=$2`
	result := []entity.Picture{}
	err := db.Select(&result, sql, size, hash)
	return result, err
}

// 添加图片
func PictureAdd(tx *sqlx.Tx, picture entity.Picture) error {
	sql := `insert into t_picture (id,name,path,hash,size,create_time,user_id) values (:id,:name,:path,:hash,:size,:create_time,:user_id)`
	_, err := tx.NamedExec(sql, picture)
	return err
}
