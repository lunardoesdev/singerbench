create table if not exists proxies (
    id integer primary key,
    link text
);

create table if not exists measurements (
    datewhen integer,
    serverid integer,
    firstbyte integer,
    lastbyte integer,
    ping integer
);

