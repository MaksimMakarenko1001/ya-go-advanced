DROP FUNCTION metric.list_metrics();

DROP FUNCTION metric.counters_list_by_metric_names(text[]);
DROP FUNCTION metric.gauges_list_by_metric_names(text[]);

DROP FUNCTION metric.counters_upsert(json);
DROP FUNCTION metric.gauges_upsert(json);