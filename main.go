package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func siegfried(ctx context.Context, sf, home string) {
	cmd := exec.CommandContext(ctx, sf, "-home", home, "-fpr")

	var err error

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	sout, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(sout)
	}

	serr, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Fatal(serr)
	}

	if err = cmd.Wait(); err != nil {
		fmt.Println(string(sout), string(serr))
		log.Fatal(err)
	}
}

func handleErr(w http.ResponseWriter, status int, e error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("server error; got %v\n", e))
}

func decodePath(path string) ([]byte, error) {
	if len(path) < 2 {
		return nil, fmt.Errorf("path is empty")
	}
	data, err := base64.URLEncoding.DecodeString(path[1:])
	if err != nil {
		return nil, fmt.Errorf("Error base64 decoding file path, error message %v", err)
	}
	return data, nil
}

func httpd(ctx context.Context, addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			handleErr(w, http.StatusNotFound, fmt.Errorf("valid paths are /, /identify and /identify/*"))
			return
		}
		path, err := decodePath(r.URL.Path)
		if err != nil {
			handleErr(w, http.StatusBadRequest, err)
			return
		}
		puid, err := identify(ctx, path)
		if err != nil {
			handleErr(w, http.StatusBadRequest, err)
			return
		}
		fmt.Fprint(w, puid)
	})
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on ", addr)
	log.Fatal(server.ListenAndServe())
}

func identify(ctx context.Context, path []byte) (string, error) {
	const addr = "/tmp/siegfried"
	var d net.Dialer
	c, err := d.DialContext(ctx, "unix", addr)
	if err != nil {
		return "", err
	}
	defer c.Close()

	// Send query and read
	c.Write(path)
	buf := make([]byte, 1024)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}

	return string(buf[0:n]), nil
}

func main() {
	var (
		addr = flag.String("addr", ":8080", "tcp network address")
		sf   = flag.String("sf", "/sf", "sf binary")
		home = flag.String("home", "/siegfried", "siegfried data diretory")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go siegfried(ctx, *sf, *home)
	go httpd(ctx, *addr)

	// Subscribe to signals and wait
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan // Block until a signal is received

	cancel()
}
