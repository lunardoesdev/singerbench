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

-- name: SelectUnmeasuredProxies :many
select p.*
from proxies p
where not exists (
    select 1
    from measurements m
    where m.serverid = p.id
);

-- name: SelectBestProxies :many
with ranked as (
    select
        m.serverid,
        m.lastbyte,
        row_number() over (
            partition by m.serverid
            order by m.lastbyte
        ) as rn,
        count(*) over (
            partition by m.serverid
        ) as cnt
    from measurements m
),
middle_values as (
    select
        serverid,
        lastbyte
    from ranked
    where rn in ((cnt + 1) / 2, (cnt + 2) / 2)
),
medians as (
    select
        serverid,
        avg(lastbyte * 1.0) as median_lastbyte
    from middle_values
    group by serverid
)
select
    p.id,
    p.link,
    medians.median_lastbyte
from proxies p
join medians on medians.serverid = p.id
order by medians.median_lastbyte asc, p.id asc;
