## `init-exporter-converter` [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/funbox/init-exporter-converter)](https://goreportcard.com/report/github.com/funbox/init-exporter-converter) [![License](https://gh.kaos.st/mit.svg)](LICENSE)

Utility for converting [`init-exporter`](https://github.com/funbox/init-exporter) procfiles from v1 to v2 format.

### Installation

#### From sources

To build the `init-exporter-converter` from scratch, make sure you have a working Go 1.16+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/funbox/init-exporter-converter
```

If you want to update `init-exporter-converter` to latest stable release, do:

```
go get -u github.com/funbox/init-exporter-converter
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install init-exporter-converter
```

### Usage

```
Usage: init-exporter-converter {options} procfile…

Options

  --config, -c file    Path to init-exporter config
  --in-place, -i       Edit procfile in place
  --no-colors, -nc     Disable colors in output
  --help, -h           Show this help message
  --version, -v        Show version

Examples

  init-exporter-converter -i config/Procfile.production
  Convert Procfile.production to version 2 in-place

  init-exporter-converter -c /etc/init-exporter.conf Procfile.*
  Convert all procfiles to version 2 with defaults from init-exporter config and print result to console

```

### Build status

| Branch | Status |
|--------|--------|
| Stable | [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) |
| Unstable | [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) |

### License

[MIT](LICENSE)

[![Sponsored by FunBox](https://funbox.ru/badges/sponsored_by_funbox_grayscale.svg)](https://funbox.ru)
