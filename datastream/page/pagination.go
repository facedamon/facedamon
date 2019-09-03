package page

import (
	"errors"
	"strconv"
	"strings"
)

const PAGE_SIZE = 100

type SqlType struct {
	MYSQL  string `enum:"mysql"`
	ORACLE string `enum:"oracle"`
	HIVE   string `enum:"hive"`
}

func (s SqlType) String() string {
	buf := make([]string, 0)
	if s.MYSQL != "" {
		buf = append(buf, s.MYSQL)
	}
	if s.ORACLE != "" {
		buf = append(buf, s.ORACLE)
	}
	if s.HIVE != "" {
		buf = append(buf, s.HIVE)
	}
	return strings.Join(buf, ",")
}

type pagination struct {
	TotalRows   int `总记录数`
	TotalPages  int `总页数`
	CurrentPage int `当前页码`
	StartIndex  int `起始行数`
	LastIndex   int `结束行数`
}

func New() *pagination {
	return &pagination{CurrentPage: 1}
}

func (p *pagination) SetCurrentPage(c int) {
	if c < 0 {
		p.CurrentPage = 1
	} else {
		p.CurrentPage = c
	}
}

//计算总页数
func (p *pagination) SetTotalPages() {
	if PAGE_SIZE == 0 {
		p.TotalPages = 0
	} else {
		if p.TotalRows%PAGE_SIZE == 0 {
			p.TotalPages = p.TotalRows / PAGE_SIZE
		} else {
			p.TotalPages = (p.TotalRows / PAGE_SIZE) + 1
		}
	}
}

func (p *pagination) SetStartIndex() {
	if p.CurrentPage > p.TotalPages {
		p.CurrentPage = p.TotalPages
	}
	p.StartIndex = (p.CurrentPage - 1) * PAGE_SIZE
}

func (p *pagination) SetLastIndex() {
	if PAGE_SIZE != 0 {
		if p.TotalRows < PAGE_SIZE {
			p.LastIndex = p.TotalRows
		} else if (p.TotalRows%PAGE_SIZE == 0) || (p.TotalRows%PAGE_SIZE != 0) && p.CurrentPage < p.TotalPages {
			p.LastIndex = p.CurrentPage * PAGE_SIZE
		} else if p.TotalRows%PAGE_SIZE != 0 && p.CurrentPage == p.TotalPages {
			//最后一页
			p.LastIndex = p.TotalRows
		}
	}
}

func (p *pagination) SqlLimit(sql string, sqlType SqlType) (string, error) {
	var mysql = SqlType{MYSQL: "mysql"}
	var oracle = SqlType{ORACLE: "oracle"}
	var hive = SqlType{HIVE: "hive"}

	switch sqlType {
	case mysql:
		return sql + " limit " + strconv.Itoa(p.StartIndex) + "," + strconv.Itoa(PAGE_SIZE), nil
	case hive:
		return sql + " limit " + strconv.Itoa(p.StartIndex) + "," + strconv.Itoa(PAGE_SIZE), nil
	case oracle:
		buf := make([]string, 0)
		buf = append(buf, sql)
		rear := append([]string{}, buf[6:]...)
		buf = append(append(buf[:6], " rownum rn "), rear...)
		if strings.Contains(sql, "where") {
			buf = append(buf, " rownum <= "+strconv.Itoa(p.LastIndex))
		} else {
			buf = append(buf, " where rownum <= "+strconv.Itoa(p.LastIndex))
		}
		buf2 := make([]string, 0)
		buf2 = append(buf2, "select * from (")
		buf2 = append(buf2, strings.Join(buf, ""))
		buf2 = append(buf2, ") tt where tt.rn > "+strconv.Itoa(p.StartIndex))
		return strings.Join(buf2, ""), nil
	default:
		return "", errors.New("there is no sqlType compatible with " + sqlType.String())
	}
}
