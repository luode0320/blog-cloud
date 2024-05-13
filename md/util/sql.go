// sql语句拼接工具类
package util

import (
	"regexp"
	"strconv"
	"strings"
)

type SqlCompletion struct {
	initSql      string
	initCountSql string
	whereSql     strings.Builder
	groupSql     string
	havingSql    string
	orderSql     strings.Builder
	limitSql     string
	paramIndex   int
	whereParams  []interface{}
	limitParams  []interface{}
}

// InitSql 用于初始化 SqlCompletion 结构体的 SQL 语句和计数 SQL 语句
// 参数 sql 为要初始化的 SQL 语句
// 返回一个指向 SqlCompletion 结构体的指针
func (s *SqlCompletion) InitSql(sql string) *SqlCompletion {
	// 将传入的 SQL 语句赋值给结构体的 initSql 字段
	s.initSql = sql
	// 使用正则表达式编译一个匹配 select 语句的正则表达式
	r, _ := regexp.Compile("(?i)(?s)select(.*?)from")
	// 将 select 语句替换为 select count(*) as count from，生成计数 SQL 语句
	s.initCountSql = r.ReplaceAllString(sql, "select count(*) as count from")
	// 返回指向当前 SqlCompletion 结构体的指针
	return s
}

// 设置初始sql语句和总行数语句
func (s *SqlCompletion) InitSqlAndCount(sql, countSql string) *SqlCompletion {
	s.initSql = sql
	s.initCountSql = countSql
	return s
}

// 获取sql语句
func (s *SqlCompletion) GetSql() string {
	return s.initSql + s.whereSql.String() + s.groupSql + s.havingSql + s.orderSql.String() + s.limitSql
}

// 获取行数sql语句
func (s *SqlCompletion) GetCountSql() string {
	return s.initCountSql + s.whereSql.String() + s.groupSql + s.havingSql
}

// 获取条件列表
func (s *SqlCompletion) GetParams() []interface{} {
	return append(s.whereParams, s.limitParams...)
}

// 获取行数条件列表
func (s *SqlCompletion) GetCountParams() []interface{} {
	return s.whereParams
}

// =
func (s *SqlCompletion) Eq(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, "=", false, isAnd)
	return s
}

// !=
func (s *SqlCompletion) Ne(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, "!=", false, isAnd)
	return s
}

// >
func (s *SqlCompletion) Gt(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, ">", false, isAnd)
	return s
}

// <
func (s *SqlCompletion) Lt(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, "<", false, isAnd)
	return s
}

// >=
func (s *SqlCompletion) Ge(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, ">=", false, isAnd)
	return s
}

// <=
func (s *SqlCompletion) Le(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, "<=", false, isAnd)
	return s
}

// like
func (s *SqlCompletion) Like(field string, param interface{}, isAnd bool) *SqlCompletion {
	s.where(field, param, "like", true, isAnd)
	return s
}

// in
func (s *SqlCompletion) In(field string, params []interface{}, isAnd bool) *SqlCompletion {
	s.whereIn(field, params, "in", isAnd)
	return s
}

// not in
func (s *SqlCompletion) NotIn(field string, params []interface{}, isAnd bool) *SqlCompletion {
	s.whereIn(field, params, "not in", isAnd)
	return s
}

// is null
func (s *SqlCompletion) IsNull(field string, isAnd bool) *SqlCompletion {
	s.whereHyphen(isAnd)
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" is null")
	return s
}

// is not null
func (s *SqlCompletion) IsNotNull(field string, isAnd bool) *SqlCompletion {
	s.whereHyphen(isAnd)
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" is not null")
	return s
}

// 分组，只设置最后一次
func (s *SqlCompletion) Group(fields string) *SqlCompletion {
	s.groupSql = " group by " + fields
	return s
}

// having条件，只设置最后一次
func (s *SqlCompletion) Having(fields string) *SqlCompletion {
	s.havingSql = " having " + fields
	return s
}

