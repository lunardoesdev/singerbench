-- name: SaveMeasurement :exec
insert into measurements (serverid, datewhen, ping, firstbyte, lastbyte)
values (?, ?, ?, ?, ?);

-- name: RemoveMeasurement :exec
delete from measurements
where id = ?;

-- name: GetProxyIdByLink :one
select * from proxies where link = ?;

-- name: AddProxy :exec
insert into proxies(link) values(?);

-- name: RemoveProxy :exec
delete from proxies where link = ?;

-- name: AddSubscription :exec
insert into subscriptions(link) values(?);

-- name: GetSubscriptionIdByLink :one
select * from subscriptions where link = ?;

-- name: RemoveSubscription :exec
delete from subscriptions where link = ?;