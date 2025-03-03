package repositories

import (
    "database/sql"
    "userservices/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `INSERT INTO users (id, username, email, password, role) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.Role)
    return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, username, email, password, role FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// GetAll returns all users from the database
func (r *UserRepository) GetAll() ([]*models.User, error) {
    query := `SELECT id, username, email, role FROM users`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        user := &models.User{}
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}