# Healthcheck

[![Automate](https://github.com/viral32111/healthcheck/actions/workflows/automate.yml/badge.svg?branch=main)](https://github.com/viral32111/healthcheck/actions/workflows/automate.yml)
[![CodeQL](https://github.com/viral32111/healthcheck/actions/workflows/codeql.yml/badge.svg)](https://github.com/viral32111/healthcheck/actions/workflows/codeql.yml)
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/viral32111/healthcheck?label=Latest)
![GitHub repository size](https://img.shields.io/github/repo-size/viral32111/healthcheck?label=Size)
![GitHub release downloads](https://img.shields.io/github/downloads/viral32111/healthcheck/total?label=Downloads)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/viral32111/healthcheck?label=Commits)

This is a utility for Docker containers to periodically check their health by checking the status code of HTTP requests.

## 📥 Usage

Download the [latest release](https://github.com/viral32111/healthcheck/releases/latest) for your platform. There are builds available for Linux and Windows, on 32-bit and 64-bit architectures of x86 and ARM. There are extra Linux builds to accommodate glibc and musl libraries. This should cover the majority of Docker images.

The utility expects a target HTTP URL as the only argument (multiple arguments will be joined together to form the URL). It will exit with a status code of `0` if the check was successful (i.e. the HTTP response matched what was expected), or `1` if there was any error (invalid flags, destination unreachable, mismatching HTTP status code, etc.).

There are optional flags for fine-tuning functionality:

* `--expect <number>`: The HTTP response status code to consider successful. Defaults to `200`.
* `--method <string>`: The HTTP request method. Defaults to `GET`.
* `--proxy <ip:port>`: The IP address and port number of a proxy server. Useful for checking .onion sites.

These flags can be prefixed with either a single (`-`) or double (`--`) hyphen.

Use the `--help` (`-h`) flag for more information.

### 🐳 Docker

Use either with the [Dockerfile `HEALTHCHECK` instruction](https://docs.docker.com/engine/reference/builder/#healthcheck) or the [docker run `--health-*`](https://docs.docker.com/engine/reference/run/#healthcheck) flags.

For example, `HEALTHCHECK... CMD healthcheck http://127.0.0.1`, or `docker run..... --health-cmd healthcheck http://127.0.0.1/ .....image:tag`.

### 🖼️ Examples

Checking if the `/metrics` endpoint at `localhost` on port `5000` will respond with a `200` status code when sending a `GET` request:

```
$ healthcheck http://localhost:5000/metrics
SUCCESS, 200 OK
```

Checking if the `/betrics` endpoint at `localhost` on port `5000` will respond with a `204` status code when sending a `GET` request:

```
$ healthcheck -expect 204 http://localhost:5000/betrics
FAILURE, 404 Not Found
```

Checking if the onion site at `hiddenservice.onion` on port `80` will respond with a `100` status code when sending a `GET` request, through the local Tor SOCKS5 proxy:

```
$ healthcheck --expect 100 --proxy 127.0.0.1:9050 http://hiddenservice.onion
SUCCESS, 100 Continue
```

Checking if the `/health` endpoint at `example.com` on port `443` will respond with a `200` status code when sending a `POST` request:

```
$ healthcheck --method POST https://example.com/health
SUCCESS, 200 OK
```

## ⚖️ License

Copyright (C) 2022-2023 [viral32111](https://viral32111.com).

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see https://www.gnu.org/licenses.
