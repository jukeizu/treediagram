package user

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	"github.com/jukeizu/treediagram/user/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_user"
)

type Repository interface {
	Preference(*pb.PreferenceRequest) (*pb.Preference, error)
	SetServer(*pb.SetServerRequest) (*pb.Preference, error)
	Migrate() error
}

type repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	r := repository{
		Db: db,
	}

	return &r, err
}

func (r *repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(migration.CreateTablePreference20190113020925{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) Preference(req *pb.PreferenceRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `SELECT userId, serverId FROM preference WHERE userId = $1`

	err := r.Db.QueryRow(q, req.UserId).Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}

func (r *repository) SetServer(req *pb.SetServerRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `INSERT INTO preference (userId, serverId) 
			VALUES($1, $2) 
			ON CONFLICT (userId) DO UPDATE SET serverId = excluded.serverId, updated = NOW()
			RETURNING userId, serverId`

	err := r.Db.QueryRow(q, req.UserId, req.ServerId).
		Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}
