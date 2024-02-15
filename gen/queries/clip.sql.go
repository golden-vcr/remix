// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: clip.sql

package queries

import (
	"context"
)

const syncClip = `-- name: SyncClip :exec
insert into remix.clip (
    id,
    title,
    duration,
    tape_id,
    created_at
) values (
    $1,
    $2,
    $3,
    $4,
    now()
)
on conflict (id) do update set
    title = excluded.title,
    duration = excluded.duration,
    tape_id = excluded.tape_id
`

type SyncClipParams struct {
	ID       string
	Title    string
	Duration int32
	TapeID   int32
}

func (q *Queries) SyncClip(ctx context.Context, arg SyncClipParams) error {
	_, err := q.db.ExecContext(ctx, syncClip,
		arg.ID,
		arg.Title,
		arg.Duration,
		arg.TapeID,
	)
	return err
}
