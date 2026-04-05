create table if not exists proxies (
    id integer primary key autoincrement,
    link text
);

create table if not exists measurements (
    id integer primary key autoincrement,
    datewhen integer,
    serverid integer,
    firstbyte integer,
    lastbyte integer,
    ping integer
);

