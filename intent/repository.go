package intent

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	"github.com/jukeizu/treediagram/intent/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_intent"
)

type Repository interface {
	Save(pb.Intent) error
	Disable(string) error
	Query(pb.QueryIntentsRequest) ([]*pb.Intent, error)
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

	i := repository{
		Db: db,
	}

	return &i, err
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

	err = g.RegisterMigrations(migration.CreateTableIntent20190113072028{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) Save(pbIntent pb.Intent) error {
	q := `INSERT INTO intent (
		serverId,
		name,
		regex,
		mention,
		response,
		endpoint,
		help,
		enabled
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := r.Db.Exec(q,
		pbIntent.ServerId,
		pbIntent.Name,
		pbIntent.Regex,
		pbIntent.Mention,
		pbIntent.Response,
		pbIntent.Endpoint,
		pbIntent.Help,
		pbIntent.Enabled,
	)

	return err
}

func (r *repository) Disable(id string) error {
	q := `UPDATE intent SET enabled = false WHERE id = $1`

	_, err := r.Db.Exec(q, id)

	return err
}

func (r *repository) Query(query pb.QueryIntentsRequest) ([]*pb.Intent, error) {
	pbIntents := []*pb.Intent{}

	q := `SELECT id,
			serverId,
			name,
			regex,
			mention,
			response,
			endpoint,
			help,
			enabled
		FROM intent 
		WHERE serverId = $1 OR serverId = '' 
		AND enabled = true`

	rows, err := r.Db.Query(q, query.ServerId)
	if err == sql.ErrNoRows {
		return pbIntents, nil
	}
	if err != nil {
		return pbIntents, err
	}

	defer rows.Close()
	for rows.Next() {
		pbIntent := pb.Intent{}
		err := rows.Scan(
			&pbIntent.Id,
			&pbIntent.ServerId,
			&pbIntent.Name,
			&pbIntent.Regex,
			&pbIntent.Mention,
			&pbIntent.Response,
			&pbIntent.Endpoint,
			&pbIntent.Help,
			&pbIntent.Enabled,
		)
		if err != nil {
			return pbIntents, err
		}

		pbIntents = append(pbIntents, &pbIntent)
	}

	return pbIntents, nil
}
