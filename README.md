## `init-exporter-converter` [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/funbox/init-exporter-converter)](https://goreportcard.com/report/github.com/funbox/init-exporter-converter) [![License](https://gh.kaos.st/mit.svg)](LICENSE)

Utility for converting [`init-exporter`](https://github.com/funbox/init-exporter) procfiles from v1 to v2 format.

### Installation

#### From sources

To build the `init-exporter-converter` from scratch, make sure you have a working Go 1.23+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go install github.com/funbox/init-exporter-converter@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```
sudo yum install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install init-exporter-converter
```

### Usage

<img src=".github/images/usage.svg" />

### Build status

| Branch | Status |
|--------|--------|
| Stable | [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) |
| Unstable | [![CI](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/funbox/init-exporter-converter/actions/workflows/ci.yml) |

### License

[MIT](LICENSE)

[![Sponsored by FunBox](.github/images/sponsored.svg)](https://funbox.ru)
