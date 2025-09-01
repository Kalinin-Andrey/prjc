-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table cmc.token_address
(
    currency_id                 bigint                  not null,
    blockchain                  text                    not null,
    address                     text                    not null
);
create unique index token_address__currency_id__blockchain__address__ux ON cmc.token_address (currency_id, blockchain, address);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop table cmc.token_address;
-- +goose StatementEnd
