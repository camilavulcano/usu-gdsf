package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/jak103/usu-gdsf/auth"
	"github.com/jak103/usu-gdsf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// _db, _ = db.NewDatabaseFromEnv()

	game0 = models.Game{
		Name:         "game0",
		Rating:       3.5,
		TimesPlayed:  1,
		ImagePath:    "path/0",
		Description:  "dummy game 0",
		Developer:    "tester",
		CreationDate: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		Version:      "0.0.0",
		Tags:         []string{"tag0", "tag1"},
		Downloads:    35,
		DownloadLink: "dummy.test",
	}

	game1 = models.Game{
		Name:         "game1",
		Rating:       3.9,
		TimesPlayed:  2,
		ImagePath:    "path/1",
		Description:  "dummy game 1",
		Developer:    "tester",
		CreationDate: time.Date(1900, 1, 2, 0, 0, 0, 0, time.UTC),
		Version:      "0.0.1",
		Tags:         []string{"tag1", "tag2"},
		Downloads:    36,
		DownloadLink: "dummy1.test",
	}
)

func TestGame(t *testing.T) {
	e := echo.New()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/game", nil)
	c := e.NewContext(request, recorder)

	if assert.NoError(t, getAllGames(c)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	}
}

func TestGetAllGames(t *testing.T) {
	e := echo.New()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/games", nil)
	c := e.NewContext(request, recorder)

	if assert.NoError(t, getAllGames(c)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	}
}

func TestGetGamesWithTags(t *testing.T) {
	//cleanup
	e := echo.New()

	t.Cleanup(func() {
		_db.RemoveGame(game0)
		_db.RemoveGame(game1)
	})

	// add the games, assign ids
	id0, _ := _db.AddGame(game0)
	id1, _ := _db.AddGame(game1)
	game0.Id = id0
	game1.Id = id1

	params := auth.TokenParams{
		Type:      auth.ACCESS_TOKEN,
		UserId:    42,
		UserEmail: "tst@example.com",
	}

	token := auth.GenerateToken(params)

	q := make(url.Values)
	q.Set("tags", "tag0-tag1")

	req := httptest.NewRequest("http.MethodGet", "/games/tags?"+q.Encode(), nil)
	req.Header.Set("accessToken", token)
	//response writer
	// we can inspect the ResponseRecorder output which is response generated by handler
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	assert.NoError(t, getGamesWithTags(c))
	response := recorder.Body.String()
	gameObjectResponse := []models.Game{}

	in := []byte(response)
	err := json.Unmarshal(in, &gameObjectResponse)
	if err != nil {
		println(err)
	}
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, 2, len(gameObjectResponse))
}

func TestGetAllGamesReturnsCorrectNumberOfGames(t *testing.T) {
	//cleanup
	e := echo.New()

	t.Cleanup(func() {
		_db.RemoveGame(game0)
		_db.RemoveGame(game1)
	})

	// add the games, assign ids
	id0, _ := _db.AddGame(game0)
	id1, _ := _db.AddGame(game1)
	game0.Id = id0
	game1.Id = id1

	params := auth.TokenParams{
		Type:      auth.ACCESS_TOKEN,
		UserId:    42,
		UserEmail: "tst@example.com",
	}

	token := auth.GenerateToken(params)

	// q := make(url.Values)
	// q.Set("tags", "tag0-tag1")

	req := httptest.NewRequest("http.MethodGet", "/games", nil)
	req.Header.Set("accessToken", token)
	//response writer
	// we can inspect the ResponseRecorder output which is response generated by handler
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)

	assert.NoError(t, getAllGames(c))
	response := recorder.Body.String()
	gameObjectResponse := []models.Game{}

	in := []byte(response)
	err := json.Unmarshal(in, &gameObjectResponse)
	if err != nil {
		println(err)
	}

	require.Equal(t, http.StatusOK, recorder.Code)

	assert.Equal(t, 10, len(gameObjectResponse))
}
