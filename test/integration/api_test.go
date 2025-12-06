package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/controller"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/internal/server"
	"github.com/serediukit/civix-backend/internal/services"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	ctx        context.Context
	dbSetup    *TestDatabaseSetup
	router     *gin.Engine
	jwtService *jwt.JWT
	cacheRepo  repository.CacheRepository
}

func (suite *APITestSuite) SetupSuite() {
	suite.ctx = context.Background()
	gin.SetMode(gin.TestMode)

	var err error
	suite.dbSetup, err = SetupTestDatabase(suite.ctx)
	suite.Require().NoError(err)

	// Initialize repositories
	userRepo := repository.NewUserRepository(suite.dbSetup.Store)
	reportRepo := repository.NewReportRepository(suite.dbSetup.Store)
	cityRepo := repository.NewCityRepository(suite.dbSetup.Store)
	suite.cacheRepo = repository.NewCacheRepository(suite.dbSetup.RedisClient)

	// Initialize JWT
	suite.jwtService = jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	// Initialize services
	authService := services.NewAuthService(userRepo, cityRepo, suite.cacheRepo, suite.jwtService)
	userService := services.NewUserService(userRepo)
	reportService := services.NewReportService(reportRepo, cityRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	reportController := controller.NewReportController(reportService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(suite.jwtService, suite.cacheRepo)

	// Setup router
	logger := logrus.New()
	suite.router = server.SetupRouter(authController, userController, reportController, authMiddleware, logger)
}

func (suite *APITestSuite) TearDownSuite() {
	if suite.dbSetup != nil {
		suite.dbSetup.Teardown(suite.ctx)
	}
}

func (suite *APITestSuite) TearDownTest() {
	// Clean up Redis between tests
	suite.dbSetup.RedisClient.GetClient().FlushAll(suite.ctx)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
}

func (suite *APITestSuite) TestRegisterAndLogin() {
	// First insert a city for testing
	insertCitySQL := `
		INSERT INTO cities (city_id, name, region, location)
		VALUES ('123e4567-e89b-12d3-a456-426614174000', 'Test City', 'Test Region',
				ST_SetSRID(ST_MakePoint(30.5234, 50.4501), 4326))
		ON CONFLICT DO NOTHING
	`
	_, err := suite.dbSetup.Store.GetDB().Exec(suite.ctx, insertCitySQL)
	suite.Require().NoError(err)

	// Test Registration
	registerReq := contracts.RegisterRequest{
		Email:       "integration@test.com",
		Password:    "password123",
		Name:        "Integration",
		Surname:     "Test",
		PhoneNumber: "+1234567890",
		Location: model.Location{
			Lat: 50.4501,
			Lng: 30.5234,
		},
	}

	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var registerResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), registerResp["success"].(bool))

	// Test Login
	loginReq := contracts.LoginRequest{
		Email:    "integration@test.com",
		Password: "password123",
	}

	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), loginResp["success"].(bool))

	// Check for tokens in response
	data, ok := loginResp["data"].(map[string]interface{})
	assert.True(suite.T(), ok, "response should have data field")
	if ok {
		accessToken, _ := data["access_token"].(map[string]interface{})
		refreshToken, _ := data["refresh_token"].(map[string]interface{})
		assert.NotNil(suite.T(), accessToken)
		assert.NotNil(suite.T(), refreshToken)
		assert.NotEmpty(suite.T(), accessToken["token"])
		assert.NotEmpty(suite.T(), refreshToken["token"])
	}
}

func (suite *APITestSuite) TestLoginWithInvalidCredentials() {
	loginReq := contracts.LoginRequest{
		Email:    "nonexistent@test.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *APITestSuite) TestUnauthorizedAccess() {
	// Try to access protected endpoint without token
	req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *APITestSuite) TestAuthenticatedEndpointsFlow() {
	// First insert a city for testing
	insertCitySQL := `
		INSERT INTO cities (city_id, name, region, location)
		VALUES ('223e4567-e89b-12d3-a456-426614174000', 'Auth City', 'Test Region',
				ST_SetSRID(ST_MakePoint(30.5234, 50.4501), 4326))
		ON CONFLICT DO NOTHING
	`
	_, err := suite.dbSetup.Store.GetDB().Exec(suite.ctx, insertCitySQL)
	suite.Require().NoError(err)

	// Create a user
	registerReq := contracts.RegisterRequest{
		Email:       "authflow@test.com",
		Password:    "password123",
		Name:        "Auth",
		Surname:     "Flow",
		PhoneNumber: "+9999999999",
		Location: model.Location{
			Lat: 50.4501,
			Lng: 30.5234,
		},
	}

	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Require().Equal(http.StatusCreated, w.Code)

	// Login to get tokens
	loginReq := contracts.LoginRequest{
		Email:    "authflow@test.com",
		Password: "password123",
	}

	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	suite.Require().Equal(http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	suite.Require().NoError(err)

	// Extract token
	data := loginResp["data"].(map[string]interface{})
	accessToken := data["access_token"].(map[string]interface{})
	tokenString := accessToken["token"].(string)

	// Test accessing protected endpoint with token
	req, _ = http.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var userResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), userResp["success"].(bool))
}
