package pg

import (
	"bytes"
	"database/sql/driver"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(driver, dataSource string) *sqlx.DB {
	db, err := sqlx.Connect(driver, dataSource)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

type Int64slice []int64

func (s Int64slice) Value() (driver.Value, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(strconv.FormatInt(val, 10))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return string(buffer.Bytes()), nil
}

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(strconv.Quote(val))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return string(buffer.Bytes()), nil
}
