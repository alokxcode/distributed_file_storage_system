// main_test.go
package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAndRetrieve(t *testing.T) {
    // setup
    s1 := makeNodeServer(":3000")
    s2 := makeNodeServer(":4000", ":3000")

    go s1.Start()
    go s2.Start()
    time.Sleep(50 * time.Millisecond) // let both boot and connect

    defer s1.Stop()
    defer s2.Stop()

    // act
    testData := []byte("hello from server2! : PORT - 4000")
    err := s2.Store("test.jpeg", bytes.NewReader(testData))
    require.NoError(t, err)
}

func TestBootstrapConnects(t *testing.T) {
    s1 := makeNodeServer(":3001")
    s2 := makeNodeServer(":4001", ":3001")

    go s1.Start()
    go s2.Start()
    time.Sleep(50 * time.Millisecond)

    defer s1.Stop()
    defer s2.Stop()

    // assert s2 actually connected to s1
    assert.Equal(t, 1, s2.Peers)
}
