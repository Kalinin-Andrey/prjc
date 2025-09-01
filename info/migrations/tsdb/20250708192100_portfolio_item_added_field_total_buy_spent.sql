-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

alter table cmc.portfolio_item add column total_buy_spent double precision not null default 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

alter table cmc.portfolio_item drop column total_buy_spent;
-- +goose StatementEnd
