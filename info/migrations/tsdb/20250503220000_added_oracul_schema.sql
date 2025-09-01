-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create schema oracul;

create table oracul.analytics
(
    currency_id                 bigint                  not null,
    whales_concentration        double precision        not null,
    worm_index                  double precision        not null,
    growth_fuel                 double precision        not null,
    ts                          timestamp               not null,
    CONSTRAINT analytics__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index analytics__currency_id__ts__pk ON oracul.analytics (currency_id, ts);
select public.create_hypertable('oracul.analytics', 'ts', chunk_time_interval => INTERVAL '1 year');


create table oracul.speedometers
(
    currency_id                 bigint                  not null,
    whales_buy_rate             double precision        not null,
    whales_sell_rate            double precision        not null,
    whales_volume               double precision        not null,
    investors_buy_rate          double precision        not null,
    investors_sell_rate         double precision        not null,
    investors_volume            double precision        not null,
    retailers_buy_rate          double precision        not null,
    retailers_sell_rate         double precision        not null,
    retailers_volume            double precision        not null,
    ts                          timestamp               not null,
    CONSTRAINT speedometers__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index speedometers__currency_id__ts__pk ON oracul.speedometers (currency_id, ts);
select public.create_hypertable('oracul.speedometers', 'ts', chunk_time_interval => INTERVAL '1 year');


create table oracul.holder_stats
(
    currency_id                 bigint                  not null,
    whales_volume               double precision        not null,
    whales_total_holders        bigint                  not null,
    investors_volume            double precision        not null,
    investors_total_holders     bigint                  not null,
    retailers_volume            double precision        not null,
    retailers_total_holders     bigint                  not null,
    ts                          timestamp               not null,
    CONSTRAINT holder_stats__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index holder_stats__currency_id__ts__pk ON oracul.holder_stats (currency_id, ts);
select public.create_hypertable('oracul.holder_stats', 'ts', chunk_time_interval => INTERVAL '1 year');


create table oracul.daily_balance_stats
(
    currency_id                 bigint                  not null,
    whales_balance              double precision        not null,
    whales_total_holders        bigint                  not null,
    investors_balance           double precision        not null,
    investors_total_holders     bigint                  not null,
    retailers_balance           double precision        not null,
    retailers_total_holders     bigint                  not null,
    d                           date                    not null,
    CONSTRAINT daily_balance_stats__currency_id__fk FOREIGN KEY (currency_id) REFERENCES cmc.currency(id)
);

create unique index daily_balance_stats__currency_id__d__pk ON oracul.daily_balance_stats (currency_id, d);
select public.create_hypertable('oracul.daily_balance_stats', 'd', chunk_time_interval => INTERVAL '1 year');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop schema oracul;
-- +goose StatementEnd
