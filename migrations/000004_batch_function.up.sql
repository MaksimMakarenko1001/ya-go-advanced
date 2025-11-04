CREATE OR REPLACE FUNCTION metric.upsert_metrics(_counter_items json, _gauge_items json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with cte(metric_name) as (
        select * from json_array_elements(metric.counters_upsert(_counter_items))
        union all
        select * from json_array_elements(metric.gauges_upsert(_gauge_items))
    )
   select json_agg(cte.metric_name) from cte
	    into _res;

    return coalesce(_res, '[]'::json);
end;
$function$
;