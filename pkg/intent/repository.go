package intent

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jukeizu/treediagram/api/protobuf-spec/intentpb"
	"github.com/jukeizu/treediagram/pkg/intent/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_intent"
)

type Repository interface {
	Save(*intentpb.Intent) error
	Disable(string) error
	Query(intentpb.QueryIntentsRequest) ([]*intentpb.Intent, error)
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

	err = g.RegisterMigrations(migrations.CreateTableIntent20190113072028{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) Save(pbIntent *intentpb.Intent) error {
	if pbIntent == nil {
		return nil
	}

	q := squirrel.Insert("intent").
		Columns(`
		serverId,
		name,
		regex,
		mention,
		response,
		endpoint,
		help,
		enabled`).
		Values(
			pbIntent.ServerId,
			pbIntent.Name,
			pbIntent.Regex,
			pbIntent.Mention,
			pbIntent.Response,
			pbIntent.Endpoint,
			pbIntent.Help,
			pbIntent.Enabled).
		Suffix("RETURNING id, created::INT")

	err := q.PlaceholderFormat(squirrel.Dollar).RunWith(r.Db).QueryRow().Scan(
		&pbIntent.Id,
		&pbIntent.Created,
	)

	return err
}

func (r *repository) Disable(id string) error {
	q := squirrel.Update("intent").Set("enabled", false).Where("id = ?", id)

	_, err := q.PlaceholderFormat(squirrel.Dollar).RunWith(r.Db).Exec()

	return err
}

func (r *repository) Query(query intentpb.QueryIntentsRequest) ([]*intentpb.Intent, error) {
	pbIntents := []*intentpb.Intent{}
	q := squirrel.Select(`
			id,
			serverId,
			name,
			regex,
			mention,
			response,
			endpoint,
			help,
			enabled,
			created::INT`).
		From("intent").
		Where("(serverId = ? OR serverId = '')", query.ServerId).
		Where("enabled = ?", true)

	if query.Name != "" {
		q = q.Where("name = ?", query.Name)
	}

	rows, err := q.PlaceholderFormat(squirrel.Dollar).RunWith(r.Db).Query()
	if err != nil {
		return pbIntents, err
	}

	defer rows.Close()
	for rows.Next() {
		pbIntent := intentpb.Intent{}
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
			&pbIntent.Created,
		)
		if err != nil {
			return pbIntents, err
		}

		pbIntents = append(pbIntents, &pbIntent)
	}

	return pbIntents, nil
}
