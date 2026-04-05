-- name: GetA :one
select * from proxies
where id = ?;