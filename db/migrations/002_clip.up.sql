begin;

create table remix.clip (
    id         text primary key,
    title      text not null,
    duration   integer not null,
    tape_id    integer not null,
    created_at timestamptz not null default now()
);

comment on table remix.clip is
    'A video clip that can be queued for playback on demand.';
comment on column remix.clip.id is
    'Unique identifier for the clip; should match its filename without extension.';
comment on column remix.clip.title is
    'User-facing title for this clip.';
comment on column remix.clip.duration is
    'Total duration of the clip, in seconds.';
comment on column remix.clip.tape_id is
    'ID of the tape from which this clip was pulled.';
comment on column remix.clip.created_at is
    'Time at which this clip was first registered.';

commit;
