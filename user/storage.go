package user

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
	_ "github.com/jukeizu/treediagram/user/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	DatabaseName = "treediagram_user"
	TableName    = "preferences"
)

type UserStorage interface {
	Preference(*pb.PreferenceRequest) (*pb.Preference, error)
	SetServer(*pb.SetServerRequest) (*pb.Preference, error)
	Migrate() error
}

type storage struct {
	Db *sql.DB
}

func NewUserStorage(url string) (UserStorage, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	r := storage{
		Db: db,
	}

	return &r, err
}

func (s *storage) Migrate() error {
	return goose.Up(s.Db, "user/migrations")
}

func (s *storage) Preference(req *pb.PreferenceRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `SELECT userId, serverId FROM $1 WHERE userId = $2`

	err := s.Db.QueryRow(q, TableName, req.UserId).Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}

func (s *storage) SetServer(req *pb.SetServerRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	q := `INSERT INTO $1 (userId, serverId) VALUES($2, $3) RETURNING userId, serverId`

	err := s.Db.QueryRow(q, TableName, req.UserId, req.ServerId).
		Scan(&preference.UserId, &preference.ServerId)

	return preference, err
}
