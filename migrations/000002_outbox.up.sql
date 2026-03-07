CREATE SCHEMA IF NOT EXISTS outbox;

-- Create custom sequence for outbox table
CREATE SEQUENCE IF NOT EXISTS outbox.outbox_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS outbox.outbox(
    id TEXT DEFAULT nextval('outbox.outbox_id_seq'::regclass)::TEXT NOT NULL,
    destination TEXT NOT NULL,
    segment TEXT NOT NULL,
    payload JSON NOT NULL,
    lock_until TIMESTAMPTZ
);

ALTER TABLE outbox.outbox
    ADD CONSTRAINT outbox_pkey PRIMARY KEY (id);

CREATE OR REPLACE FUNCTION outbox.outbox_add_new(_items json, _segment text) RETURNS void
    LANGUAGE plpgsql
    AS $$
begin
    perform pg_advisory_xact_lock(hashtext('outbox_'||_segment));

    with cte as (
        select * from json_populate_recordset(null::outbox.outbox, _items)
    )
    insert into outbox.outbox (destination, segment, payload)
        select src.destination, src.segment, src.payload
            from cte as src
    ;
end;
$$;

CREATE OR REPLACE FUNCTION outbox.outbox_get_next(_destination text, _segment text, _limit integer = 100) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare _res json;
begin
    perform pg_advisory_xact_lock(hashtext('outbox_'||_segment));

    with
        cte as (
            select * 
                from outbox.outbox as src
                    where src.destination = _destination
                        and src.segment = _segment
                        and (src.lock_until is null or src.lock_until < now())
                limit _limit
        ),
        upd_cte as (
            update outbox.outbox as upd set
                lock_until = now() + '30sec'::interval
            from cte as src
                where upd.id = src.id
        )
    select json_agg(src.*)
        into _res
        from cte as src
    ;

    return coalesce(_res, '[]'::json);
end;
$$;

CREATE OR REPLACE FUNCTION outbox.outbox_commit(_ok_ids text[], _failed_ids text[], _segment text) RETURNS void
    LANGUAGE plpgsql
    AS $$
begin
    perform pg_advisory_xact_lock(hashtext('outbox_'||_segment));

    delete from outbox.outbox as del
        where del.id = any(_ok_ids)
    ;

    update outbox.outbox as upd set
        lock_until = now()
        where upd.id = any(_failed_ids)
    ;
end;
$$;

DROP FUNCTION metric.metrics_upsert(json, json);
CREATE OR REPLACE FUNCTION metric.metrics_upsert(_counter_items json, _gauge_items json, _outbox_items json = NULL::json, _outbox_segment text = ''::text)
 RETURNS json
 LANGUAGE plpgsql
AS $$
declare
    _res json;
begin
    with cte(metric_name) as (
        select * from json_array_elements(metric.counters_upsert(_counter_items))
        union all
        select * from json_array_elements(metric.gauges_upsert(_gauge_items))
    )
    select json_agg(cte.metric_name)
	    into _res
        from cte
    ;

    perform outbox.outbox_add_new(_outbox_items, _outbox_segment);

    return coalesce(_res, '[]'::json);
end;
$$
;