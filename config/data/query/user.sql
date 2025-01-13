-- name: FindUsers :many
select u.*
from user u
limit @limit
offset @offset;

-- name: FindUserByID :one
select u.*
from user u
where u.id = ?
limit 1;

-- name: FindUserByUsername :one
select u.*
from user u
where u.username = ?
limit 1;

-- name: InsertUser :one
insert into user (username, password, hash, active, role_id)
values (@username, @password, @hash, @active, @role_id)
returning id;

-- name: UpdateUser :exec
update user
set username = @username,
    password = @password,
    hash = @hash,
    active = @active,
    role_id = @roleID
where user.id = @id;

-- name: DeleteUser :exec
delete from user where id = @id;