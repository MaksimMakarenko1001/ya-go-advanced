DROP FUNCTION outbox.outbox_add_new(json, text);
DROP FUNCTION outbox.outbox_get_next(text, text, integer);
DROP FUNCTION outbox.outbox_commit(text[],text[], text);

DROP FUNCTION metric.metrics_upsert(json, json, json);
CREATE OR REPLACE FUNCTION metric.metrics_upsert(_counter_items json, _gauge_items json)
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

DROP SEQUENCE IF EXISTS outbox.outbox_id_seq;

DROP TABLE IF EXISTS outbox.outbox;

DROP SCHEMA IF EXISTS outbox;