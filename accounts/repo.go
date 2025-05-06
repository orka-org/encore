package accounts

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"encore.dev/storage/sqldb"
)

type UsersRepo interface {
	Create(ctx context.Context, username, email, passwordHash string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit int, offset int) ([]*User, error)
	Count(ctx context.Context) (int, error)
	Filter(ctx context.Context, filter *Filter) ([]*User, error)
}

type usersRepo struct {
	db *sqldb.Database
}

func NewUsersRepo(db *sqldb.Database) UsersRepo {
	return &usersRepo{
		db: db,
	}
}

func (r *usersRepo) Create(ctx context.Context, username, email, passwordHash string) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		INSERT INTO accounts (username, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, password, role, created_at, updated_at`,
		// INSERT INTO accounts (username, email, password, role, first_name, last_name)
		// VALUES ($1, $2, $3, $4, $5, $6)
		// RETURNING id, username, email, role, first_name, last_name, created_at, updated_at`,
		username, email, passwordHash, "user", // user.FirstName, user.LastName,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		// &user.FirstName,
		// &user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, role, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, password, phone, role, first_name, last_name, created_at, updated_at
		FROM accounts
		WHERE username = $1
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, email, role, created_at, updated_at
		FROM accounts
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) Update(ctx context.Context, user *User) error {
	err := r.db.QueryRow(ctx, `
		UPDATE accounts
		SET username = $1, email = $2, password = $3, role = $4
		WHERE id = $5
		RETURNING id, username, email, role, created_at, updated_at
	`, user.Username, user.Email, user.Password, user.Role, user.ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM accounts
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepo) List(ctx context.Context, limit int, offset int) ([]*User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, username, email, role, created_at, updated_at
		FROM accounts
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *usersRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM accounts
	`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type Filter struct {
	Username string // accepts regex, partial match and case-insensitive

	Email string // accepts regex, partial match and case-insensitive
	Phone string // accepts regex, partial match and case-insensitive

	CreatedBefore time.Time
	CreatedAfter  time.Time
	UpdatedBefore time.Time
	UpdatedAfter  time.Time

	Limit  int
	Offset int
}

func (r *usersRepo) Filter(ctx context.Context, filter *Filter) ([]*User, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM accounts
		WHERE 1=1
	`

	if filter.Username != "" {
		conditions = append(conditions, fmt.Sprintf(" AND username ILIKE $%d", argNum))
		args = append(args, "%"+filter.Username+"%")
		argNum++
	}

	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf(" AND email ILIKE $%d", argNum))
		args = append(args, "%"+filter.Email+"%")
		argNum++
	}

	if filter.Phone != "" {
		conditions = append(conditions, fmt.Sprintf(" AND phone ILIKE $%d", argNum))
		args = append(args, "%"+filter.Phone+"%")
		argNum++
	}

	if !filter.CreatedBefore.IsZero() {
		conditions = append(conditions, fmt.Sprintf(" AND created_at < $%d", argNum))
		args = append(args, filter.CreatedBefore)
		argNum++
	}

	if !filter.CreatedAfter.IsZero() {
		conditions = append(conditions, fmt.Sprintf(" AND created_at > $%d", argNum))
		args = append(args, filter.CreatedAfter)
		argNum++
	}

	if !filter.UpdatedBefore.IsZero() {
		conditions = append(conditions, fmt.Sprintf(" AND updated_at < $%d", argNum))
		args = append(args, filter.UpdatedBefore)
		argNum++
	}

	if !filter.UpdatedAfter.IsZero() {
		conditions = append(conditions, fmt.Sprintf(" AND updated_at > $%d", argNum))
		args = append(args, filter.UpdatedAfter)
		argNum++
	}

	query += strings.Join(conditions, "")
	query += `
		ORDER BY created_at DESC
		LIMIT $` + strconv.Itoa(argNum) + ` OFFSET $` + strconv.Itoa(argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
