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

func (pg PostgresRepo) Insert(graph domain.TrafficGraph) error {
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

	res, err := statement.Exec(graph.Reporter, _vertices, _edges, graph.PacketCount)
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

func (s PostgresRepo) FetchVertices() ([]domain.Vertex, error) {
	records := []domain.Vertex{}

	rows, err := database.Query("SELECT obj FROM vertices")
	if err != nil {
		return records, errors.New("failed retrieving VIEW [vertices]; " + err.Error())
	}

	for rows.Next() {
		var record string
		err := rows.Scan(&record)
		if err != nil {
			return records, errors.New("failed mapping vertices from storage; " + err.Error())
		}

		var vertex domain.Vertex
		err = json.Unmarshal([]byte(record), &vertex)
		if err != nil {
			return records, errors.New("failed unmarshalling to vertex object; " + err.Error())
		}

		records = append(records, vertex)
	}

	return records, nil
}

func (s PostgresRepo) FetchEdges() ([]domain.Edge, error) {
	records := []domain.Edge{}

	rows, err := database.Query("SELECT obj FROM edges_1min ORDER BY time desc limit 1")
	if err != nil {
		return records, errors.New("failed retrieving VIEW [edges_1m]; " + err.Error())
	}

	for rows.Next() {
		var record string
		err := rows.Scan(&record)
		if err != nil {
			return records, errors.New("failed mapping edges from storage; " + err.Error())
		}

		var edges []domain.Edge
		err = json.Unmarshal([]byte(record), &edges)
		if err != nil {
			return records, errors.New("failed unmarshalling to edge object; " + err.Error())

		}
		for _, edge := range edges {
			records = append(records, edge)
		}
	}

	return records, nil
}

func (s PostgresRepo) FetchHistory() ([]domain.GraphHistory, error) {
	records := []domain.GraphHistory{}

	rows, err := database.Query("SELECT time, sum FROM edges_1min ORDER BY time")
	if err != nil {
		return records, errors.New("failed retrieving VIEW [edges_1m]; " + err.Error())
	}

	for rows.Next() {
		var record domain.GraphHistory
		err := rows.Scan(&record.Time, &record.PacketCount)
		if err != nil {
			return records, errors.New("failed mapping records to GraphHistory object; " + err.Error())
		}

		records = append(records, record)
	}

	return records, nil
}
