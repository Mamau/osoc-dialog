package dialog

import (
	"context"
	"github.com/Masterminds/squirrel"
	"osoc-dialog/internal/entity"
	"osoc-dialog/pkg/mysql"
	"time"
)

const maxMessageLimit = 1000

type Repository struct {
	db      *mysql.DB
	proxyDB *mysql.ProxyMysql
}

func New(db *mysql.DB, proxyDB *mysql.ProxyMysql) *Repository {
	return &Repository{
		db:      db,
		proxyDB: proxyDB,
	}
}

func (r *Repository) Save(ctx context.Context, message entity.Message) error {
	sqlQuery, args, err := r.proxyDB.Builder.Insert("messages").
		Columns("user_id", "author_id", "text", "created_at").
		Values(message.UserID, message.AuthorID, message.Text, time.Now()).
		ToSql()

	if err != nil {
		return err
	}

	if _, err = r.proxyDB.ExecContext(ctx, sqlQuery, args...); err != nil {
		return err
	}

	return nil
}
func (r *Repository) GetList(ctx context.Context, authorID int, userID int) ([]entity.Message, error) {
	sqlQuery, args, err := r.proxyDB.Builder.
		Select("id", "text", "user_id", "author_id", "created_at").
		From("messages").
		Where(
			squirrel.Or{
				squirrel.Eq{"user_id": userID},
				squirrel.Eq{"author_id": userID},
			},
		).
		Where(
			squirrel.Or{
				squirrel.Eq{"user_id": authorID},
				squirrel.Eq{"author_id": authorID},
			},
		).
		Limit(maxMessageLimit).
		ToSql()

	if err != nil {
		return nil, err
	}

	data := make([]entity.Message, 0, maxMessageLimit)
	if err = r.proxyDB.SelectContext(ctx, &data, sqlQuery, args...); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, entity.ErrNotFound
	}
	return data, nil
}
