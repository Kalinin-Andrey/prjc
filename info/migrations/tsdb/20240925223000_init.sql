-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

create schema cmc;

create table cmc.currency
(
    id                                  bigint                  not null,
    symbol                              text                    not null,
    slug                                text                    not null,
    name                                text                    not null,
    is_for_observing                    bool                    not null,
    CONSTRAINT currency__id__pk PRIMARY KEY (id) include (symbol, slug, name, is_for_observing)
);

create unique index currency__slug__ux ON cmc.currency (slug) include (id, symbol, name, is_for_observing);


create table cmc.import_max_time
(
    currency_id                         bigint                  not null,
    price_and_cap                       timestamp               null,
    concentration                       timestamp               null,
    CONSTRAINT import_max_time__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index import_max_time__currency_id__pk ON cmc.import_max_time (currency_id) include (price_and_cap, concentration);


create table cmc.price_and_cap
(
    currency_id                         bigint                  not null,
    price                               double precision        not null,
    daily_volume                        double precision        not null,
    cap                                 double precision        not null,
    ts                                  timestamp               not null,
    CONSTRAINT price_and_cap__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index price_and_cap__currency_id__ts__pk ON cmc.price_and_cap (currency_id, ts) include (price, daily_volume, cap);

select public.create_hypertable('cmc.price_and_cap', 'ts', chunk_time_interval => INTERVAL '1 year');


create table cmc.concentration
(
    currency_id                         bigint                  not null,
    whales                              double precision        not null,
    investors                           double precision        not null,
    retail                              double precision        not null,
    d                                   date                    not null,
    CONSTRAINT concentration__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index concentration__currency_id__ts__pk ON cmc.concentration (currency_id, d) include (whales, investors, retail);

select public.create_hypertable('cmc.concentration', 'd', chunk_time_interval => INTERVAL '1 year');



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop schema cmc cascade;
-- +goose StatementEnd
