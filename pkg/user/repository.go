package user

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
	"github.com/jukeizu/treediagram/pkg/user/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_user"
)

type Repository interface {
	Preference(*userpb.PreferenceRequest) (*userpb.Preference, error)
	SetServer(*userpb.SetServerRequest) (*userpb.Preference, error)
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

func (r *repository) Preference(req *userpb.PreferenceRequest) (*userpb.Preference, error) {
	preference := &userpb.Preference{}

	q := `SELECT userId, serverId FROM preference WHERE userId = $1`

	err := r.Db.QueryRow(q, req.UserId).Scan(&preference.UserId, &preference.ServerId)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return preference, err
}

func (r *repository) SetServer(req *userpb.SetServerRequest) (*userpb.Preference, error) {
	preference := &userpb.Preference{}

	q := `INSERT INTO preference (userId, serverId) 
			VALUES($1, $2) 
			ON CONFLICT (userId) DO UPDATE SET serverId = excluded.serverId, updated = NOW()
			RETURNING userId, serverId`

	err := r.Db.QueryRow(q, req.UserId, req.ServerId).
		Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}
