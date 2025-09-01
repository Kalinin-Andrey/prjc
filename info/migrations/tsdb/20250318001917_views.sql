-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE VIEW cmc.whales_prc_and_min_max_price AS
with d as (
    select currency_id, max(d) as last_d
    from cmc.concentration
    group by currency_id),
     whales as (
         select currency_id, d, whales * 100 / (whales + investors + retail) as whales_prc
         from cmc.concentration),
     ts as (
         select currency_id, max(ts) as last_ts
         from cmc.price_and_cap
         group by currency_id),
     min_max as (
         select currency_id, min(price) as min_price, max(price) as max_price
         from cmc.price_and_cap
         group by currency_id),
     price as (
         select currency_id, price, ts
         from cmc.price_and_cap
     )
select c.id, c.symbol, round(cast(whales.whales_prc AS numeric), 2) as whales_prc, price.price, min_max.min_price, min_max.max_price,
       round(cast(min_max.max_price/price.price AS numeric)) as to_ath, round(cast(price.price/min_max.min_price AS numeric)) as from_atl
from cmc.currency c
         inner join ts on c.id = ts.currency_id
         inner join price on ts.currency_id = price.currency_id and ts.last_ts = price.ts
         left join d on c.id = d.currency_id
         left join whales on d.currency_id = whales.currency_id and d.last_d = whales.d
         inner join min_max on c.id = min_max.currency_id
where c.is_for_observing = true;

CREATE VIEW cmc.cw1 AS
with m2 as (
    select currency_id, max(d) as d
    from cmc.concentration
    where d <= (now() - interval '2 month')::date
        group by currency_id),
        m1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 month')::date
        group by currency_id),
        w2 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 week')::date
        group by currency_id),
        w1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 week')::date
        group by currency_id),
        n as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 day')::date		-- из-за глюков в данных
        group by currency_id)
select c.id, c.symbol, round(cast(((cn.whales - cw1.whales) * 100)/cw1.whales * 1000 AS numeric)) as bonus,
       round(cast(((cn.whales - cw1.whales) * 100)/cw1.whales AS numeric), 4) as week1, round(cast(((cn.whales - cw2.whales) * 100)/cw2.whales AS numeric), 4) as week2,
       round(cast(((cn.whales - cm1.whales) * 100)/cm1.whales AS numeric), 4) as month1, round(cast(((cn.whales - cm2.whales) * 100)/cm2.whales AS numeric), 4) as month2
from cmc.currency c
         inner join m2 on c.id = m2.currency_id
         inner join m1 on c.id = m1.currency_id
         inner join w2 on c.id = w2.currency_id
         inner join w1 on c.id = w1.currency_id
         inner join n on c.id = n.currency_id
         inner join cmc.concentration cm2 on m2.currency_id = cm2.currency_id and m2.d = cm2.d
         inner join cmc.concentration cm1 on m1.currency_id = cm1.currency_id and m1.d = cm1.d
         inner join cmc.concentration cw2 on w2.currency_id = cw2.currency_id and w2.d = cw2.d
         inner join cmc.concentration cw1 on w1.currency_id = cw1.currency_id and w1.d = cw1.d
         inner join cmc.concentration cn on n.currency_id = cn.currency_id and n.d = cn.d
where c.is_for_observing = true
order by week1, week2, month1, month2;

CREATE VIEW cmc.cw2 AS
with m2 as (
    select currency_id, max(d) as d
    from cmc.concentration
    where d <= (now() - interval '2 month')::date
        group by currency_id),
        m1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 month')::date
        group by currency_id),
        w2 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 week')::date
        group by currency_id),
        w1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 week')::date
        group by currency_id),
        n as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 day')::date		-- из-за глюков в данных
        group by currency_id)
select c.id, c.symbol, round(cast(((cn.whales - cw2.whales) * 100)/cw2.whales * 1000 AS numeric)) as bonus,
       round(cast(((cn.whales - cw1.whales) * 100)/cw1.whales AS numeric), 4) as week1, round(cast(((cn.whales - cw2.whales) * 100)/cw2.whales AS numeric), 4) as week2,
       round(cast(((cn.whales - cm1.whales) * 100)/cm1.whales AS numeric), 4) as month1, round(cast(((cn.whales - cm2.whales) * 100)/cm2.whales AS numeric), 4) as month2
from cmc.currency c
         inner join m2 on c.id = m2.currency_id
         inner join m1 on c.id = m1.currency_id
         inner join w2 on c.id = w2.currency_id
         inner join w1 on c.id = w1.currency_id
         inner join n on c.id = n.currency_id
         inner join cmc.concentration cm2 on m2.currency_id = cm2.currency_id and m2.d = cm2.d
         inner join cmc.concentration cm1 on m1.currency_id = cm1.currency_id and m1.d = cm1.d
         inner join cmc.concentration cw2 on w2.currency_id = cw2.currency_id and w2.d = cw2.d
         inner join cmc.concentration cw1 on w1.currency_id = cw1.currency_id and w1.d = cw1.d
         inner join cmc.concentration cn on n.currency_id = cn.currency_id and n.d = cn.d
where c.is_for_observing = true
order by week2, week1, month1, month2;

CREATE VIEW cmc.cm1 AS
with m2 as (
    select currency_id, max(d) as d
    from cmc.concentration
    where d <= (now() - interval '2 month')::date
        group by currency_id),
        m1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 month')::date
        group by currency_id),
        w2 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 week')::date
        group by currency_id),
        w1 as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '1 week')::date
        group by currency_id),
        n as (
        select currency_id, max(d) as d
        from cmc.concentration
        where d <= (now() - interval '2 day')::date		-- из-за глюков в данных
        group by currency_id)
select c.id, c.symbol, round(cast(((cn.whales - cm1.whales) * 100)/cm1.whales * 1000 AS numeric), 2) as bonus,
       round(cast(((cn.whales - cw1.whales) * 100)/cw1.whales AS numeric), 4) as week1, round(cast(((cn.whales - cw2.whales) * 100)/cw2.whales AS numeric), 4) as week2,
       round(cast(((cn.whales - cm1.whales) * 100)/cm1.whales AS numeric), 4) as month1, round(cast(((cn.whales - cm2.whales) * 100)/cm2.whales AS numeric), 4) as month2
from cmc.currency c
         inner join m2 on c.id = m2.currency_id
         inner join m1 on c.id = m1.currency_id
         inner join w2 on c.id = w2.currency_id
         inner join w1 on c.id = w1.currency_id
         inner join n on c.id = n.currency_id
         inner join cmc.concentration cm2 on m2.currency_id = cm2.currency_id and m2.d = cm2.d
         inner join cmc.concentration cm1 on m1.currency_id = cm1.currency_id and m1.d = cm1.d
         inner join cmc.concentration cw2 on w2.currency_id = cw2.currency_id and w2.d = cw2.d
         inner join cmc.concentration cw1 on w1.currency_id = cw1.currency_id and w1.d = cw1.d
         inner join cmc.concentration cn on n.currency_id = cn.currency_id and n.d = cn.d
where c.is_for_observing = true
order by month1, week1, week2, month2;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP VIEW cmc.whales_prc_and_min_max_price;
DROP VIEW cmc.cw1;
DROP VIEW cmc.cw2;
DROP VIEW cmc.cm1;
-- +goose StatementEnd
