package xmlparse

import (
	"bytes"
	"encoding/xml"
	"facedamon/datastream/page"
	"io/ioutil"
	"log"
	"strings"
)

type mapper struct {
	XMLName xml.Name `xml:"mapper"`
	Src     src      `xml:"src"`
	Det     det      `xml:"det"`
}

type det struct {
	XMLName    xml.Name
	SqlId      string     `xml:"sqlId,attr"`
	DataSource datasource `xml:"datasource"`
	Table      table      `xml:"table"`
}

type src struct {
	XMLName    xml.Name
	SqlId      string     `xml:"sqlId,attr"`
	DataSource datasource `xml:"datasource"`
	Table      table      `xml:"table"`
}

type datasource struct {
	XMLName     xml.Name
	PoolName    string `xml:poolName,attr`
	ClassDriver string `xml:"class"`
	Url         string `xml:"url"`
	UserName    string `xml:"username"`
	Pwd         string `xml:"pwd"`
}

type table struct {
	XMLName    xml.Name
	TableName  string   `xml:"name,attr"`
	TableWhere string   `xml:"where"`
	Results    []result `xml:"result"`
}

type result struct {
	XMLName xml.Name
	Column  string `xml:"column,attr"`
	Type    string `xml:"jdbcType,attr"`
}

func New(path string) mapper {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	var mapper mapper
	err = xml.Unmarshal(file, &mapper)
	if err != nil {
		log.Fatalln(err)
	}
	return mapper
}

func (m mapper) SrcResults() string {
	r := m.Src.Table.Results
	str := make([]string, len(r), len(r))
	for i, v := range r {
		str[i] = v.Column
	}
	return strings.Join(str, ",")
}

func (m mapper) SrcSelect() string {
	s := m.SrcResults()
	t := m.Src.Table.TableName
	buf := bytes.Buffer{}
	buf.WriteString("select ")
	buf.WriteString(s)
	buf.WriteString(" from ")
	buf.WriteString(t)
	w := m.Src.Table.TableWhere
	if w != "" {
		buf.WriteString(" where ")
		buf.WriteString(w)
	}
	return buf.String()
}

func (m mapper) DetResults() string {
	r := m.Det.Table.Results
	str := make([]string, len(r), len(r))
	for i, v := range r {
		str[i] = v.Column
	}
	return strings.Join(str, ",")
}

func (m mapper) DetInsert() string {
	s := m.DetResults()
	t := m.Det.Table.TableName
	buf := bytes.Buffer{}
	buf.WriteString("insert into ")
	buf.WriteString(t)
	buf.WriteString(" (")
	buf.WriteString(s)
	buf.WriteString(") values(")
	ss := strings.Split(s, ",")
	for i := 0; i < len(ss); i++ {
		ss[i] = "?"
	}
	buf.WriteString(strings.Join(ss, ","))
	buf.WriteString(")")
	return buf.String()
}

func (m mapper) SrcUrl() string {
	return m.Src.DataSource.UserName + ":" + m.Src.DataSource.Pwd + m.Src.DataSource.Url
}

func (m mapper) DetUrl() string {
	return m.Det.DataSource.UserName + ":" + m.Det.DataSource.Pwd + m.Det.DataSource.Url
}

func (m mapper) SqlCount() string {
	return "select count(1) from " + m.Src.Table.TableName
}

func (m mapper) SqlType() page.SqlType {
	id := m.Src.SqlId

	my := page.SqlType{MYSQL: "mysql"}
	or := page.SqlType{ORACLE: "oracle"}
	hi := page.SqlType{HIVE: "hive"}

	if my.String() == id {
		return my
	}
	if or.String() == id {
		return or
	}
	if hi.String() == id {
		return hi
	}
	return my
}
