-- name: SyncClip :exec
insert into remix.clip (
    id,
    title,
    duration,
    tape_id,
    created_at
) values (
    @id,
    @title,
    @duration,
    @tape_id,
    now()
)
on conflict (id) do update set
    title = excluded.title,
    duration = excluded.duration,
    tape_id = excluded.tape_id;

-- name: GetClips :many
select
    clip.id,
    clip.title,
    clip.duration,
    clip.tape_id
from remix.clip
order by clip.created_at;
