-- views

create view role_view as
select r.*,
       p.id      as permission_id,
       p.name    as permission_name,
       p.role_id as permission_role_id
from role r
         join (select id, rank from role) as rr on rr.rank <= r.rank
         left join permission p on rr.id = p.role_id;