package http_test

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	palikka "palikka-go"
	"testing"
)

func (s *GameTestSuite) TestGame() {
	req, err := http.NewRequest(http.MethodGet, s.URL()+"/ws", nil)
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "upgrade")
	req.Header.Add("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")
	req.Header.Add("Sec-WebSocket-Version", "13")
	s.Require().NoError(err)

	s.Server.MustAuthenticate(s.T(), req, &palikka.User{ID: 1})

	res, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
}

type GameTestSuite struct {
	suite.Suite
	*Server
}

func (s *GameTestSuite) SetupSuite() {
	s.Server = MustOpenServer(s.T())
}

func (s *GameTestSuite) TearDownSuite() {
	MustCloseServer(s.T(), s.Server)
}

func (s *GameTestSuite) TearDownTest() {
	s.Server.MustResetMocks(s.T())
}

func (s *GameTestSuite) TearDownSubTest() {
	s.Server.MustResetMocks(s.T())
}

func TestGameTestSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
