package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type DBUser struct {
	Id       int64
	Username string
	Status   int
}

type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{
		db: db,
	}
}

func (m *DBManager) SetUser(u *DBUser) error {
	if m.db == nil {
		return errors.Errorf("database connection is nil")
	}

	query := `INSERT INTO users (user_id, username, status) VALUES ($1, $2, $3)`

	_, err := m.db.Exec(query, u.Id, u.Username, u.Status)
	if err != nil {
		return errors.Wrap(err, "err to insert user")
	}
	return nil
}

func (m *DBManager) GetUser(id int64) (*DBUser, error) {
	u := DBUser{}
	err := m.db.QueryRow("SELECT user_id, username, status FROM users WHERE user_id = $1", id).
		Scan(&u.Id, &u.Username, &u.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
