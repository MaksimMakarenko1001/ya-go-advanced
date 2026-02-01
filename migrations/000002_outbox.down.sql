DROP FUNCTION outbox.outbox_add_new(json, text);
DROP FUNCTION outbox.outbox_get_next(text, text, integer);
DROP FUNCTION outbox.outbox_set_failed(text[], text);
DROP FUNCTION outbox.outbox_set_completed(text[], text);

DROP FUNCTION metric.metrics_upsert(json, json, json)

DROP SEQUENCE IF EXISTS outbox.outbox_id_seq;

DROP TABLE IF EXISTS outbox.outbox;

DROP SCHEMA IF EXISTS outbox;