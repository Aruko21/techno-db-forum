package apiserver

import (
	"net/http"
	"github.com/jackc/pgx"
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
	}

	dbConn, err := newDB(pgxConfig)
	if err != nil {
		return err
	}
	defer dbConn.Close()

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

	//if err := create.CreateTables(conn); err != nil {
	//	return nil, err
	//}

	return db, nil
}