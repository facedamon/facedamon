package sql

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

type ModelBaseWordInfo struct {
	UniqueNum     string         `db:"Unique_Num" json:"UniqueNum"`
	BaseWordId    string         `db:"Base_Word_Id" json:"BaseWordId"`
	BaseWordCn    string         `db:"Base_Word_Cn" json:"BaseWordCn"`
	BaseWordEn    string         `db:"Base_Word_En" json:"BaseWordEn"`
	Abbreviation1 string         `db:"Abbreviation1" json:"Abbreviation1"`
	Abbreviation2 string         `db:"Abbreviation2" json:"Abbreviation2"`
	UserId        string         `db:"User_Id" json:"UserId"`
	Creator       string         `db:"Creator" json:"Creator"`
	CreateTime    time.Time      `db:"Create_Time" json:"CreateTime"`
	Modifier      sql.NullString `db:"Modifier" json:"Modifier"`
	ModifyTime    sql.NullString `db:"Modify_Time" json:"ModifyTime"`
}

type MBWI string

func NewMbwi() MBWI {
	return "mbwi"
}

// queryByNum select ModelBaseWordInfo with UniqueNum = n
func (mbwi MBWI) QueryByNum(n string) (*ModelBaseWordInfo, error) {
	/*rows, err := db.Query("select * from Model_Base_Word_Info where Unique_Num = ?", n)
	if nil != err {
		log.Fatalf("Fail to get Model_Base_Word_Info data '%v'", err)
	}
	defer rows.Close()
	for rows.Next() {
		//there needs reflect to scan m
		err = rows.Scan(m)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Fatalf("there is no rows for Model_Base_Word_Info with Unique_Num=%s", n)
			}
			log.Fatalf("get Model_Base_Word_Info fail '%v'", err)
		}
	}*/
	m := ModelBaseWordInfo{}
	err := db.Get(&m, db.Rebind("select * from Model_Base_Word_Info where Unique_Num = ?"), n)
	if nil != err {
		log.Printf("Fail to get Model_Base_Word_Info data '%v'", err)
		return nil, err
	}
	return &m, nil
}

// Count statistical ModelBaseWorldInfo
func (mbwi MBWI) Count() (int, error) {
	n := 0
	err := db.QueryRow("select count(1) from Model_Base_Word_Info").Scan(&n)
	if nil != err {
		log.Printf("Fail to get Count Model_Base_Word_Info '%v'", err)
		return 0, err
	}
	return n, nil
}

// QueryByStruct fuzzy query by struct
func (mbwi MBWI) QueryByStruct(m ModelBaseWordInfo) ([]ModelBaseWordInfo, error) {
	sql, args := buildMbwiSql(m)
	var err error
	n := 0
	ms := []ModelBaseWordInfo{}
	if n, err = mbwi.Count(); err == nil && n > 0 {
		err = db.Select(&ms, db.Rebind(sql), args...)
		a := make([]string, 0)
		for _, p := range args {
			a = append(a, p.(string))
		}
		log.Println(db.Rebind(sql), strings.Join(a, ","))
	}
	if nil != err {
		return nil, err
	}
	if n <= 0 {
		return nil, errors.New("the table Model_Base_Word_Info is empty")
	}
	return ms, nil
}

func buildMbwiSql(m ModelBaseWordInfo) (string, []interface{}) {
	if (ModelBaseWordInfo{}) == m {
		return "", nil
	}
	buf := bytes.Buffer{}

	//args := make([]string, reflect.TypeOf(m).NumField())
	args := make([]interface{}, 0)
	buf.WriteString("select * from Model_Base_Word_Info where 1=1")
	if m.UniqueNum != "" {
		buf.WriteString(" and Unique_Num = ?")
		args = append(args, m.UniqueNum)
	}
	if m.BaseWordId != "" {
		buf.WriteString(" and Base_Word_Id like concat('%',?,'%')")
		args = append(args, m.BaseWordId)
	}
	if m.BaseWordCn != "" {
		buf.WriteString(" and Base_Word_Cn like concat('%',?,'%')")
		args = append(args, m.BaseWordCn)
	}
	if m.BaseWordEn != "" {
		buf.WriteString(" and Base_Word_En like concat('%',?,'%')")
		args = append(args, m.BaseWordEn)
	}
	if m.Abbreviation1 != "" {
		buf.WriteString(" and Abbreviation1 like concat('%',?,'%')")
		args = append(args, m.Abbreviation1)
	}
	if m.Abbreviation2 != "" {
		buf.WriteString(" and Abbreviation2 like concat('%',?,'%')")
		args = append(args, m.Abbreviation2)
	}
	if m.UserId != "" {
		buf.WriteString(" and User_Id  = ?")
		args = append(args, m.UserId)
	}
	if m.Creator != "" {
		buf.WriteString(" and Creator like concat('%',?,'%')")
		args = append(args, m.Creator)
	}
	if m.CreateTime != (time.Time{}) {
		buf.WriteString(" and Create_Time = ?")
		args = append(args, m.CreateTime.String())
	}
	if m.Modifier.Valid && m.Modifier.String != "" {
		buf.WriteString(" and Modifier like concat('%',?,'%')")
		args = append(args, m.Modifier.String)
	}
	if m.ModifyTime.Valid && m.ModifyTime.String != "" {
		buf.WriteString(" and Modify_Time = ?")
		args = append(args, m.ModifyTime.String)
	}
	return buf.String(), args
}
