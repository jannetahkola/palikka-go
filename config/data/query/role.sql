-- name: FindRoles :many
select * from role_view;

-- name: FindRoleByID :many
select * from role_view where id = @id;

-- name: InsertRole :one
insert into role (name, rank)
values (@name, @rank)
returning id;

-- name: UpdateRole :exec
update role
set name = @name,
    rank = @rank
where id = @id
