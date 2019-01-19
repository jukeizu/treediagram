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
	DatabaseName = "treediagram_intents"
)

type IntentDb interface {
	Save(pb.Intent) error
	Disable(string) error
	Query(pb.QueryIntentsRequest) ([]*pb.Intent, error)
	Migrate() error
}

type intentDb struct {
	Db *sql.DB
}

func NewIntentDb(url string) (IntentDb, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	i := intentDb{
		Db: db,
	}

	return &i, err
}

func (i *intentDb) Migrate() error {
	g, err := gossage.New(i.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(migration.CreateTableIntents20190113072028{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (i *intentDb) Save(pbIntent pb.Intent) error {
	q := `INSERT INTO intents (
		serverId,
		name,
		regex,
		mention,
		response,
		endpoint,
		help,
		enabled
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := i.Db.Exec(q,
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

func (i *intentDb) Disable(id string) error {
	q := `UPDATE intents SET enabled = false WHERE id = $1`

	_, err := i.Db.Exec(q, id)

	return err
}

func (i *intentDb) Query(query pb.QueryIntentsRequest) ([]*pb.Intent, error) {
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
		FROM intents 
		WHERE serverId = $1 OR serverId = '' 
		AND enabled = true`

	rows, err := i.Db.Query(q, query.ServerId)
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
