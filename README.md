### dendrite-service-worker

This is an experimental JS library which wraps a WASM-compiled
[Dendrite](https://github.com/matrix-org/dendrite) (golang Matrix homeserver)
into a web browser service worker.

It's intended to be imported by a Matrix webclient such as Riot/Web in order
to embed a homeserver in the browser for use with experimental P2P Matrix.

It relies on a branch of Dendrite which uses the
[go-sqlite3-js](https://github.com/matrix-org/go-sqlite3-js) driver for
storage (thus delegating SQL to sql.js, running in-browser), and the
[go-http-js-libp2p](https://github.com/matrix-org/go-http-js-libp2p) transport
for HTTP (thus tunnelling HTTP over libp2p via js-libp2p, running
in-browser).

In order to bootstrap, we start off with a mock dendrite in main.go, which is built via:

`GOOS=js GOARCH=wasm go build -o main.wasm`

It's unclear at this stage whether the library should provide the service worker
script (sw.js) or if the calling application should.


To build dendrite, go to the dendrite repo and: (may need to be on `kegan/wasm`):
```
$ GOOS=js GOARCH=wasm go build -o main.wasm ./cmd/dendritejs
$ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

If you update go-sqlite-js you need to:
```
$ yarn add "https://github.com/matrix-org/go-sqlite3-js#master"
```

And to pull in the latest Go code in dendrite:
```
$ go get github.com/matrix-org/go-sqlite3-js@master
```
