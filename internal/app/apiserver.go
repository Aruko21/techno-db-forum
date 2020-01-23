package apiserver

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/store/create"
	"net/http"
	"time"
)

func Start() error {
	config := NewConfig()

	server, err := NewServer(config)
	if err != nil {
		return err
	}

	pgxConfig := pgx.ConnConfig{
		User:              "technopark",
		Password:          "park",
		Host:              "localhost",
		Port:              5432,
		Database:          "db-forum",
		//RuntimeParams: map[string]string{
		//	"standard_conforming_strings": "on",
		//},
		//PreferSimpleProtocol: true,
	}

	dbConn, err := newDB(pgxConfig)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	server.ConfigureServer(dbConn)

	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(connectionConfig pgx.ConnConfig) (*pgx.ConnPool, error) {
	// Необходимо указать макимальное количество одновременных
	// соединений и время ожидания при занятости. Так как иначе база
	// будет деградировать. Также используется Pool коннектов, а не чистый коннект
	pgxConnPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: connectionConfig,
		MaxConnections: 25,
		AcquireTimeout: time.Minute * 20,
	}

	db, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		return nil, err
	}

	var result string
	err = db.QueryRow("SHOW fsync;").Scan(&result)

	fmt.Println("CHECK fsync: ", result)

	if err := create.CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}