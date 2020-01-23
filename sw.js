// -*- coding: utf-8 -*-
// Copyright 2019, 2020 The Matrix.org Foundation C.I.C.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

self.importScripts("wasm_exec.js", "bundles/go_http_bridge.js", "bundles/sqlite_bridge.js")

self.addEventListener('install', function(event) {
    console.log("installing SW")
})

self.addEventListener('activate', function(event) {
    console.log("SW activated")

    event.waitUntil(
        sqlite_bridge.init().then(()=>{
            const go = new Go()
            WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
                go.run(result.instance)
            });
        })
    )
})

self.addEventListener('fetch', function(event) {
    console.log("intercepted " + event.request.url)
    if (event.request.url.match(/\/_matrix\/client/)) {
        if (global.fetchListener) {
            event.respondWith(global.fetchListener.onFetch(event))
        }
        else {
            console.log("no fetch listener present for " + event.request.url)
        }
    }
    else {
        return fetch(event.request)
    }
})
