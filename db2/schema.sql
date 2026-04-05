create table if not exists proxies (
    id integer primary key autoincrement,
    link text unique
);

create table if not exists measurements (
    id integer primary key autoincrement,
    datewhen integer,
    serverid integer,
    firstbyte integer,
    lastbyte integer,
    ping integer
);

create table if not exists subscriptions (
    id integer primary key autoincrement,
    link text unique
);