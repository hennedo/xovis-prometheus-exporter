# Prometheus exporter for Xovis PC2S sensors

Leave this software running and configure your PC2S sensor to push line count data to http://addressofhost:8080/xovis, then let your prometheus scrape http://addressofhost:8080/metrics and you'll have xovis_entries, xovis_exits and xovis_sum as data points