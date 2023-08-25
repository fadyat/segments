create table segment
(
    id   uuid primary key,
    slug varchar(300) not null,

    constraint segment_slug_unique
        unique (slug)
)