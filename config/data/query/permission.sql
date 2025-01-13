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
