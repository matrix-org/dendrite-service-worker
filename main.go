// -*- coding: utf-8 -*-
// Copyright 2020 The Matrix.org Foundation C.I.C.
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

package main

// This is a mock Dendrite used for testing the service worker
// without actually having to go through the pain of getting Dendrite
// to use go-sqlite-js and go-http-js-libp2p.

// It can also be used as a guide to getting Dendrite doing the right thing.

import "database/sql"
import "log"
import "fmt"
import "net/http"
import "io/ioutil"


import "github.com/matrix-org/go-http-js-libp2p/go_http_js_libp2p"
import _ "github.com/matrix-org/go-sqlite3-js/sqlite3_js"

var c chan struct{}

func init() {
    c = make(chan struct{})
}

func main() {
	// we listen for federation traffic via go-http-libp2p-js
	// and for now respond to /ping with the contents of the
	// ping table (incrementing the row every time we get a ping).

	db := initDb()
	node := go_http_js_libp2p.NewP2pLocalNode("org.matrix.p2p.experiment")
	federationServer(node, db)
	clientServer(db)

	// due to https://github.com/golang/go/issues/27495 we can't override the DialContext
	// instead we have to provide a whole custom transport.
	client := &http.Client{
		Transport: go_http_js_libp2p.NewP2pTransport(node),
	}

	// try to ping every peer that we discover which supports this service
	node.RegisterFoundProvider(func(pi *go_http_js_libp2p.PeerInfo) {
		go func() {
			log.Printf("Trying to GET libp2p-http://%s/ping", pi.Id)
			resp, err := client.Get(fmt.Sprintf("libp2p-http://%s/_matrix/client/ping", pi.Id))
			if err != nil {
				log.Fatal("Can't make request")
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				bodyString := string(bodyBytes)
				log.Printf("Received body: %s", bodyString)
			}
		}()
	})

    <-c
}

func initDb() (*sql.DB) {
	var db *sql.DB
    var err error
    if db, err = sql.Open("sqlite3_js", "test.db"); err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("create table ping(id int)")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("insert into ping values(1)")
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func federationServer(node *go_http_js_libp2p.P2pLocalNode, db *sql.DB) {
	http.HandleFunc("/_matrix/federation/ping", func(w http.ResponseWriter, r *http.Request) {
		var err error
	    _, err = db.Exec("update ping set id=id+1")
	    if err != nil {
	        log.Fatal(err)
	    }

	    rows, err := db.Query("select id from ping")
	    if err != nil {
	        log.Fatal(err)
	    }
	    var id int
	    defer rows.Close()
        err = rows.Scan(&id)
	    if err != nil {
	        log.Fatal(err)
	    }
	    fmt.Fprintf(w, "pong: %d", id)
	})

	log.Println("starting SS API server")

	listener := go_http_js_libp2p.NewP2pListener(node)
	s := &http.Server{}
	go s.Serve(listener)
}

func clientServer(db *sql.DB) {
	http.HandleFunc("/_matrix/client/ping", func(w http.ResponseWriter, r *http.Request) {
		var err error
	    _, err = db.Exec("update ping set id=id+1")
	    if err != nil {
	        log.Fatal(err)
	    }

	    rows, err := db.Query("select id from ping")
	    if err != nil {
	        log.Fatal(err)
	    }
	    if rows.Next() {
		    var id int
		    defer rows.Close()
	        err = rows.Scan(&id)
		    if err != nil {
		        log.Fatal(err)
		    }
		    fmt.Fprintf(w, "pong: %d", id)
		}
	})

	log.Println("starting CS API server")

	listener := go_http_js_libp2p.NewFetchListener()
	s := &http.Server{}
	go s.Serve(listener)
}
