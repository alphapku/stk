package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	mdw "StakeBackendGoTest/api/middleware"
	resp "StakeBackendGoTest/api/response"
	adt "StakeBackendGoTest/internal/adapter"
	mdl "StakeBackendGoTest/internal/model"
	def "StakeBackendGoTest/pkg/const"
	"StakeBackendGoTest/pkg/log"
)

// userPositions is the one users are supposed to use when parsing response on the endpoint
type userPositions struct {
	ErrorCode    int                  `json:"errCode,omitempty"`
	ErrorMessage string               `json:"errMessage,omitempty"`
	Data         *resp.StakePositions `json:"data"`
}

type apiTestSuite struct {
	suite.Suite

	*gin.Engine

	dataManager *mdl.DataManager
}

func (s *apiTestSuite) SetupSuite() {
	_ = log.Init(def.DevMode)

	s.Engine = gin.Default()

	s.dataManager = mdl.NewDataManager()

	AddRouters(s.Engine, s.dataManager)
}

func (s *apiTestSuite) TestInvalidToken() {
	req, _ := http.NewRequest("POST", equityPositionURI, nil)
	recorder := httptest.NewRecorder()
	s.Engine.ServeHTTP(recorder, req)

	resp, _ := ioutil.ReadAll(recorder.Body)
	expected := `{"errCode":-1,"data":"invalid token"}`
	s.Equal(expected, string(resp))
}

// TestValidToken tests empty data with valida token
func (s *apiTestSuite) TestValidTokenWithoutPositions() {
	s.dataManager.Reset()

	req, _ := http.NewRequest("POST", equityPositionURI, nil)
	req.Header.Set(mdw.TokenKey, mdw.TestToken)
	recorder := httptest.NewRecorder()
	s.Engine.ServeHTTP(recorder, req)

	data, _ := ioutil.ReadAll(recorder.Body)
	expected := &userPositions{}
	_ = json.Unmarshal(data, expected)

	s.Equal(expected.ErrorCode, 0)

	s.Equal(len(expected.Data.StakePositions), 0)
}

func (s *apiTestSuite) TestValidTokenWithPositions() {
	intlPositions, intlPrices, _ := adt.LoadAndParseMockData("../")
	// feed some data
	s.dataManager.Reset()
	s.dataManager.OnMessage(intlPositions)
	s.dataManager.OnMessage(intlPrices)

	req, _ := http.NewRequest("POST", equityPositionURI, nil)
	req.Header.Set(mdw.TokenKey, mdw.TestToken)
	recorder := httptest.NewRecorder()
	s.Engine.ServeHTTP(recorder, req)

	data, _ := ioutil.ReadAll(recorder.Body)
	expected := userPositions{}
	err := json.Unmarshal(data, &expected)
	s.Nil(err)
	// we just verify it does return positions, instead of checking the content as we already did this in data_manager_test.go
	s.Equal(expected.ErrorCode, 0)
	s.NotNil(expected.Data)
	s.Greater(len(expected.Data.StakePositions), 0)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}
