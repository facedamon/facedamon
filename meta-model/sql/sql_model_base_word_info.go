package sql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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
		sqlLayoutWithArgs(sql, args)
	}
	if nil != err {
		return nil, err
	}
	if n <= 0 {
		return nil, errors.New("the table Model_Base_Word_Info is empty")
	}
	return ms, nil
}

func (mbwi MBWI) Inserts(ms []ModelBaseWordInfo) error {
	if len(ms) < 1 {
		return errors.New(fmt.Sprintf("Model_Base_Word_Info Inserts Should Not Be Exec, Because The Model Is Empty"))
	}
	tx, err := db.Begin()
	if nil != err {
		return errors.New(fmt.Sprintf("Model_Base_Word_Info Inserts Begin Transaction Fail '%v'", err))
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("insert into Model_Base_Word_Info() values(?,?,?,?,?,?,?,?,?,?,?)")
	if nil != err {
		return errors.New(fmt.Sprintf("Model_Base_Word_Info Inserts Prepare Fail '%v'", err))
	}
	for _, v := range ms {
		_, err = stmt.Exec(v.UniqueNum, v.BaseWordId, v.BaseWordCn, v.BaseWordEn, v.Abbreviation1, v.Abbreviation2, v.UserId, v.Creator, time.Now(), v.Modifier, v.ModifyTime)
		if nil != err {
			return errors.New(fmt.Sprintf("Model_Base_Word_Info Inserts Exec Fail '%v'", err))
		}
	}
	err = tx.Commit()
	if nil != err {
		return errors.New(fmt.Sprintf("Model_Base_Word_Info Insert Commit Fail '%v'", err))
	}
	stmt.Close()
	return nil
}

func (mbwi MBWI) Insert(m ModelBaseWordInfo) (string, error) {
	tx, err := db.Begin()
	if nil != err {
		log.Printf("Model_Base_Word_Info Insert Begin Transaction Fail '%v'", err)
		return "", err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("insert into Model_Base_Word_Info() values(?,?,?,?,?,?,?,?,?,?,?)")
	if nil != err {
		log.Printf("Model_Base_Word_Info Insert Prepare Fail '%v'", err)
		return "", err
	}
	res, err := stmt.Exec(m.UniqueNum, m.BaseWordId, m.BaseWordCn, m.BaseWordEn, m.Abbreviation1, m.Abbreviation2, m.UserId, m.Creator, time.Now(), m.Modifier, m.ModifyTime)
	if nil != err {
		log.Printf("Model_Base_Word_Info Insert Exec Fail '%v'", err)
		return "", err
	}
	num, _ := res.LastInsertId()
	err = tx.Commit()
	if nil != err {
		log.Printf("Model_Base_Word_Info Insert Commit Fail '%v'", err)
		return "", err
	}
	stmt.Close()
	return strconv.FormatInt(num, 10), nil
}

func (mbwi MBWI) UpdateByNum(m ModelBaseWordInfo) (string, error) {
	n := m.UniqueNum

	tx, err := db.Begin()
	if nil != err {
		log.Printf("Model_Base_Word_Info UpdateByNum(%s) Begin Transaction Fail '%v'", n, err)
		return "", err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("update Model_Base_Word_Info set (Base_Word_Id = ?, Base_Word_Cn = ?, Base_Word_En = ?, Abbreviation1 = ?, Abbreviation2 = ?, User_Id = ?,Modifier = ?, Modify_Time = ?) where Unique_Num = ?")
	if nil != err {
		log.Printf("Model_Base_Word_Info UpdateByNum(%s) Begin Prepare Fail '%v'", n, err)
		return "", err
	}
	res, err := stmt.Exec(m.BaseWordId, m.BaseWordCn, m.BaseWordEn, m.Abbreviation1, m.Abbreviation2, m.UserId, m.Modifier, time.Now(), n)
	if nil != err {
		log.Printf("Model_Base_Word_Info UpdateByNum(%s) Begin Result Fail '%v'", n, err)
		return "", err
	}
	num, _ := res.RowsAffected()
	err = tx.Commit()
	if nil != err {
		log.Printf("Model_Base_Word_Info UpdateByNum(%s) Commit Fail '%v'", n, err)
		return "", err
	}
	// *sql.Tx一旦释放,连接就回到连接池中，这里stmt在关闭时，必须在commit或rollback之前
	stmt.Close()
	return strconv.FormatInt(num, 10), nil
}

func sqlLayoutWithArgs(sql string, args []interface{}) {
	a := make([]string, 0)
	for _, p := range args {
		a = append(a, p.(string))
	}
	log.Println(db.Rebind(sql), strings.Join(a, ","))
}

func sqlLayout(sql string) {
	log.Println(db.Rebind(sql))
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
