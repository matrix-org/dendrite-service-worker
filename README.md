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
