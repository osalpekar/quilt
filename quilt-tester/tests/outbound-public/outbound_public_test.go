package main

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/quilt/quilt/api"
	"github.com/quilt/quilt/api/client"
	"github.com/quilt/quilt/connection/credentials"
	"github.com/quilt/quilt/db"
	"github.com/quilt/quilt/stitch"
)

func TestOutboundPublic(t *testing.T) {
	clnt, err := client.New(api.DefaultSocket, credentials.Insecure{})
	if err != nil {
		t.Fatalf("couldn't get quiltctl client: %s", err)
	}
	defer clnt.Close()

	containers, err := clnt.QueryContainers()
	if err != nil {
		t.Fatalf("couldn't query containers: %s", err)
	}

	connections, err := clnt.QueryConnections()
	if err != nil {
		t.Fatalf("couldn't query connections: %s", err)
	}

	test(t, containers, connections)
}

var testPort = 80
var testHost = fmt.Sprintf("google.com:%d", testPort)

func test(t *testing.T, containers []db.Container, connections []db.Connection) {
	connected := map[string]struct{}{}
	for _, conn := range connections {
		if conn.To == stitch.PublicInternetLabel &&
			inRange(testPort, conn.MinPort, conn.MaxPort) {
			connected[conn.From] = struct{}{}
		}
	}

	for _, c := range containers {
		shouldErr := !containsAny(connected, c.Labels)

		fmt.Printf("Fetching %s from container %s\n", testHost, c.StitchID)
		if shouldErr {
			fmt.Println(".. It should fail")
		} else {
			fmt.Println(".. It should not fail")
		}

		out, err := exec.Command("quilt", "ssh", c.StitchID,
			"wget", "-T", "2", "-O", "-", testHost).CombinedOutput()

		errored := err != nil
		if !shouldErr && errored {
			t.Errorf("Fetch failed when it should have succeeded: %s", err)
			fmt.Println(string(out))
		} else if shouldErr && !errored {
			t.Error("Fetch succeeded when it should have failed")
			fmt.Println(string(out))
		}
	}
}

func inRange(candidate, min, max int) bool {
	return min <= candidate && candidate <= max
}

func containsAny(m map[string]struct{}, keys []string) bool {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}
