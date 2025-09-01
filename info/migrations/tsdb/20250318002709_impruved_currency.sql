-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE cmc.currency ADD COLUMN circulating_supply                  double precision    null;
ALTER TABLE cmc.currency ADD COLUMN self_reported_circulating_supply    double precision    null;
ALTER TABLE cmc.currency ADD COLUMN total_supply                        double precision    null;
ALTER TABLE cmc.currency ADD COLUMN max_supply                          double precision    null;
ALTER TABLE cmc.currency ADD COLUMN latest_price                        double precision    null;
ALTER TABLE cmc.currency ADD COLUMN cmc_rank                            integer             null;
ALTER TABLE cmc.currency ADD COLUMN date_added                          timestamp           null;
ALTER TABLE cmc.currency ADD COLUMN platform                            jsonb               null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE cmc.currency DROP COLUMN circulating_supply;
ALTER TABLE cmc.currency DROP COLUMN self_reported_circulating_supply;
ALTER TABLE cmc.currency DROP COLUMN total_supply;
ALTER TABLE cmc.currency DROP COLUMN max_supply;
ALTER TABLE cmc.currency DROP COLUMN latest_price;
ALTER TABLE cmc.currency DROP COLUMN cmc_rank;
ALTER TABLE cmc.currency DROP COLUMN date_added;
ALTER TABLE cmc.currency DROP COLUMN platform;
-- +goose StatementEnd
