// SPDX-License-Identifier: BlueOak-1.0.0
// Copyright Russell Hernandez Ruiz <qrpnxz@hyperlife.xyz>

// Command wkdserver serves a WKD as specified by
// https://tools.ietf.org/html/draft-koch-openpgp-webkey-service-11
//
// The first argument is the address on which the server will listen
// for connections. It is optional.
//
// Keys are taken from files in pgpKeyDir. For example, the keys for
// the address mister@example.org are in the file named
// r3ptdiy83btqwgjkooeprx3udzwcr34a in pgpKeyDir.
package main

import (
	"crypto/sha1"
	"encoding/base32"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	domain      = "example.org"
	pgpKeyDir   = "/var/lib/wkdserver"
	certFile    = "/etc/letsencrypt/live/example.org/fullchain.pem"
	certKeyFile = "/etc/letsencrypt/live/example.org/privkey.pem"
)

var users = []string{
	"john.doe",
	"jane.doe",
}

var zBase32 = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

func sum(p []byte) []byte {
	hash := sha1.Sum(p)
	dst := make([]byte, zBase32.EncodedLen(len(hash)))
	zBase32.Encode(dst, hash[:])
	return dst
}

func keyRequestHandler(w http.ResponseWriter, keyFile string) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	f, err := os.Open(keyFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	const wellknown = "/.well-known/openpgpkey/" + domain
	http.HandleFunc(
		wellknown+"/policy",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
		},
	)
	for _, user := range users {
		hash := string(sum([]byte(user)))
		keyFile := pgpKeyDir + "/" + hash
		log.Printf(
			"Serving the keys of %s@%s from %s\n",
			user, domain, keyFile,
		)
		http.HandleFunc(
			wellknown+"/hu/"+hash,
			func(w http.ResponseWriter, r *http.Request) {
				keyRequestHandler(w, keyFile)
			},
		)
	}
	var addr string
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	log.Fatal(http.ListenAndServeTLS(addr, certFile, certKeyFile, nil))
}
