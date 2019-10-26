create table products
(
    product_id bigint                              not null
        primary key,
    title      text                                null,
    category   text                                null,
    vendor     text                                null,
    quantity   bigint                              null,
    price      bigint                              null,
    code       text                                null,
    url        text                                null,
    date       timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);
