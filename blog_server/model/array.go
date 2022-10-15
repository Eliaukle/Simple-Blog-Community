package model

import (
	"database/sql/driver"
	"strings"
)

type Array []string

// Scan 从数据库读取数据后，对其进行处理，获得Go类型的变量
func (m *Array) Scan(val interface{}) error {
	s := val.([]uint8)
	ss := strings.Split(string(s), "|")
	*m = ss
	return nil
}

// Value 将数据存到数据库时，对数据进行处理，获得数据库支持的类型
func (m Array) Value() (driver.Value, error) {
	str := strings.Join(m, "|")
	return str, nil
}
