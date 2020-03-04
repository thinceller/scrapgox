package client

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

var testClient *Client

const testToken = "test_token"
const testProject = "test"

func setup() {
	parsedUrl, _ := url.ParseRequestURI(DefaultHost)

	c, _ := NewClient(parsedUrl, testToken, DefaultUserAgent)

	testClient = c
}

func TestNewClient(t *testing.T) {
	if testClient.URL.String() != DefaultHost {
		t.Fatal("failed test")
	}
	if testClient.Token != testToken {
		t.Fatal("failed test")
	}
}

func TestBuildRequest(t *testing.T) {
	path := "api/help-jp?skip=0&limi=30"
	req, err := testClient.buildRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	expectedUrl := fmt.Sprintf("%s/%s", testClient.URL.String(), path)
	if req.URL.String() != expectedUrl {
		t.Fatal("failed to build request url")
	}
	if req.Method != "GET" {
		t.Fatal("failed to set request method")
	}
	cookie := req.Header.Get("Cookie")
	if cookie != "connect.sid="+testToken {
		t.Fatal("failed to set request token")
	}
}

func TestBuildPath(t *testing.T) {
	emptyQuery := ""
	existingQuery := "test"

	path := buildPath(testProject, 0, 30, emptyQuery)
	pathByQuery := buildPath(testProject, 0, 30, existingQuery)

	if path != fmt.Sprintf("api/pages/%s?skip=0&limit=30&q=", testProject) {
		t.Fatal("failed to build empty query path")
	}
	if pathByQuery != fmt.Sprintf("api/pages/%s/search/query?skip=0&limit=30&q=%s", testProject, existingQuery) {
		t.Fatal("failed to build query existing path")
	}
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()

	os.Exit(exitCode)
}
