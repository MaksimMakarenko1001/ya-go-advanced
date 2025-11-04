CREATE OR REPLACE FUNCTION metric.list_metrics()
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

CREATE OR REPLACE FUNCTION metric.counters_upsert(_items json)
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
        )
    select json_agg(ins_cte.metric_name) from ins_cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION metric.gauges_upsert(_items json)
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
        )
    select json_agg(ins_cte.metric_name) from ins_cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;
