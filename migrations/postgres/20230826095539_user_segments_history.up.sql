alter table user_segment
    add column joined_at timestamp not null default now();

alter table user_segment
    add column left_at timestamp default null;