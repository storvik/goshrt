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

When doing local development/testing postgres has to be running.
While there are several ways to achieve this, VM / docker / podman, I myself use Nix.
Spinning up a development database is very simple in Nix shell.
After installing Nix, devshell is entered through the command `nix develop`.

``` shell
$ pgnix-init     # initiate database and start it
$ pgnix-start    # start database
$ pgnix-status   # check if database is running
$ pgnix-restart  # restart database
$ pgnix-stop     # stop postgresql database
$ pgnix-purge    # stop database and delete it
$ pgnix-pgcli    # start pgcli and connect to database
$ pgnix-psql     # start psql and connect to database
```

### Unit testing

Nix shell and `pgnix-` wrappers makes running unit test in a clean environment very simple.
Inside `nix develop` the following oneliner runs all unit tests:

``` shell
$ pgnix-purge && pgnix-init && go clean -testcache && go test -v ./...
```

> `go clean -testcache` ensures that all tests are run.
> Without it tests will be cached and for instance database migragions will not be run.

## Todo
- [x] Add rest api for adding and getting shrts
- [x] Client, goshrtc
- [x] Add support for multiple domains
- [x] Add support for random generated slugs
- [x] Add support for user specified slugs
- [x] Authentication
- [ ] Add delete shrt should use id
- [ ] Add getting list of shrts
- [ ] Redis cache in front of postgresql
- [ ] Metrics on visited urls
- [ ] Add instructions for
  - [ ] Setup
  - [ ] Nginx proxy
- [ ] Add package to flake.nix
