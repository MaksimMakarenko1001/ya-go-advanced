DROP FUNCTION metric.metrics_list();
DROP FUNCTION metric.counters_list_by_metric_names(text[]);
DROP FUNCTION metric.gauges_list_by_metric_names(text[]);
DROP FUNCTION metric.counters_upsert(json);
DROP FUNCTION metric.gauges_upsert(json);
DROP FUNCTION metric.metrics_upsert(json, json);

DROP TABLE IF EXISTS metric.counters;
DROP TABLE IF EXISTS metric.gauges;

DROP SCHEMA IF EXISTS metric;