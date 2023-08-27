create table segment
(
    id                        uuid primary key,
    slug                      varchar(300) not null,

    -- auto distribution is a percentage of the total users that will be
    -- automatically assigned to this segment.
    auto_distribution_percent int          not null default 0,

    constraint segment_slug_unique
        unique (slug)
);

create table "user"
(
    id int primary key
);

create table user_segment
(
    user_id    int       not null,
    segment_id uuid      not null,
    joined_at  timestamp not null default now(),
    left_at    timestamp          default null,
    due_at     timestamp          default null,

    constraint user_segment_user_id_fk
        foreign key (user_id) references "user" (id),
    constraint user_segment_segment_id_fk
        foreign key (segment_id) references segment (id)
);

create or replace function set_join_at_default()
    returns trigger as
$$
begin
    new.joined_at := now();
    return new;
end;
$$
    language plpgsql;

create trigger set_join_at_default
    before insert
    on user_segment
    for each row
    when ( new.joined_at is null )
execute procedure set_join_at_default();

create unique index user_segment_unique
    on user_segment (user_id, segment_id)
    where left_at is null;

insert into "user"
values (1),
       (2),
       (3),
       (4),
       (5),
       (6),
       (7),
       (8),
       (9),
       (10),
       (322);
