package data

import (
	"database/sql"
	"errors"
)

type PostgresTestRepository struct {
	Conn    *sql.DB
	users   []*User
	idCount int
}

func NewPostgresTestRepository(db *sql.DB) PostgresTestRepository {
	return PostgresTestRepository{
		Conn:    db,
		users:   []*User{},
		idCount: 0,
	}
}

func (p PostgresTestRepository) GetAll() ([]*User, error) {
	return p.users, nil
}

func (p PostgresTestRepository) GetByEmail(email string) (*User, error) {
	for _, v := range p.users {
		if v.Email == email {
			return v, nil
		}
	}

	return &User{}, nil
}

func (p PostgresTestRepository) GetOne(id int) (*User, error) {
	for _, v := range p.users {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, errors.New("User not found")
}

func (p *PostgresTestRepository) Update(user *User) error {
	for i, v := range p.users {
		if v.ID == user.ID {
			p.users[i] = user

			return nil
		}
	}

	return errors.New("User not found")
}

func (p *PostgresTestRepository) Delete(id int) error {
	index := -1

	for i, v := range p.users {
		if v.ID == id {
			index = i
		}
	}

	if index == -1 {
		return errors.New("User not found")
	}

	tmp := p.users
	p.users = append(tmp[:index], tmp[index+1:]...)

	return nil
}

func (p *PostgresTestRepository) Insert(user *User) (int, error) {
	p.idCount++
	newId := p.idCount

	user.ID = newId

	p.users = append(p.users, user)

	return newId, nil
}

func (p *PostgresTestRepository) ResetPassword(password string, user *User) error {
	for i, v := range p.users {
		if v.ID == user.ID {
			user.Password = password
			p.users[i] = user

			return nil
		}
	}

	return errors.New("User not found")
}

func (p PostgresTestRepository) PasswordMatches(plainText string, user *User) (bool, error) {
	for _, v := range p.users {
		if v.ID == user.ID {
			matches := user.Password == plainText

			return matches, nil
		}
	}

	return true, nil
}
