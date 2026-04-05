-- name: SaveMeasure :exec
insert into measurements (serverid, datewhen, ping, firstbyte, lastbyte)
values (?, ?, ?, ?, ?);