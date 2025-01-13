insert into role (id, name, rank)
values (1, 'reader', 1);

insert into role (id, name, rank)
values (2, 'admin', 2);

insert into permission (id, name, role_id)
values

    -- feature: users
    (1, 'users_create', 2),
    (2, 'users_read', 1),
    (3, 'users_update', 2),
    (4, 'users_delete', 2),

    -- feature: roles
    (5, 'roles_create', 2),
    (6, 'roles_read', 1),
    (7, 'roles_update', 2),
    (8, 'roles_delete', 2);

insert into user (id, username, password, active, role_id)
values
    (1, 'adminuser', '$2a$10$lZgoXXYMDT0y3ZT60FS/9ObOBCaB3F8Oo8Ghh9tzjd3L3.s/AL/fe', true, 2),
    (2, 'readeruser', '$2a$10$lZgoXXYMDT0y3ZT60FS/9ObOBCaB3F8Oo8Ghh9tzjd3L3.s/AL/fe', true, 1);