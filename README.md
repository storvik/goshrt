# goshrt - Go URL shortener

> Work in progress!

This is my attempt at creating a self hosted URL shortener written in Go.
The goal is to support multiple domains, cache and a simple API for creating new entries.

## Keywords
- Postgresql
- API first

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
