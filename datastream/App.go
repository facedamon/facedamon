package main

import (
	"database/sql"
	"fmt"
	"github.com/facedamon/datastream/page"
	"github.com/facedamon/datastream/xmlparse"
	_ "github.com/go-sql-driver/mysql"
)

type container []map[string]interface{}

var c container

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func unit8string(u []uint8) string {
	r := make([]byte, 0)
	//r := []byte{}
	for _, b := range u {
		r = append(r, byte(b))
	}
	return string(r)
}

// when the type of sql in (char, varchar)
// then the go sql type is []uint8
// so this handler trans []inut8 to string
func FormatContainer(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch v.(type) {
		case []uint8:
			t := v.([]uint8)
			m[k] = unit8string(t)
		case nil:
			m[k] = ""
		}
	}
	return m
}

// parse the sql rows to []map[column]value
func ParseRows(rows *sql.Rows) container {
	col, err := rows.Columns()
	CheckErr(err)

	pts := make([]interface{}, len(col))
	in := make([]interface{}, len(col))

	for i := range pts {
		pts[i] = &in[i]
	}

	for rows.Next() {
		rows.Scan(pts...)
		m := make(map[string]interface{}, len(col))
		for i, column := range col {
			m[column] = in[i]
			m = FormatContainer(m)
		}
		c = append(c, m)
	}

	return c
}

func main() {
	//fmt.Println(os.Args[1])
	m := xmlparse.New("/Users/facedamon/go/src/facedamon/datastream/mysql2hive.xml")
	p := page.New()

	fmt.Println(m.SrcSelect())

	db, err := sql.Open(m.Src.SqlId, m.SrcUrl())
	CheckErr(err)
	defer db.Close()

	//计算总记录数
	count := 0
	row := db.QueryRow(m.SqlCount())
	row.Scan(&count)
	p.TotalRows = count
	p.SetTotalPages()

	var rows *sql.Rows

	for pa := p.CurrentPage; pa <= p.TotalPages; p.SetCurrentPage(pa) {
		p.SetStartIndex()
		p.SetLastIndex()
		selectSql, err := p.SqlLimit(m.SrcSelect(), m.SqlType())
		CheckErr(err)
		rows, err = db.Query(selectSql)
		CheckErr(err)
		r := ParseRows(rows)
		if len(r) != 0 {
			var args []interface{}
			for _, v := range r {
				for i := 0; i < len(v); i++ {
					args = append(args, v[m.Det.Table.Results[i].Column])
				}
			}
			sql.Open(m.Det.SqlId, m.DetUrl())
			CheckErr(err)
			result, err := db.Exec(m.DetInsert(), args...)
			CheckErr(err)
			fmt.Println(result.RowsAffected())
		}
		pa++
	}

	if rows != nil {
		defer rows.Close()
	}
}