// Order 方法用于构建 SQL 语句中的排序子句
// 参数 field 表示排序的字段名
// 参数 isAsc 表示是否按升序排序
// 返回一个指向当前 SqlCompletion 结构体的指针
func (s *SqlCompletion) Order(field string, isAsc bool) *SqlCompletion {
	// 判断是否需要添加逗号
	needComma := true

	// 如果 orderSql 字符串长度为0，表示还没有添加排序子句，需要先添加 "order by" 关键字
	if s.orderSql.Len() == 0 {
		s.orderSql.WriteString(" order by ")
		needComma = false
	}

	if needComma {
		s.orderSql.WriteString(",")
	}
	s.orderSql.WriteString(field)
	if !isAsc {
		s.orderSql.WriteString(" desc")
	}
	return s
}

// Limit 方法用于设置 SQL 语句的分页参数
// 参数 page 表示页码，从1开始
// 参数 size 表示每页的数据条数
// 返回一个指向当前 SqlCompletion 结构体的指针
func (s *SqlCompletion) Limit(page, size int) *SqlCompletion {
	// 如果页码和每页数据条数均大于0，则设置分页参数
	if page > 0 && size > 0 {
		// 计算偏移量
		offset := size * (page - 1)

		// 递增参数索引，并构建参数占位符
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder1 := "$" + strconv.Itoa(s.paramIndex)

		// 递增参数索引，并构建参数占位符
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder2 := "$" + strconv.Itoa(s.paramIndex)

		// 设置 limit 子句和参数
		s.limitSql = " limit " + paramPlaceholder1 + " offset " + paramPlaceholder2
		s.limitParams = []interface{}{size, offset}
	}

	// 返回指向当前 SqlCompletion 结构体的指针
	return s
}

// where 方法用于构建 SQL 语句中的 WHERE 子句
// 参数 field 表示字段名
// 参数 param 表示字段的值
// 参数 symbol 表示比较符号，如 "="、">"、"<" "like" 等
// 参数 isLike 表示是否使用模糊匹配, 如 "%s%"
// 参数 isAnd 表示是否使用 AND 连接多个条件
func (s *SqlCompletion) where(field string, param interface{}, symbol string, isLike bool, isAnd bool) {
	// 递增参数索引
	s.paramIndex = s.paramIndex + 1
	// 构建参数占位符
	paramPlaceholder := "$" + strconv.Itoa(s.paramIndex)
	// 根据 isAnd 参数决定是否添加 AND 连接符
	s.whereHyphen(isAnd)
	// 添加字段名、比较符号和参数占位符到 whereSql 字符串中
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" ")
	s.whereSql.WriteString(symbol)
	if isLike {
		// 如果 isLike 为 true，则使用模糊匹配，将参数占位符包裹在 % 符号中
		s.whereSql.WriteString(" '%'||" + paramPlaceholder + "||'%'")
	} else {
		// 如果 isLike 为 false，则直接使用参数占位符
		s.whereSql.WriteString(" " + paramPlaceholder)
	}

	// 将参数值添加到 whereParams 切片中
	s.whereParams = append(s.whereParams, param)
}

// 添加where条件（in / not in）
func (s *SqlCompletion) whereIn(field string, params []interface{}, symbol string, isAnd bool) {
	if len(params) == 0 {
		return
	}
	s.whereHyphen(isAnd)
	s.whereSql.WriteString(field)
	s.whereSql.WriteString(" ")
	s.whereSql.WriteString(symbol)
	s.whereSql.WriteString(" (")
	for i, v := range params {
		if i != 0 {
			s.whereSql.WriteString(",")
		}
		s.paramIndex = s.paramIndex + 1
		paramPlaceholder := "$" + strconv.Itoa(s.paramIndex)
		s.whereSql.WriteString(paramPlaceholder)
		s.whereParams = append(s.whereParams, v)
	}
	s.whereSql.WriteString(")")
}

// whereHyphen 方法用于拼接 WHERE 子句中的连接符（AND 或 OR）
// 参数 isAnd 表示是否使用 AND 连接符
func (s *SqlCompletion) whereHyphen(isAnd bool) {
	// 判断是否需要添加连接符
	needHyphen := true
	// 如果 whereSql 字符串长度为0，表示还没有添加 WHERE 子句，需要先添加 WHERE 关键字
	if s.whereSql.Len() == 0 {
		s.whereSql.WriteString(" where ")
		needHyphen = false
	}

	// 如果需要添加连接符
	if needHyphen {
		// 根据 isAnd 参数决定添加 AND 还是 OR 连接符
		if isAnd {
			s.whereSql.WriteString(" and ")
		} else {
			s.whereSql.WriteString(" or ")
		}
	}
}
