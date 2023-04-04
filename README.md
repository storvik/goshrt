# goshrt - Go URL shortener

[![build](https://github.com/storvik/goshrt/actions/workflows/build.yml/badge.svg)](https://github.com/storvik/goshrt/actions/workflows/build.yml)
[![go test](https://github.com/storvik/goshrt/actions/workflows/gotest.yml/badge.svg)](https://github.com/storvik/goshrt/actions/workflows/gotest.yml)

> Work in progress!

This is my attempt at creating a self hosted URL shortener written in Go.
The goal is to support multiple domains, cache and a simple API for creating new entries.

## Keywords
- Postgresql
- API first

## Development

### Postgres

When doing local development postgres has to be running.
While there are several ways to achieve this, VM / docker / podman, I myself use Nix.
Spinning up a development database is simple in Nix develop shell:

``` shell
  $ nix develop
  $ pgnix-init     # initiate database and start it
  $ pgnix-status   # check if database is running
  $ pgnix-stop     # stop postgresql database
```

## Todo
- [ ] Add rest api for adding and deleting shrts
- [ ] Client, goshrtc
- [ ] Add support for multiple domains
- [ ] Add support for random generated slugs
- [ ] Add support for user specified slugs
- [ ] Redis cache in front of postgresql
- [ ] Metrics on visited urls
- [ ] Add instructions for
  - [ ] Setup
  - [ ] Nginx proxy
- [ ] Add package to flake.nix
