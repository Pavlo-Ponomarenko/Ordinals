-- +migrate Up
create table address (
    id varchar(255) primary key
);

create table inscription (
    id bigint primary key,
    content varchar(255),
    tx_hash varchar(255),
    date varchar(255),
    address_id varchar(255) references address(id)
);
-- +migrate Down
drop table inscription;
drop table address;
