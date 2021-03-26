package graphrepo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"netmap-bastion/internal/domain"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DbLocation string
}

var database *sql.DB

func NewPostgresRepo(loc string) (*PostgresRepo, error) {
	pg := new(PostgresRepo)
	pg.DbLocation = loc
	err := pg.Init()
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (pg PostgresRepo) Init() error {
	var err error
	database, err = sql.Open("postgres", pg.DbLocation)
	if err != nil {
		return err
	}
	return nil
}

func (pg PostgresRepo) Destroy() error {
	err := database.Close()
	if err != nil {
		return err
	}
	return nil
}

func (pg PostgresRepo) Insert(reporter string, graph domain.TrafficGraph) error {
	statement, err := database.Prepare("INSERT INTO graphs (time,reporter,vertices,edges,packets) values(NOW(),$1,$2,$3,$4) RETURNING time")
	if err != nil {
		return errors.New("failed database upsert; " + err.Error())
	}
	// var _reporter, _vertices, _edges, _packets sql.NullString
	// _reporter = sql.NullString{String: reporter, Valid: (reporter != "")}
	// _vertices = sql.NullString{String: graph.Vertices, Valid: (graph.Password != "")}
	// _edges = sql.NullString{String: graph.Email, Valid: (graph.Email != "")}
	// _packets = sql.NullString{String: graph.Email, Valid: (graph.Email != "")}

	_vertices, err := json.Marshal(graph.Vertices)
	if err != nil {
		return errors.New("failed json marshalling; " + err.Error())
	}

	_edges, err := json.Marshal(graph.Edges)
	if err != nil {
		return errors.New("failed json marshalling; " + err.Error())
	}

	res, err := statement.Exec(reporter, _vertices, _edges, graph.Properties.PacketCount)
	if err != nil {
		return err
	}
	var affected string
	a, err := res.RowsAffected()
	if err != nil {
		affected = err.Error()
	}
	affected = string(fmt.Sprint(a))
	fmt.Println("DEBUG: inserted record; rows affected: " + affected)
	return nil
}
