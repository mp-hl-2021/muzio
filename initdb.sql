-- DROP TABLE IF EXISTS accounts CASCADE;
CREATE TABLE accounts (
    id serial primary key,
    login varchar(255) not null,
    password varchar(255) not null,
    createdAt timestamp without time zone default now(),
    updatedAt timestamp without time zone default now(),

    unique(login)
);

CREATE TYPE link AS (
    serviceName varchar(255),
    url varchar(255),
    isAvailable boolean not null
);

-- DROP TABLE IF EXISTS entities CASCADE;
CREATE TABLE entities (
    id serial primary key,
    artist varchar(255) not null,
    album varchar(255),
    track varchar(255),
    links link[] not null
);

-- DROP TABLE IF EXISTS playlists CASCADE;
CREATE TABLE playlists (
    id serial primary key,
    owner integer not null,
    name varchar(255) not null,
    content integer[] not null,

    constraint fk_account foreign key (owner) references accounts(id)
);
