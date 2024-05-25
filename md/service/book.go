package service

import (
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"path/filepath"
	"strings"
	"time"
)

// 添加目录
func BookAdd(book entity.Book) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	book.Name = strings.TrimSpace(book.Name)
	if book.Name == "" {
		panic(common.NewError("目录名称不可为空"))
	}
	if util.StringLength(book.Name) > 100 {
		panic(common.NewError("目录名称过长, 请小于100个字符"))
	}

	// 根据名称查询一目录列表
	books, err := dao.BookListByName(tx, book.Name, book.UserId)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}
	if len(books) > 0 {
		panic(common.NewError("已存在同名目录"))
	}

	// 保存
	book.Id = util.SnowflakeString()
	book.CreateTime = time.Now().UnixMilli()
	err = dao.BookAdd(tx, book)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	go func() {
		var rootBook entity.Book
		if book.ParentId != "" {
			rootBook = Book(book.ParentId)
		}
		util.CreateDir(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
	}()

	middleware.Log.Infof("成功添加一级目录: {%s}", book.Name)
}

// 修改目录
func BookUpdate(book entity.Book) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	oldBook := Book(book.Id)

	book.Name = strings.TrimSpace(book.Name)
	if book.Name == "" {
		panic(common.NewError("目录名称不可为空"))
	}

	if util.StringLength(book.Name) > 100 {
		panic(common.NewError("目录名称过长，请小于100个字符"))
	}

	// 根据名称查询目录列表
	books, err := dao.BookListByName(tx, book.Name, book.UserId)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	for _, v := range books {
		if v.Id != book.Id {
			panic(common.NewError("已存在同名目录"))
		}
	}

	// 更新名称
	err = dao.BookUpdate(tx, book)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	go func() {
		var rootBook entity.Book
		if oldBook.ParentId != "" {
			rootBook = Book(oldBook.ParentId)
		}
		oldPath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, oldBook.Name)
		newPath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
		util.RenameDir(oldPath, newPath)
	}()

	middleware.Log.Infof("成功更新目录名称: {%s}", book.Name)
}

// 删除目录
func BookDelete(id, userId string) {
	documents, err := dao.DocumentList(middleware.Db, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}
	if len(documents) > 0 {
		panic(common.NewError("目录不为空, 无法删除"))
	}

	twoBooks, err := dao.BookByParentId(middleware.Db, userId, id)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}
	if len(twoBooks) > 0 {
		panic(common.NewError("目录不为空, 无法删除"))
	}

	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	book := Book(id)

	// 删除
	err = dao.BookDeleteById(tx, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	go func() {
		var rootBook entity.Book
		if book.ParentId != "" {
			rootBook = Book(book.ParentId)
		}
		util.RemoveDir(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
	}()

	middleware.Log.Infof("成功删除一级目录: {%s}", book.Name)
}

// 查询一级目录列表
func BookList(userId string) []entity.Book {
	books, err := dao.BookList(middleware.Db, userId)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	// 将全部加到首位
	books = append([]entity.Book{}, books...)
	return books
}

// 查询一级目录
func Book(id string) entity.Book {
	book, err := dao.Book(middleware.Db, id)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}

	return book
}
