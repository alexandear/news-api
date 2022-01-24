package test

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	serverURL = "http://localhost:8080"
)

type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	cmd := exec.Command("docker-compose", "-f", "../docker-compose.dev.yaml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) TearDownSuite() {
	cmd := exec.Command("docker-compose", "-f", "../docker-compose.dev.yaml", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) Test_EndToEnd_GetAllPosts() {
	s.AssertRequestResponse(http.MethodGet, "/posts", "",
		http.StatusOK, `{"posts":[]}`)
}

func (s *e2eTestSuite) AssertRequestResponse(reqMethod, reqPath, reqBody string, expectedStatus int, expectedBody string) {
	s.T().Helper()
	req := s.NewRequest(reqMethod, reqPath, reqBody)
	resp := s.DoRequest(req)
	s.EqualResponse(expectedStatus, expectedBody, resp)
}

func (s *e2eTestSuite) NewRequest(method, path, body string) *http.Request {
	req, err := http.NewRequest(method, serverURL+path, strings.NewReader(body))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (s *e2eTestSuite) DoRequest(req *http.Request) *http.Response {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *e2eTestSuite) EqualResponse(expectedStatusCode int, expectedBody string, actual *http.Response) {
	s.T().Helper()
	s.Require().NotNil(actual)
	s.Require().NotNil(actual.Body)

	s.Equal(expectedStatusCode, actual.StatusCode)

	byteBody, err := io.ReadAll(actual.Body)
	s.Require().NoError(err)
	s.Equal(expectedBody, strings.Trim(string(byteBody), "\n"))

	s.Require().NoError(actual.Body.Close())
}
