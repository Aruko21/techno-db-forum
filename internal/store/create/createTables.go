package create

import (
	"fmt"
	"github.com/jackc/pgx"
	"io/ioutil"
)

func CreateTables(db *pgx.ConnPool) error {
	//_ = dropAllTables(db)

	file, err := ioutil.ReadFile("./internal/store/create/init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return err
	} else {
		fmt.Println("init.sql done!")
	}

	return nil
}

func dropAllTables(db *pgx.ConnPool) error {
	query := `DROP TABLE IF EXISTS forums;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}