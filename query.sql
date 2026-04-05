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

-- name: ListSubscriptoins :many
select * from subscriptions limit ? offset ?;

-- name: ListProxies :many
select * from proxies limit ? offset ?;

-- name: ListMeasurements :many
select * from measurements limit ? offset ?;

-- name: CountSubscriptions :one
select count(*) from subscriptions;

-- name: CountProxies :one
select count(*) from proxies;

-- name: CountMeasurements :one
select count(*) from measurements;

-- name: CountMeasurementsByProxy :one
select serverid, count(*)
from measurements where serverid = ?;