package pkg

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go StartApp()
	waitServer()
	run := m.Run()
	println("Test end")
	os.Exit(run)
}

func TestWithoutTokenItShouldReturnNotAuthorize(t *testing.T) {
	code, _ := getCall("/v1/auth", NoToken)
	if code != http.StatusUnauthorized {
		t.Errorf("Incorrect Status code %v", code)
	}
}

func TestTokenEndpointShouldReturnValidEndpoint(t *testing.T) {
	code, body := postCall("/token", NoToken)
	token := unmarshal[tokenResponse](body)
	if code != http.StatusOK {
		t.Errorf("Incorrect Status code %v", code)
	}

	code, _ = getCall("/v1/auth", token.Token)
	if code != http.StatusOK {
		t.Errorf("Incorrect Status code %v", code)
	}
}

const NoToken = ""

func waitServer() {
	code, _ := getCall("/health", NoToken)
	for code != http.StatusOK {
		code, _ = getCall("/health", NoToken)
	}
}

func getCall(path string, token string) (int, string) {
	return call("GET", path, token)
}

func postCall(path string, token string) (int, string) {
	return call("POST", path, token)
}

func call(method string, path string, token string) (int, string) {
	fullUrl := fmt.Sprintf("http://localhost:%d%v", 8080, path)
	var httpClient = http.Client{
		Timeout: time.Duration(24) * time.Hour,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest(method, fullUrl, nil)
	if token != NoToken {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Print(err.Error())
		return 0, err.Error()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return resp.StatusCode, string(body)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func unmarshal[K any](value string) K {
	var result = new(K)
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return *result
}
