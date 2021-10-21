package data

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type User struct {
	Id              int64     `json:"id"`
	Email           string    `json:"email"`
	Pass            string    `json:"pass"`
	Name            string    `json:"name"`
	Surname         string    `json:"surname"`
	Birthdate       time.Time `json:"-"`
	BirthdateString string    `json:"birthdate"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetAll() ([]*User, error) {
	stmt := `SELECT id, email, pass, name, surname, birthdate FROM users ORDER BY id;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		err = rows.Scan(&u.Id, &u.Email, &u.Pass, &u.Name, &u.Surname, &u.Birthdate)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	stmt := `SELECT id, email, pass, name, surname, birthdate FROM users WHERE id = ?;`
	rows := m.DB.QueryRow(stmt, id)

	u := &User{}
	err := rows.Scan(&u.Id, &u.Email, &u.Pass, &u.Name, &u.Surname, &u.Birthdate)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return u, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	stmt := `SELECT id, email, pass, name, surname, birthdate FROM users WHERE email = ?;`
	rows := m.DB.QueryRow(stmt, email)

	u := &User{}
	err := rows.Scan(&u.Id, &u.Email, &u.Pass, &u.Name, &u.Surname, &u.Birthdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return u, nil
}

func (m *UserModel) Insert(email, pass, name, surname string, birthdate time.Time) (int, error) {
	stmt := `INSERT INTO users (email, pass, name, surname, birthdate) VALUES (?,?,?,?,?)`

	result, err := m.DB.Exec(stmt, email, pass, name, surname, birthdate.String())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Update(id int, email, name, surname string, birthdate time.Time) error {
	stmt := `UPDATE users SET email=?, name=?, surname=?, birthdate=? WHERE id = ?`

	result, err := m.DB.Exec(stmt, email, name, surname, birthdate, id)

	if err != nil {
		log.Printf("Unable to EDIT: %v\n", err)
		return err
	}

	if temp, _ := result.RowsAffected(); temp == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (m *UserModel) UpdatePassword(pass string, id int) error {
	stmt := `UPDATE users SET pass=? WHERE id = ?`

	result, err := m.DB.Exec(stmt, pass, id)

	if err != nil {
		log.Printf("Unable to EDIT: %v\n", err)
		return err
	}

	if temp, _ := result.RowsAffected(); temp == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (m *UserModel) Delete(id int) error {
	stmt := `DELETE FROM users WHERE id = ?;`

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
