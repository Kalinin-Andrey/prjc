-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table cmc.portfolio_item
(
    portfolio_source_id         text                    not null,
    currency_id                 bigint                  not null,
    amount                      double precision        not null,
    current_price               double precision        not null,
    crypto_holdings             double precision        not null,
    holdings_percent            double precision        not null,
    buy_avg_price               double precision        not null,
    pl_percent_value            double precision        not null,
    pl_value                    double precision        not null,
    updated_at                  timestamp               not null
);
create unique index portfolio_item__portfolio_source_id__currency_id__ux ON cmc.portfolio_item (portfolio_source_id, currency_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop table cmc.portfolio_item;
-- +goose StatementEnd
