CREATE SCHEMA IF NOT EXISTS metric;

CREATE TABLE IF NOT EXISTS metric.counters (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    metric_type TEXT NOT NULL DEFAULT 'counter',
    metric_name TEXT UNIQUE NOT NULL,
    metric_value BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS metric.gauges (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    metric_type TEXT NOT NULL DEFAULT 'gauge',
    metric_name TEXT UNIQUE NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE metric.logs (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	source TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
	audited_at TIMESTAMPTZ NULL
);

CREATE OR REPLACE FUNCTION metric.metrics_list()
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with 
        counter_data as (
            select c.* from metric.counters as c
        ),
        gauge_data as (
            select g.* from metric.gauges as g
        )
    select
		json_build_object(
			'counters', (select json_agg(r.*) from counter_data as r),
			'gauges', (select json_agg(r.*) from gauge_data as r)
		)
	    into _res;

    return _res;
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.counters_list_by_metric_names(_metric_names text[])
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with 
        cte as (
            select c.* from metric.counters as c
                where c.metric_name = any(_metric_names)
        )
    select json_agg(cte.*) from cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.gauges_list_by_metric_names(_metric_names text[])
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with 
        cte as (
            select g.* from metric.gauges as g
                where g.metric_name = any(_metric_names)
        )
    select json_agg(cte.*) from cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.counters_upsert(_ip_address text, _items json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with 
        cte as (
            select * from json_populate_recordset(null::metric.counters, _items)
        ),
        ins_cte as (
            insert into metric.counters as c (metric_type, metric_name, metric_value,
                    created_at, updated_at)
            select cte.metric_type, cte.metric_name, cte.metric_value,
                    cte.created_at, cte.updated_at
                from cte
            on conflict (metric_name) do update
                set metric_value = c.metric_value + excluded.metric_value,
                    updated_at = excluded.updated_at
            returning c.metric_name
        ),
        logs_cte as (
            insert into metric.logs (source, ip_address, metric_name, created_at)
                select 'metric.counters_upsert', _ip_address, src.metric_name, now()
                    from ins_cte as src
        )
    select json_agg(ins_cte.metric_name) from ins_cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.gauges_upsert(_ip_address text, _items json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with 
        cte as (
            select * from json_populate_recordset(null::metric.gauges, _items)
        ),
        ins_cte as (
            insert into metric.gauges as g (metric_type, metric_name, metric_value,
                    created_at, updated_at)
            select src.metric_type, src.metric_name, src.metric_value,
                    src.created_at, src.updated_at
                from cte as src
            on conflict (metric_name) do update
                set metric_value = excluded.metric_value,
                    updated_at = excluded.updated_at
            returning g.metric_name
        ),
        logs_cte as (
            insert into metric.logs (source, ip_address, metric_name, created_at)
                select 'metric.gauges_upsert', _ip_address, src.metric_name, now()
                    from ins_cte as src
        )
    select json_agg(ins_cte.metric_name) from ins_cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.metrics_upsert(_ip_address text, _counter_items json, _gauge_items json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with cte(metric_name) as (
        select * from json_array_elements(metric.counters_upsert(_ip_address, _counter_items))
        union all
        select * from json_array_elements(metric.gauges_upsert(_ip_address, _gauge_items))
    )
    select json_agg(cte.metric_name) from cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;