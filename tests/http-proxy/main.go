// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

// #!/usr/bin/env bash
// case `uname -s` in
//     Linux*)     sslConfig=/etc/ssl/openssl.cnf;;
//     Darwin*)    sslConfig=/System/Library/OpenSSL/openssl.cnf;;
// esac
// openssl req \
//     -newkey rsa:2048 \
//     -x509 \
//     -nodes \
//     -keyout server.key \
//     -new \
//     -out server.pem \
//     -subj /CN=localhost \
//     -reqexts SAN \
//     -extensions SAN \
//     -config <(cat $sslConfig \
//         <(printf '[SAN]\nsubjectAltName=DNS:localhost')) \
//     -sha256 \
//     -days 3650

package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTunneling", r.Host)
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	byteRequest, err := httputil.DumpRequest(r, true)
	fmt.Println("handleTunneling byteRequest len", len(byteRequest))
	fmt.Println("handleTunnel Request\n", string(byteRequest))
	if err != nil {
		http.Error(w, "DumpRequest failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn, 1)
	go transfer(clientConn, destConn, 2)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, n int) {
	fmt.Println("transfer", n)
	defer destination.Close()
	defer source.Close()
	var buf bytes.Buffer
	w := io.MultiWriter(destination, &buf)
	io.Copy(w, source)
	fmt.Println("transfer", n, "writer buffer bytes len", len(buf.Bytes()))
	if n == 1 {
		fmt.Println("****************************************************************************************************************************************")
		fmt.Println("transfer", n, "after io.Copy", string(buf.Bytes()))
		fmt.Println("****************************************************************************************************************************************")
	} else {
		if len(buf.Bytes()) > 10000 {
			fmt.Println("transfer", n, "after io.Copy", string(buf.Bytes()[:10000]))
		} else {
			fmt.Println("transfer", n, "after io.Copy", string(buf.Bytes()))
		}
		fmt.Println("****************************************************************************************************************************************")
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleHTTP")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	fmt.Println("handleHTTP io.Copy")
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func onConnStateEvent(conn net.Conn, state http.ConnState) {
	fmt.Println("onChangeConnState LocalAddr", conn.LocalAddr().String(), "RemoteAddr", conn.RemoteAddr().String(), "state", state.String())
}

func main() {
	var pemPath string
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	var keyPath string
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	var proto string
	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
	flag.Parse()
	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}
	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		ConnState: onConnStateEvent,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	if proto == "http" {
		fmt.Println("ListenAndServe")
		log.Fatal(server.ListenAndServe())
		fmt.Println("after ListenAndServe")
	} else {
		fmt.Println("ListenAndServeTLS")
		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
		fmt.Println("after ListenAndServeTLS")
	}
}
