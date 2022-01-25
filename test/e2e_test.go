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
	cmd := newDockerComposeCmd("up", "-d", "--build")
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
	req := s.NewRequest(http.MethodPost, "/posts", `{"title":"Post Title","content":"Post Content"}`)

	resp := s.DoRequest(req)

	s.EqualStatusCode(http.StatusOK, resp)

	s.Require().NotNil(resp.Body)
	type respJSON struct {
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		ID        string    `json:"id"`
	}

	actual := &respJSON{}
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(actual))
	s.True(actual.UpdatedAt.IsZero())
	s.False(actual.CreatedAt.IsZero())
	s.IsUUID(actual.ID)

	s.Require().NoError(resp.Body.Close())
}

func (s *e2eTestSuite) Test_EndToEnd_GetPost() {
	s.Run("post not found", func() {
		req := s.NewRequest(http.MethodGet, "/posts/"+uuid.NewString(), ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)
	})

	s.Run("invalid post id", func() {
		req := s.NewRequest(http.MethodGet, "/posts/abc", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) Test_EndToEnd_UpdatePost() {
	s.Run("post not found", func() {
		req := s.NewRequest(http.MethodPut, "/posts/"+uuid.NewString(), ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)
	})

	s.Run("invalid post id", func() {
		req := s.NewRequest(http.MethodPut, "/posts/abc", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) Test_EndToEnd_DeletePost() {
	s.Run("post not found", func() {
		req := s.NewRequest(http.MethodDelete, "/posts/"+uuid.NewString(), ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)
	})

	s.Run("invalid post id", func() {
		req := s.NewRequest(http.MethodDelete, "/posts/abc", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) NewRequest(method, path, body string) *http.Request {
	s.T().Helper()

	req, err := http.NewRequest(method, serverURL+path, strings.NewReader(body))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (s *e2eTestSuite) DoRequest(req *http.Request) *http.Response {
	s.T().Helper()

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *e2eTestSuite) EqualStatusCode(expectedStatusCode int, actual *http.Response) {
	s.T().Helper()

	s.Require().NotNil(actual)

	s.Equal(expectedStatusCode, actual.StatusCode)
}

func (s *e2eTestSuite) EqualResponse(expectedStatusCode int, expectedBody string, actual *http.Response) {
	s.T().Helper()

	s.EqualStatusCode(expectedStatusCode, actual)

	s.Require().NotNil(actual.Body)
	byteBody, err := io.ReadAll(actual.Body)
	s.Require().NoError(err)

	s.Equal(expectedBody, strings.Trim(string(byteBody), "\n"))

	s.Require().NoError(actual.Body.Close())
}

func (s *e2eTestSuite) IsUUID(actual string) {
	s.T().Helper()

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
