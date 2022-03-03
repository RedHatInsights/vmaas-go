package database

import (
	"io/ioutil"
)

func ExecFile(filename string) error {
	sql, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	sqldb, _ := Db.DB()
	_, err = sqldb.Exec(string(sql))
	return err
}
