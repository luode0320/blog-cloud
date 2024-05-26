package service

import (
	"errors"
	"md/dao"
	"md/middleware"
	"md/model/common"
	"md/model/entity"
	"md/util"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func RefreshDocument(document entity.Document, parentName string) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	user, err := dao.UserGetByName(middleware.Db, "admin")
	if err != nil {
		return
	}

	document.Id = util.SnowflakeString()
	document.Type = entity.DocMd
	document.CreateTime = util.CreateStamp()
	document.UpdateTime = util.CreateStamp()
	document.UserId = user.Id
	document.Name = strings.TrimSpace(document.Name)
	// 根据名称查询文档是否存在
	docs, err := dao.DocumentGetName(middleware.Db, document.Name, document.UserId)
	if err != nil {
		return
	}
	if len(docs) > 0 {
		return
	}

	if parentName != "" {
		// 根据名称查询一目录列表
		parentName = strings.TrimSpace(parentName)
		books, err := dao.BookListByName(tx, parentName, document.UserId)
		if err != nil {
			return
		}
		if len(books) == 0 {
			return
		}
		document.BookId = books[0].Id
	}

	err = dao.DocumentAdd(tx, document)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}
}

// 添加文档
func DocumentAdd(document entity.Document) entity.Document {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	document.Name = strings.TrimSpace(document.Name)
	if document.Name == "" {
		panic(common.NewError("文档名称不可为空"))
	}
	if util.StringLength(document.Name) > 100 {
		panic(common.NewError("文档名称过长, 请小于100个字符"))
	}
	if util.StringLength(document.Content) > 10000000 {
		panic(common.NewError("文档内容过多，请小于1000万个字符"))
	}
	if document.BookId == "" {
		panic(common.NewErr("请先选择的目录", errors.New("请先选择的目录")))
	}
	if document.Type != entity.DocMd {
		panic(common.NewError("不支持的文档类型"))
	}

	docs, err := dao.DocumentGetName(middleware.Db, document.Name, document.UserId)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}
	if len(docs) > 0 {
		panic(common.NewError("已存在同名文档"))
	}

	document.Id = util.SnowflakeString()
	document.CreateTime = util.CreateStamp()
	document.UpdateTime = util.CreateStamp()
	err = dao.DocumentAdd(tx, document)
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("添加失败", err))
	}

	go func() {
		book := Book(document.BookId)
		var rootBook entity.Book
		if book.ParentId != "" {
			rootBook = Book(book.ParentId)
		}

		// 生成文件
		filePath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
		util.CreateFile(filePath, document.Name+entity.MdExt, []byte(""))
		util.RefreshDir()
	}()

	middleware.Log.Infof("添加文档成功: {%s}", document.Name)
	return document
}

// 修改文档基础信息
func DocumentUpdate(document entity.Document) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	book := Book(document.BookId)

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

	go func() {
		var rootBook entity.Book
		if book.ParentId != "" {
			rootBook = Book(book.ParentId)
		}

		// 重命名
		dirPath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
		util.RenameFile(dirPath, doc.Name+entity.MdExt, document.Name+entity.MdExt)
		util.RefreshDir()
	}()

	middleware.Log.Infof("成功更新文档基础信息: {%s}", document.Name)
}

// 修改文档内容
func DocumentUpdateContent(document entity.Document) entity.Document {
	doc := DocumentGet(document.Id, document.UserId)
	book := Book(doc.BookId)
	var rootBook entity.Book
	if book.ParentId != "" {
		rootBook = Book(book.ParentId)
	}

	// 正则表达式模式，匹配图片URL
	pattern := `\((https?://[^)]*/` + common.PictureName + `/[^"\s]+)\)`
	re := regexp.MustCompile(pattern)

	// 用于存储替换后的结果
	var modifiedContent strings.Builder

	lastEnd := 0
	matches := re.FindAllStringIndex(document.Content, -1)
	for _, match := range matches {
		start, end := match[0], match[1]

		// 添加未匹配部分
		modifiedContent.WriteString(document.Content[lastEnd:start])

		// 处理匹配的图片URL
		matchStr := document.Content[start:end]
		splitURL := strings.Split(matchStr, "/"+common.PictureName+"/")
		if len(splitURL) > 1 {
			if book.ParentId != "" {
				modifiedURL := "(../../" + common.PictureName + "/" + strings.Join(splitURL[1:], "")
				modifiedContent.WriteString(modifiedURL)
			} else {
				modifiedURL := "(../" + common.PictureName + "/" + strings.Join(splitURL[1:], "")
				modifiedContent.WriteString(modifiedURL)
			}
		} else {
			// 如果没有找到/picture，原样保留
			modifiedContent.WriteString(matchStr)
		}

		// 更新lastEnd为当前匹配的结束位置
		lastEnd = end
	}

	// 添加剩余内容
	modifiedContent.WriteString(document.Content[lastEnd:])
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	document.Content = modifiedContent.String()

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

	go func() {
		// 将文档写入markdown文件
		filePath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
		util.CreateFile(filePath, doc.Name+entity.MdExt, []byte(document.Content))
		util.RefreshDir()
	}()

	middleware.Log.Infof("成功更新文档内容: {%s}", doc.Name)
	return document
}

// 删除文档
func DocumentDelete(id, userId string) {
	tx := middleware.DbW.MustBegin()
	defer tx.Rollback()

	doc := DocumentGet(id, userId)

	err := dao.DocumentDeleteById(tx, id, userId)
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(common.NewErr("删除失败", err))
	}

	go func() {
		book := Book(doc.BookId)
		var rootBook entity.Book
		if book.ParentId != "" {
			rootBook = Book(book.ParentId)
		}

		// 删除文档
		filePath := filepath.Join(common.DataPath, common.ResourceName, rootBook.Name, book.Name)
		util.RemoveFile(filePath, doc.Name+entity.MdExt)
		util.RefreshDir()
	}()

	middleware.Log.Infof("成功删除文档: {%s}", id)
}

// 查询文档列表
func DocumentList(bookId, userId string) []entity.Document {
	if bookId == "" {
		return []entity.Document{}
	}

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
