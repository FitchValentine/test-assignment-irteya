package repository

import (
	"database/sql"
	"ta/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserDB struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Age       int
	IsMarried bool
	Password  string
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, firstname, lastname, age, is_married, password) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, user.ID, user.Firstname, user.Lastname, user.Age, user.IsMarried, user.Password)
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, firstname, lastname, age, is_married, password FROM users WHERE id = $1`
	var dbUser UserDB
	err := r.db.QueryRow(query, id).Scan(
		&dbUser.ID, &dbUser.Firstname, &dbUser.Lastname, &dbUser.Age, &dbUser.IsMarried, &dbUser.Password,
	)
	if err != nil {
		return nil, err
	}
	return toUserDomain(&dbUser), nil
}

func toUserDomain(db *UserDB) *domain.User {
	return &domain.User{
		ID:        db.ID,
		Firstname: db.Firstname,
		Lastname:  db.Lastname,
		Fullname:  db.Firstname + " " + db.Lastname,
		Age:       db.Age,
		IsMarried: db.IsMarried,
		Password:  db.Password,
	}
}

