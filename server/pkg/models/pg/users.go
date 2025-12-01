package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/psat/shiksha-prathap/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Update(name, phone string, age int32, id int) error {
	fmt.Printf("In UserModel Update: name=%s, phone=%s, age=%d, id=%d\n", name, phone, age, id)

	stmt := `UPDATE users SET username = $1, age = $2, phone=$3 WHERE id = $4`
	_, err := m.DB.Exec(stmt, name, age, phone, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Insert(email, password string) (int, error) {
	var pqErr *pq.Error

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return -1, err
	}

	stmt := `INSERT INTO users (email, hashed_password) VALUES($1, $2) RETURNING id`

	var newID int
	err = m.DB.QueryRow(stmt, email, hashedPassword).Scan(&newID)

	if err != nil && errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return -1, models.ErrDuplicateEmail
	}

	return newID, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = $1", email)

	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, username, email, age, phone FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.Phone)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}
