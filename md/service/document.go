package service

import (
	"errors"
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"os"
	"strings"
	"time"
)

// 添加文档
func DocumentAdd(document entity.Document) entity.Document {
	if document.BookId == "" {
		panic(common.NewErr("请先选择的文集", errors.New("请先选择的文集")))
	}

	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	document.Name = strings.TrimSpace(document.Name)
	if document.Name == "" {
		panic(common.NewError("文档名称不可为空"))
	}

	if util.StringLength(document.Name) > 1000 {
		panic(common.NewError("文档名称过长，请小于1000个字符"))
	}

	if util.StringLength(document.Content) > 10000000 {
		panic(common.NewError("文档内容过多，请小于1000万个字符"))
	}

	if document.Type != entity.DocMd && document.Type != entity.DocOpenApi {
		panic(common.NewError("不支持的文档类型"))
	}

	document.Id = util.SnowflakeString()
	document.CreateTime = time.Now().UnixMilli()
	document.UpdateTime = time.Now().UnixMilli()
	err := dao.DocumentAdd(tx, document)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	// 将文档写入markdown文件
	go func() {
		book := entity.Book{}
		if document.BookId == "" {
			book.Name = "其它"
		} else {
			book = Book(document.BookId)
		}

		// 创建目录
		fileName := document.Name + entity.MdExt
		filePath := common.DataPath + common.ResourceName + "/" + book.Name
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			middleware.Log.Errorf("创建目录失败: {%s}", err)
			return
		}

		// 生成文件: 文集名称 + 文档名称
		saveMdFile, err := os.Create(filePath + "/" + fileName)
		if err != nil {
			middleware.Log.Errorf("添加文档失败: {%s}", err)
		}
		defer saveMdFile.Close()

		middleware.Log.Infof("成功添加文档: {%s}", filePath)
	}()

	middleware.Log.Infof("添加文档成功: {%s}", document.Name)
	return document
}

// 修改文档基础信息
func DocumentUpdate(document entity.Document) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	document.Name = strings.TrimSpace(document.Name)
	if document.Name == "" {
		panic(common.NewError("文档名称不可为空"))
	}

	if util.StringLength(document.Name) > 1000 {
		panic(common.NewError("文档名称过长，请小于1000个字符"))
	}

	doc, err := dao.DocumentGetById(middleware.Db, document.Id, document.UserId)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	err = dao.DocumentUpdate(tx, document)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	// 将文档写入markdown文件
	go func() {
		book := Book(document.BookId)

		// 目录
		oldFileName := doc.Name + entity.MdExt
		newFileName := document.Name + entity.MdExt
		filePath := common.DataPath + common.ResourceName + "/" + book.Name
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			middleware.Log.Errorf("创建目录失败: {%s}", err)
			return
		}

		err = os.Rename(filePath+"/"+oldFileName, filePath+"/"+newFileName)
		if err != nil {
			middleware.Log.Errorf("重命名文档失败: {%s}", err)
			return
		}

		middleware.Log.Infof("成功更新文档基础信息: {%s}", filePath+"/"+newFileName)
	}()

	middleware.Log.Infof("成功更新文档基础信息: {%s}", document.Name)
}

// 修改文档内容
func DocumentUpdateContent(document entity.Document) entity.Document {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	if util.StringLength(document.Content) > 10000000 {
		panic(common.NewError("文档内容过多，请小于1000万个字符"))
	}

	document.UpdateTime = time.Now().UnixMilli()
	err := dao.DocumentUpdateContent(tx, document)
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("更新失败", err))
	}

	doc := DocumentGet(document.Id, document.UserId)

	// 将文档写入markdown文件
	go func() {
		book := entity.Book{}
		if doc.BookId == "" {
			book.Name = "其它"
		} else {
			book = Book(doc.BookId)
		}

		// 创建目录
		fileName := doc.Name + entity.MdExt
		filePath := common.DataPath + common.ResourceName + "/" + book.Name
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			middleware.Log.Errorf("创建目录失败: {%s}", err)
			return
		}

		// 生成文件: 文集名称 + 文档名称
		saveMdFile, err := os.Create(filePath + "/" + fileName)
		if err != nil {
			middleware.Log.Errorf("更新文档内容失败: {%s}", err)
			return
		}
		defer saveMdFile.Close()

		_, err = saveMdFile.Write([]byte(document.Content))
		if err != nil {
			middleware.Log.Errorf("更新文档内容发生错误: {%s}", err)
			return
		}

		middleware.Log.Infof("成功更新文档内容: {%s}", filePath)
	}()

	middleware.Log.Infof("成功更新文档内容: {%s}", doc.Name)
	return doc
}

// 删除文档
func DocumentDelete(id, userId string) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	err := dao.DocumentDeleteById(tx, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	doc := DocumentGet(id, userId)
	book := entity.Book{}
	if doc.BookId == "" {
		book.Name = "其它"
	} else {
		book = Book(doc.BookId)
	}

	fileName := doc.Name + entity.MdExt
	filePath := common.DataPath + common.ResourceName + "/" + book.Name
	// 删除文件: 文集名称 + 文档名称
	err = os.Remove(filePath + "/" + fileName)
	if err != nil {
		panic(common.NewErr("删除文件失败", err))
		return
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	middleware.Log.Infof("成功删除文档: {%s}", id)
}

// 查询文档列表
func DocumentList(bookId, userId string) []entity.Document {
	documents, err := dao.DocumentList(middleware.Db, bookId, userId)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}
	return documents
}

// 查询文档
func DocumentGet(id, userId string) entity.Document {
	document, err := dao.DocumentGetById(middleware.Db, id, userId)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}
	return document
}

// 查询公开发布文档
func DocumentGetPublished(id string) entity.Document {
	document, err := dao.DocumentGetPublished(middleware.Db, id)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}
	return document
}

// 分页查询公开发布文档列表
func DocumentPagePulished(pageCondition common.PageCondition[entity.DocumentPageCondition]) common.PageResult[entity.DocumentPageResult] {
	records, total, err := dao.DocumentPagePulished(middleware.Db, pageCondition)
	if err != nil {
		panic(common.NewErr("查询失败", err))
	}
	pageResult := common.PageResult[entity.DocumentPageResult]{Records: records, Total: total}
	return pageResult
}
