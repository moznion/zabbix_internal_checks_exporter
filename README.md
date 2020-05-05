# zabbix_internal_checks_exporter

A Prometheus exporter for [Zabbix internal checks](https://www.zabbix.com/documentation/current/manual/config/items/itemtypes/internal) metrics.

## Usage

```
$ zabbix-internal-checks-exporter --help
Usage of zabbix-internal-checks-exporter:
  -interval-sec uint
        [optional] an interval seconds of collecting the metrics (default 30)
  -port int
        [mandatory] a port number of exporter listens (default -1)
  -zabbix-password string
        [mandatory] a Zabbix server password for authentication to use API
  -zabbix-url string
        [mandatory] a Zabbix server URL to collect the metrics
  -zabbix-user string
        [mandatory] a Zabbix server user name for authentication to use API
```

And this exporter returns the metrics when that received a HTTP request to `GET /metrics`.

## Supported Zabbix version

- 4.4

In other versions, it hasn't confirmed what it works certainly. If you know the other Zabbix version that works well with this exporter, it would be awesome if you could report that through an issue.

## How does it work

This exporter retrieves "Zabbix internal checks" through JSON-RPC API calling: [item.get](https://www.zabbix.com/documentation/current/manual/api/reference/item/get) for `zabbix[*]`.

## Note

This exporter uses a key of a Zabbix metric as a Prometheus metric name, but it sanitizes some metric name because of the limitation of Prometheus's naming rule.

|original|sanitized|
|--------|---------|
| '['    | '\_\_'  |
| ','    | ':'     |
| ']'    | ''      |
| ' '    | '\_'    |
| '-'    | '\_'    |

