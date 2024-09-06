Bitcoind Image Build
====================

## Build

```
docker build --file Containerfile .
```

## Usage

```
docker run -it --rm --volume ./data:/data:rw --tmpfs /tmp --read-only --publish 127.0.0.1:8332:8332/tcp --publish 127.0.0.1:8333:8333/tcp ghcr.io/jmanero/bitcoind:latest
```
