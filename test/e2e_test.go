package test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const (
	serverURL = "http://localhost:8080"
)

type e2eTestSuite struct {
	suite.Suite

	dockerCompose *exec.Cmd
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	cmd := newDockerComposeCmd("up", "-d")
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) TearDownSuite() {
	cmd := newDockerComposeCmd("down")
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) TearDownTest() {
	cmd := newDockerComposeCmd("up", "news-clean-postgres")
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) Test_EndToEnd_GetAllPostsEmpty() {
	req := s.NewRequest(http.MethodGet, "/posts", "")

	resp := s.DoRequest(req)

	s.EqualResponse(http.StatusOK, `{"posts":[]}`, resp)
}

func (s *e2eTestSuite) Test_EndToEnd_CreatePost() {
	req := s.NewRequest(http.MethodPost, "/posts", `{"title":"Post Title","content":"Content"}`)

	resp := s.DoRequest(req)

	s.Require().NotNil(resp)
	s.Require().NotNil(resp.Body)
	s.Equal(http.StatusOK, resp.StatusCode)

	type respJSON struct {
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		ID        string    `json:"id"`
	}

	actual := &respJSON{}
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(actual))
	s.False(actual.UpdatedAt.IsZero())
	s.False(actual.CreatedAt.IsZero())
	s.IsUUID(actual.ID)

	s.Require().NoError(resp.Body.Close())
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

func (s *e2eTestSuite) IsUUID(actual string) {
	_, err := uuid.Parse(actual)
	s.NoError(err)
}

func newDockerComposeCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("docker-compose", "-f", "../docker-compose.dev.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append(cmd.Args, args...)
	return cmd
}
