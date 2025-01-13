insert into role (id, name, rank)
values (1, 'mock-role-1', 1),
       (2, 'mock-role-2', 2),
       (3, 'mock-role-3', 0);

-- permissions

insert into permission (id, name, role_id)
values (1, 'read', 1),
       (2, 'update', 2),
       (3, 'create', 2);

insert into permission (id, name)
values (4, 'unused');

-- users

insert into user (id, username, password, hash, role_id)
values (1, 'mock-user-1', 'mock-pass-1', 'mock-hash-1', 1),
       (2, 'mock-user-2', 'mock-pass-2', 'mock-hash-2', 2);

insert into user (id, username, password, hash)
values (3, 'mock-user-3', 'mock-pass-3', 'mock-hash-3');
