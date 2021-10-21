package data

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Token struct {
	Id      int64
	User_id int64
	Token   string
	Expire  time.Time
}

type TokenModel struct {
	DB *sql.DB
}

func (m *TokenModel) GetByToken(token string) (*Token, error) {
	stmt := `SELECT id, user_id, token, expire FROM tokens WHERE token = ?`
	rows := m.DB.QueryRow(stmt, token)

	t := &Token{}
	err := rows.Scan(&t.Id, &t.User_id, &t.Token, &t.Expire)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return t, nil
}

func (m *TokenModel) GetByUserId(id int) (*Token, error) {
	stmt := `SELECT id, user_id, token, expire FROM tokens WHERE user_id = ?`
	rows := m.DB.QueryRow(stmt, id)

	t := &Token{}
	err := rows.Scan(&t.Id, &t.User_id, &t.Token, &t.Expire)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return t, nil
}

func (m *TokenModel) Insert(user_id int, token string, expire time.Time) (int, error) {
	stmt := `INSERT INTO tokens (user_id, token, expire) VALUES (?,?,?);`

	result, err := m.DB.Exec(stmt, user_id, token, expire)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TokenModel) Delete(id int) error {
	stmt := `DELETE FROM tokens WHERE id = ?;`

	result, err := m.DB.Exec(stmt, id)

	if err != nil {
		log.Printf("Unable to DELETE: %v\n", err)
		return err
	}

	if temp, _ := result.RowsAffected(); temp == 0 {
		// Work with Error
		return errors.New("no affected rows")
	}

	return nil

}
