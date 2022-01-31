package test

import (
	"encoding/json"
	"fmt"
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
}

func TestE2ETestSuite(t *testing.T) {
	var e2ets e2eTestSuite
	suite.Run(t, &e2ets)
}

func (s *e2eTestSuite) SetupSuite() {
	cmd := newDockerComposeCmd("build")
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) SetupTest() {
	cmd := newDockerComposeCmd("up", "-d")
	s.Require().NoError(cmd.Run())
	time.Sleep(4 * time.Second)
}

func (s *e2eTestSuite) TearDownTest() {
	cmd := newDockerComposeCmd("down")
	s.Require().NoError(cmd.Run())
}

func (s *e2eTestSuite) Test_EndToEnd_GetAllPosts() {
	s.Run("200 when no posts", func() {
		req := s.NewRequest(http.MethodGet, "/posts", "")

		resp := s.DoRequest(req)

		s.EqualResponse(http.StatusOK, `{"posts":[]}`, resp)
	})

	s.Run("200 when few posts", func() {
		for i := 0; i < 3; i++ {
			s.createPost(fmt.Sprintf("Title %d", i), fmt.Sprintf("Content %d", i))
		}
		req := s.NewRequest(http.MethodGet, "/posts", "")

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusOK, resp)

		s.Require().NotNil(resp.Body)
		type respPost struct {
			ID        string    `json:"id"`
			Title     string    `json:"title"`
			Content   string    `json:"content"`
			UpdatedAt time.Time `json:"updated_at"`
			CreatedAt time.Time `json:"created_at"`
		}
		type respJSON struct {
			Posts []respPost `json:"posts"`
		}

		var actual respJSON
		s.Require().NoError(json.NewDecoder(resp.Body).Decode(&actual))
		s.Len(actual.Posts, 3)
	})
}

func (s *e2eTestSuite) Test_EndToEnd_CreatePost() {
	s.Run("200 ok", func() {
		s.createPost("Post Title", "Post Content")
	})

	s.Run("400 when empty body", func() {
		req := s.NewRequest(http.MethodPost, "/posts", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})

	s.Run("400 when invalid title", func() {
		req := s.NewRequest(http.MethodPost, "/posts", `{"title":"po","content":"Post Content"}`)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) createPost(title, content string) string {
	req := s.NewRequest(http.MethodPost, "/posts", fmt.Sprintf(`{"title":"%s","content":"%s"}`, title, content))
	resp := s.DoRequest(req)
	s.EqualStatusCode(http.StatusOK, resp)
	s.Require().NotNil(resp.Body)
	type respJSON struct {
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		ID        string    `json:"id"`
	}

	var actual respJSON
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&actual))

	s.True(actual.UpdatedAt.IsZero())
	s.False(actual.CreatedAt.IsZero())
	s.IsUUID(actual.ID)

	s.Require().NoError(resp.Body.Close())

	return actual.ID
}

func (s *e2eTestSuite) Test_EndToEnd_GetPost() {
	s.Run("200 when success", func() {
		id := s.createPost("Title", "Content")

		s.getPost(id, "Title", "Content")
	})

	s.Run("404 when post not found", func() {
		req := s.NewRequest(http.MethodGet, "/posts/"+uuid.NewString(), ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)
	})

	s.Run("400 when invalid post id", func() {
		req := s.NewRequest(http.MethodGet, "/posts/abc", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) getPost(id, expectedTitle, expectedContent string) {
	req := s.NewRequest(http.MethodGet, "/posts/"+id, ``)

	resp := s.DoRequest(req)

	s.EqualStatusCode(http.StatusOK, resp)
	s.Require().NotNil(resp.Body)
	type respJSON struct {
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		ID        string    `json:"id"`
	}

	var actual respJSON
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&actual))
	s.Equal(id, actual.ID)
	s.Equal(expectedTitle, actual.Title)
	s.Equal(expectedContent, actual.Content)
}

func (s *e2eTestSuite) Test_EndToEnd_UpdatePost() {
	s.Run("200 when update", func() {
		id := s.createPost("Title", "Content")
		req := s.NewRequest(http.MethodPut, "/posts/"+id, `{"title":"New title","content":"New content"}`)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusOK, resp)
		s.Require().NotNil(resp.Body)

		s.getPost(id, "New title", "New content")
	})

	s.Run("404 when post not found", func() {
		req := s.NewRequest(http.MethodPut, "/posts/"+uuid.NewString(), `{"title":"Title","content":"Content"}`)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)
	})

	s.Run("400 when invalid post id", func() {
		req := s.NewRequest(http.MethodPut, "/posts/abc", ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusBadRequest, resp)
	})
}

func (s *e2eTestSuite) Test_EndToEnd_DeletePost() {
	s.Run("200 when success delete", func() {
		id := s.createPost("Title", "Content")
		req := s.NewRequest(http.MethodDelete, "/posts/"+id, ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusOK, resp)
		{
			req := s.NewRequest(http.MethodGet, "/posts/"+id, ``)
			resp := s.DoRequest(req)
			s.EqualStatusCode(http.StatusNotFound, resp)
		}
	})

	s.Run("404 when post not found", func() {
		req := s.NewRequest(http.MethodDelete, "/posts/"+uuid.NewString(), ``)

		resp := s.DoRequest(req)

		s.EqualStatusCode(http.StatusNotFound, resp)

	})

	s.Run("400 when invalid post id", func() {
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

	client := http.Client{
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
