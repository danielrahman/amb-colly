create table transactions
(
    product_id bigint                              not null
        primary key,
    quantity   bigint                              null,
    adjustment bigint                              null,
    date       timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);

