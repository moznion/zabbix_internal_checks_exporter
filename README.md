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
  -version
        show version information
  -zabbix-password string
        [mandatory] a Zabbix server password for authentication to use API
  -zabbix-url string
        [mandatory] a Zabbix server URL to collect the metrics
  -zabbix-user string
        [mandatory] a Zabbix server user name for authentication to use API
```

And this exporter returns the metrics when that received a HTTP request to `GET /metrics`.

## Installation

Pre built binaries are on [Releases](https://github.com/moznion/zabbix_internal_checks_exporter/releases).

## Supported Zabbix version

- 4.4

In other versions, it hasn't confirmed what it works certainly. If you know the other Zabbix version that works well with this exporter, it would be awesome if you could report that through an issue.

## How does it work

This exporter retrieves "Zabbix internal checks" through JSON-RPC API calling: [item.get](https://www.zabbix.com/documentation/current/manual/api/reference/item/get) for `zabbix[*]`.

## Special Metric

### `zabbix_internal_checks_metrics_collected_at`

This value is set epoch milliseconds when the "Zabbix internal checks" metrics were collected. You can use this value to check whether it has collected the metrics periodically.

## Note

This exporter uses a key of a Zabbix metric as a Prometheus metric name, but it sanitizes some metric name because of the limitation of Prometheus's naming rule.

|original|sanitized        |
|--------|-----------------|
| `[`    | `__`            |
| `]`    | `__`            |
| `{`    | `__`            |
| `}`    | `__`            |
| `,`    | `:`             |
| ` `    | `_`             |
| `-`    | `_`             |

And it replaces other not allowed characters (i.e. `[^a-zA-Z0-9_:]`) to empty characters.

Please see also [Metric and label naming | Prometheus](https://prometheus.io/docs/practices/naming/).

## Metrics Example

```
# HELP zabbix_internal_checks_metrics_collected_at Epoch milliseconds at Zabbix internal checks metrics were collected
# TYPE zabbix_internal_checks_metrics_collected_at gauge
zabbix_internal_checks_metrics_collected_at 1.588776472169e+12
# HELP zabbix_internal_checks_zabbix__process:alerter:avg:busy__ Zabbix $4 $2 processes, in % (original key: zabbix[process,alerter,avg,busy])
# TYPE zabbix_internal_checks_zabbix__process:alerter:avg:busy__ gauge
zabbix_internal_checks_zabbix__process:alerter:avg:busy__ 0
# HELP zabbix_internal_checks_zabbix__process:configuration_syncer:avg:busy__ Zabbix $4 $2 processes, in % (original key: zabbix[process,configuration syncer,avg,busy])
# TYPE zabbix_internal_checks_zabbix__process:configuration_syncer:avg:busy__ gauge
zabbix_internal_checks_zabbix__process:configuration_syncer:avg:busy__ 8.5211
```

## License

```
The MIT License (MIT)
Copyright © 2020 moznion, https://moznion.net/ <moznion@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

