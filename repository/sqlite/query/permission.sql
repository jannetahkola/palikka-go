-- name: FindPermissions :many
select p.*
from permission p
limit @limit
offset @offset;

-- name: FindPermissionsByRoleID :many
select p.*
from permission p
join role r on r.id = p.role_id
where r.rank <= (
    select rr.rank from role rr
    where rr.id = @role_id
);

-- name: InsertPermission :one
insert into permission (name, role_id)
values (@name, @role_id)
returning id;

-- name: UpdatePermission :exec
update permission
set role_id = @role_id
where id = @id;

-- name: DeletePermissions :exec
delete from permission where id > 0;
