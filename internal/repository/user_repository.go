package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/serediukit/civix-backend/pkg/timeutil"

	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/pkg/database"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// Update(ctx context.Context, user *model.User) error
	// Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	store *database.Store
}

func NewUserRepository(store *database.Store) UserRepository {
	return &userRepository{store: store}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	createdAt := timeutil.Now()

	sql, args, err := db.SB().
		Insert(db.TableUsers).
		Columns(
			db.TableUsersColumnEmail,
			db.TableUsersColumnPasswordHash,
			db.TableUsersColumnName,
			db.TableUsersColumnCreatedAt,
			db.TableUsersColumnUpdatedAt,
		).
		Values(
			user.Email,
			user.PasswordHash,
			user.Name,
			createdAt,
			createdAt,
		).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, "Create user [%s] ToSQL: %s, %+v", user.Email, sql, args)
	}

	_, err = r.store.GetDB().Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "Create user [%s] Exec: %s, %+v", user.Email, sql, args)
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	sql, args, err := db.SB().
		Select(
			db.TableUsersColumnUserID,
			db.TableUsersColumnEmail,
			db.TableUsersColumnName,
			db.TableUsersColumnCreatedAt,
			db.TableUsersColumnUpdatedAt,
		).
		From(db.TableUsers).
		Where(db.TableUsersColumnUserID+" = ?", id).
		Where(db.TableUsersColumnDeletedAt + " IS NULL").
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Get user by id [%d]", id)
	}

	var user model.User

	err = r.store.GetDB().
		QueryRow(ctx, sql, args...).
		Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = db.ErrNotFound
		}

		return nil, errors.Wrapf(err, "Get user by id [%d]", id)
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	sql, args, err := db.SB().
		Select(
			db.TableUsersColumnUserID,
			db.TableUsersColumnEmail,
			db.TableUsersColumnPasswordHash,
			db.TableUsersColumnName,
			db.TableUsersColumnCreatedAt,
			db.TableUsersColumnUpdatedAt,
		).
		From(db.TableUsers).
		Where(db.TableUsersColumnEmail+" = ?", email).
		Where(db.TableUsersColumnDeletedAt + " IS NULL").
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Get user by email [%s]", email)
	}

	var user model.User

	err = r.store.GetDB().
		QueryRow(ctx, sql, args...).
		Scan(
			&user.UserID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = db.ErrNotFound
		}

		return nil, errors.Wrapf(err, "Get user by email [%s]", email)
	}

	return &user, nil
}
