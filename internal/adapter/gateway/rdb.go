package gateway

import (
	"context"
	"database/sql"

	"github.com/morning-night-guild/platform-app/pkg/ent"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/pkg/errors"
)

// RDB RDBクライアント.
type RDB struct {
	*ent.Client
}

// RDBFactory RDBクライアントのファクトリ.
type RDBFactory interface {
	Of(dsn string) (*RDB, error)
}

// IsDuplicatedError 重複エラーであるかを判定する関数.
func (r *RDB) IsDuplicatedError(ctx context.Context, err error) bool {
	// https://github.com/ent/ent/issues/2176 により、
	// on conflict do nothingとしてもerror no rowsが返るため、個別にハンドリングする
	if errors.Is(err, sql.ErrNoRows) {
		log := log.GetLogCtx(ctx)

		log.Debug(err.Error())

		return true
	}

	return false
}
