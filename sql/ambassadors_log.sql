create table log
(
    id     int unsigned auto_increment
        primary key,
    type   varchar(50)                         null,
    status varchar(50)                         null,
    date   timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);

INSERT INTO ambassadors.log (id, type, status, date) VALUES (1, 'products', 'start', '2019-10-25 16:55:38');