package sqlite

import (
	"database/sql"
	"errors"
	"strings"

	"blog/internal/domain"
	"blog/internal/infrastructure/persistence/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) All() ([]domain.User, error) {
	var dbUsers []models.User
	err := r.db.Select(&dbUsers, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	users := dbUsersToDomainUsers(dbUsers)
	return users, nil
}

func (r UserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	var dbUser models.User
	err := r.db.Get(&dbUser, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user := dbUserToDomainUser(dbUser)
	return user, nil
}

func (r UserRepository) FindByEmail(email string) (*domain.User, error) {
	var dbUser models.User
	err := r.db.Get(&dbUser, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user := dbUserToDomainUser(dbUser)
	return user, nil
}

func (r UserRepository) FindByUsername(username string) (*domain.User, error) {
	var dbUser models.User
	err := r.db.Get(&dbUser, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user := dbUserToDomainUser(dbUser)
	return user, nil
}

func (r UserRepository) Exists(id domain.UserID) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r UserRepository) UsernameExists(username string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE username=?", username)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r UserRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email=?", email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r UserRepository) Create(user *domain.User) (*domain.User, error) {
	rolesStr := rolesToString(user.UserRoles())

	_, err := r.db.Exec(`
		INSERT INTO 
		users (id, email, username, password_hash, description, user_roles, join_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		user.GetID().String(),
		user.Email(),
		user.Username(),
		user.PasswordHash(),
		user.Description(),
		rolesStr,
		user.JoinDate(),
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r UserRepository) UpdateRoles(id domain.UserID, roles []domain.UserRole) error {
	rolesStr := rolesToString(roles)

	_, err := r.db.Exec(`
		UPDATE users
		SET user_roles = ?
		WHERE id = ?
	`,
		rolesStr,
		id.String(),
	)
	return err
}

func (r UserRepository) UpdateDescription(id domain.UserID, newDescription string) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET description = ?
		WHERE id = ?
	`,
		newDescription,
		id.String(),
	)
	return err
}

func (r UserRepository) UpdatePasswordHash(id domain.UserID, passwordHash string) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET password_hash = ?
		WHERE id = ?
	`,
		passwordHash,
		id.String(),
	)
	return err
}

func dbUserToDomainUser(dbUser models.User) *domain.User {
	roles := stringToRoles(dbUser.UserRoles)

	return domain.RebuildUser(
		domain.NewUserID(dbUser.ID),
		dbUser.Email,
		dbUser.PasswordHash,
		dbUser.Username,
		dbUser.Description,
		roles,
		dbUser.JoinDate,
	)
}

func dbUsersToDomainUsers(dbUsers []models.User) []domain.User {
	users := []domain.User{}
	for _, user := range dbUsers {
		users = append(users, *dbUserToDomainUser(user))
	}
	return users
}

func rolesToString(roles []domain.UserRole) string {
	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = string(role)
	}
	return strings.Join(roleStrs, ";")
}

func stringToRoles(rolesStr string) []domain.UserRole {
	if rolesStr == "" {
		return []domain.UserRole{}
	}

	roleStrs := strings.Split(rolesStr, ";")
	roles := make([]domain.UserRole, 0, len(roleStrs))

	for _, roleStr := range roleStrs {
		roleStr = strings.TrimSpace(roleStr)
		if roleStr != "" {
			roles = append(roles, domain.UserRole(roleStr))
		}
	}

	return roles
}

