// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: clip.sql

package queries

import (
	"context"
)

const getClips = `-- name: GetClips :many
select
    clip.id,
    clip.title,
    clip.duration,
    clip.tape_id
from remix.clip
order by clip.created_at
`

type GetClipsRow struct {
	ID       string
	Title    string
	Duration int32
	TapeID   int32
}

func (q *Queries) GetClips(ctx context.Context) ([]GetClipsRow, error) {
	rows, err := q.db.QueryContext(ctx, getClips)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetClipsRow
	for rows.Next() {
		var i GetClipsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Duration,
			&i.TapeID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

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
