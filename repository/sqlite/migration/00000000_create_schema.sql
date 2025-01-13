-- role

create table role
(
    id   integer primary key autoincrement,
    name text not null unique,
    rank int  not null unique
);

-- permission

create table permission
(
    id      integer primary key autoincrement,
    name    text not null unique,
    role_id int,
    foreign key (role_id) references role (id)
);

-- user

create table user
(
    id         integer primary key autoincrement,
    username   text not null unique,
    password   text not null,
    created_at text not null default current_timestamp,
    active     int  not null default 1,
    role_id    int,
    foreign key (role_id) references role (id)
);
