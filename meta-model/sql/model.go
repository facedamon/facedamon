package sql

import (
	"database/sql"
	"log"
	"time"
)

type ModelBaseWordInfo struct {
	UniqueNum     string         `db:"Unique_Num"`
	BaseWordId    string         `db:"Base_Word_Id"`
	BaseWordCn    string         `db:"Base_Word_Cn"`
	BaseWordEn    string         `db:"Base_Word_En"`
	Abbreviation1 string         `db:"Abbreviation1"`
	Abbreviation2 string         `db:"Abbreviation2"`
	UserId        string         `db:"User_Id"`
	Creator       string         `db:"Creator"`
	CreateTime    time.Time      `db:"Create_Time"`
	Modifier      sql.NullString `db:"Modifier"`
	Modify_Time   time.Time      `db:"Modify_Time"`
}

func NewModelBaseWordInfo() *ModelBaseWordInfo {
	return &ModelBaseWordInfo{}
}

func (m *ModelBaseWordInfo) QueryByNum(n string) {
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
	err := db.Get(m, "select * from Model_Base_Word_Info where Unique_Num = ?", n)
	if nil != err {
		log.Fatalf("Fail to get Model_Base_Word_Info data '%v'", err)
	}
}
