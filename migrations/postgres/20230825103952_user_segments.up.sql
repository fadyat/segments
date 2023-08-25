create table user_segment
(
    user_id    int  not null,
    segment_id uuid not null,

    constraint user_segment_pk primary key (user_id, segment_id),
    foreign key (segment_id) references segment (id)
);
