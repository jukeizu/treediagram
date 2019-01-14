package user

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	"github.com/jukeizu/treediagram/user/migrations"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_user"
)

type UserDb interface {
	Preference(*pb.PreferenceRequest) (*pb.Preference, error)
	SetServer(*pb.SetServerRequest) (*pb.Preference, error)
	Migrate() error
}

type userDb struct {
	Db     *sql.DB
	Logger zerolog.Logger
}

func NewUserDb(logger zerolog.Logger, url string) (UserDb, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	u := userDb{
		Db:     db,
		Logger: logger,
	}

	return &u, err
}

func (u *userDb) Migrate() error {
	g, err := gossage.New(u.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(migration.CreateTablePreferences20190113020925{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (u *userDb) Preference(req *pb.PreferenceRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `SELECT userId, serverId FROM preferences WHERE userId = $1`

	err := u.Db.QueryRow(q, req.UserId).Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}

func (u *userDb) SetServer(req *pb.SetServerRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `INSERT INTO preferences (userId, serverId) 
			VALUES($1, $2) 
			ON CONFLICT (userId) DO UPDATE SET serverId = excluded.serverId, updated = NOW()
			RETURNING userId, serverId`

	err := u.Db.QueryRow(q, req.UserId, req.ServerId).
		Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}
