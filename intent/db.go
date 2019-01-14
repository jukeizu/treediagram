package intent

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	_ "github.com/jukeizu/treediagram/intent/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName = "treediagram_intents"
)

type Intent struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Server   string        `json:"server"`
	Name     string        `json:"name"`
	Regex    string        `json:"regex"`
	Mention  bool          `json:"mention"`
	Response string        `json:"response"`
	Endpoint string        `json:"endpoint"`
	Help     string        `json:"help"`
	Enabled  bool          `json:"enabled"`
}

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
	//goose.AddMigration(migration.Up20190113072028, migration.Down20190113072028)
	return goose.Up(i.Db, "intent/migrations")
}

func (i *intentDb) Save(pbIntent pb.Intent) error {
	return nil
}

func (i *intentDb) Disable(id string) error {
	return nil
}

func (i *intentDb) Query(query pb.QueryIntentsRequest) ([]*pb.Intent, error) {
	pbIntents := []*pb.Intent{}

	return pbIntents, nil
}
