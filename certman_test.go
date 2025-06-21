// Copyright 2017 Dyson Simmons. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package certman_test

import (
	"bytes"
	"crypto/tls"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/abh/certman"
)

func TestValidPair(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cm, err := certman.New("./testdata/server1.crt", "./testdata/server1.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	logGot := buf.String()

	if !strings.Contains(logGot, "certificate and key loaded") {
		t.Log("log output received:", logGot)
		t.Fatal("expected certificate loaded message not found")
	}
	if !strings.Contains(logGot, "watching for cert and key change") {
		t.Log("log output received:", logGot)
		t.Fatal("expected watching message not found")
	}
}

func TestInvalidPair(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cm, err := certman.New("./testdata/server1.crt", "./testdata/server2.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	logGot := buf.String()

	if !strings.Contains(logGot, "can't load cert or key file") {
		t.Log("log output received:", logGot)
		t.Fatal("expected error message not found")
	}
	if !strings.Contains(logGot, "watching for cert and key change") {
		t.Log("log output received:", logGot)
		t.Fatal("expected watching message not found")
	}
}

func TestCertificateNotFound(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cm, err := certman.New("./testdata/nothere.crt", "./testdata/server2.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		if !strings.HasPrefix(err.Error(), "can't watch cert file:") {
			t.Fatalf("unexpected watch error: %v", err)
		}
	}
}

func TestKeyNotFound(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cm, err := certman.New("./testdata/server1.crt", "./testdata/nothere.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		if !strings.HasPrefix(err.Error(), "can't watch key file:") {
			t.Fatalf("unexpected watch error: %v", err)
		}
	}
}

func TestValidPairValidPair(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	copyPair("./testdata/server1.crt", "./testdata/server1.key")

	cm, err := certman.New("./testdata/server.crt", "./testdata/server.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	logGot := buf.String()

	if !strings.Contains(logGot, "certificate and key loaded") {
		t.Log("log output received:", logGot)
		t.Fatal("expected certificate loaded message not found")
	}
	if !strings.Contains(logGot, "watching for cert and key change") {
		t.Log("log output received:", logGot)
		t.Fatal("expected watching message not found")
	}

	buf.Reset()
	copyPair("./testdata/server2.crt", "./testdata/server2.key")

	time.Sleep(2500 * time.Millisecond)

	logGot = buf.String()

	if !strings.Contains(logGot, "certificate and key loaded") {
		t.Log("log output received:", logGot)
		t.Fatal("expected certificate reloaded message not found")
	}
}

func TestValidPairInvalidPair(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	copyPair("./testdata/server1.crt", "./testdata/server1.key")

	cm, err := certman.New("./testdata/server.crt", "./testdata/server.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	logGot := buf.String()

	if !strings.Contains(logGot, "certificate and key loaded") {
		t.Log("log output received:", logGot)
		t.Fatal("expected certificate loaded message not found")
	}
	if !strings.Contains(logGot, "watching for cert and key change") {
		t.Log("log output received:", logGot)
		t.Fatal("expected watching message not found")
	}

	buf.Reset()

	copyPair("./testdata/server1.crt", "./testdata/server2.key")

	time.Sleep(2500 * time.Millisecond)

	logGot = buf.String()

	if !strings.Contains(logGot, "can't load cert or key file") {
		t.Log("log output received:", logGot)
		t.Fatal("expected error message not found")
	}
}

func TestStop(t *testing.T) {
	buf := new(bytes.Buffer)
	l := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	copyPair("./testdata/server1.crt", "./testdata/server1.key")

	cm, err := certman.New("./testdata/server.crt", "./testdata/server.key", certman.WithLogger(l))
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}
	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	logGot := buf.String()

	if !strings.Contains(logGot, "certificate and key loaded") {
		t.Log("log output received:", logGot)
		t.Fatal("expected certificate loaded message not found")
	}
	if !strings.Contains(logGot, "watching for cert and key change") {
		t.Log("log output received:", logGot)
		t.Fatal("expected watching message not found")
	}

	buf.Reset()
	cm.Stop()

	copyPair("./testdata/server2.crt", "./testdata/server2.key")
	time.Sleep(200 * time.Millisecond)

	logGot = buf.String()

	if !strings.Contains(logGot, "stopped watching") {
		t.Log("log output received:", logGot)
		t.Fatal("expected stopped watching message not found")
	}
}

func TestGetCertificate(t *testing.T) {
	cm, err := certman.New("./testdata/server1.crt", "./testdata/server1.key")
	if err != nil {
		t.Fatalf("could not create certman: %v", err)
	}

	if err := cm.Watch(); err != nil {
		t.Fatalf("could not watch files: %v", err)
	}

	hello := &tls.ClientHelloInfo{}

	cmCert, err := cm.GetCertificate(hello)
	if err != nil {
		t.Fatalf("could not get certman certificate")
	}

	expectedCert, _ := tls.LoadX509KeyPair("./testdata/server1.crt", "./testdata/server1.key")
	if err != nil {
		t.Fatalf("could not load certificate and key files to test: %v", err)
	}

	if !reflect.DeepEqual(cmCert.Certificate, expectedCert.Certificate) {
		t.Fatalf("certman certificate doesn't match expected certificate")
	}
}

func copyPair(crt, key string) {
	// ignore error handling
	crtSource, _ := os.Open(crt)
	defer crtSource.Close()

	crtDest, _ := os.Create("./testdata/server.crt")
	defer crtDest.Close()

	io.Copy(crtDest, crtSource)

	keySource, _ := os.Open(key)
	defer keySource.Close()

	keyDest, _ := os.Create("./testdata/server.key")
	defer keyDest.Close()

	io.Copy(keyDest, keySource)
}
