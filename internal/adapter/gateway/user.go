package gateway

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/errors"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	entuser "github.com/morning-night-guild/platform-app/pkg/ent/user"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ repository.User = (*User)(nil)

type User struct {
	rdb *RDB
}

func NewUser(rdb *RDB) *User {
	return &User{
		rdb: rdb,
	}
}

func (usr *User) Save(
	ctx context.Context,
	item model.User,
) error {
	id := item.UserID.Value()

	now := time.Now().UTC()

	if err := usr.rdb.User.Create().
		SetID(id).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		OnConflict(
			sql.ConflictColumns(entuser.FieldID),
		).
		UpdateUpdatedAt().
		SetUpdatedAt(now).
		Exec(ctx); err != nil {
		log.GetLogCtx(ctx).Warn("failed to save user", log.ErrorField(err))

		return err
	}

	return nil
}

func (usr *User) Find(
	ctx context.Context,
	id user.ID,
) (model.User, error) {
	item, err := usr.rdb.User.Get(ctx, id.Value())
	if err != nil {
		return model.User{}, errors.NewNotFoundError("failed to find user", err)
	}

	return model.User{
		UserID: user.ID(item.ID),
	}, nil
}
