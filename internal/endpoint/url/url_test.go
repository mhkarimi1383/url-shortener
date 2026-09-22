package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/mhkarimi1383/url-shortener/constrains"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/internal/redirectcache"
	ivalidator "github.com/mhkarimi1383/url-shortener/internal/validator"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
	responseschemas "github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

func TestCreateReturnsPersistedURLID(t *testing.T) {
	e, engine, owner, _, _ := setupURLTest(t)

	recorder, err := createURLRequest(e, owner, `{
		"FullUrl":"https://example.com/original",
		"ShortCode":"created"
	}`)
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	var response responseschemas.Url
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Id == 0 {
		t.Fatal("create response ID must be non-zero")
	}

	stored := databasemodels.Url{Id: response.Id}
	has, err := engine.Get(&stored)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatalf("URL %d from create response was not persisted", response.Id)
	}
	if stored.ShortCode != response.ShortCode || stored.FullUrl != "https://example.com/original" {
		t.Fatalf("create response does not identify the persisted URL: %+v", stored)
	}
	if response.FullUrl != stored.FullUrl || response.Creator.Id != owner.Id {
		t.Fatalf("create response does not contain the complete URL: %+v", response)
	}
}

func TestUpdateChangesOnlyFullURL(t *testing.T) {
	e, engine, owner, other, _ := setupURLTest(t)
	entity := insertEntity(t, engine, owner, "first")
	otherEntity := insertEntity(t, engine, owner, "second")
	stored := insertURL(t, engine, owner, entity, "owned")
	before := getURL(t, engine, stored.Id)

	payload := fmt.Sprintf(`{
		"FullUrl":"https://example.com/updated",
		"ShortCode":"ignored",
		"Entity":%d,
		"Creator":{"Id":%d},
		"VisitCount":999
	}`, otherEntity.Id, other.Id)
	recorder, err := updateURLRequest(e, owner, fmt.Sprint(stored.Id), payload)
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}

	after := getURL(t, engine, stored.Id)
	if after.FullUrl != "https://example.com/updated" {
		t.Fatalf("FullUrl = %q, want updated URL", after.FullUrl)
	}
	if after.ShortCode != before.ShortCode {
		t.Fatalf("ShortCode changed from %q to %q", before.ShortCode, after.ShortCode)
	}
	if after.Entity.Id != before.Entity.Id {
		t.Fatalf("Entity changed from %d to %d", before.Entity.Id, after.Entity.Id)
	}
	if after.Creator.Id != before.Creator.Id {
		t.Fatalf("Creator changed from %d to %d", before.Creator.Id, after.Creator.Id)
	}
	if after.VisitCount != before.VisitCount {
		t.Fatalf("VisitCount changed from %d to %d", before.VisitCount, after.VisitCount)
	}
	if !sameTime(after.LastVisitedAt, before.LastVisitedAt) {
		t.Fatalf("LastVisitedAt changed from %v to %v", before.LastVisitedAt, after.LastVisitedAt)
	}
	if !after.CreatedAt.Equal(before.CreatedAt) {
		t.Fatalf("CreatedAt changed from %v to %v", before.CreatedAt, after.CreatedAt)
	}
	if after.Version != before.Version+1 {
		t.Fatalf("Version = %d, want %d", after.Version, before.Version+1)
	}
	if !after.UpdatedAt.After(before.UpdatedAt) {
		t.Fatalf("UpdatedAt = %v, want after %v", after.UpdatedAt, before.UpdatedAt)
	}
}

func TestUpdateEnforcesOwnershipAndExistence(t *testing.T) {
	e, engine, owner, other, admin := setupURLTest(t)
	entity := insertEntity(t, engine, owner, "first")
	stored := insertURL(t, engine, owner, entity, "protected")

	_, err := updateURLRequest(e, other, fmt.Sprint(stored.Id), `{"FullUrl":"https://example.com/forbidden"}`)
	requireHTTPErrorCode(t, err, http.StatusNotFound)
	if current := getURL(t, engine, stored.Id); current.FullUrl != stored.FullUrl {
		t.Fatalf("non-owner changed FullUrl to %q", current.FullUrl)
	}

	recorder, err := updateURLRequest(e, admin, fmt.Sprint(stored.Id), `{"FullUrl":"https://example.com/admin"}`)
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("admin update status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if current := getURL(t, engine, stored.Id); current.FullUrl != "https://example.com/admin" {
		t.Fatalf("admin update was not persisted: %q", current.FullUrl)
	}

	_, err = updateURLRequest(e, owner, fmt.Sprint(stored.Id+1), `{"FullUrl":"https://example.com/missing"}`)
	requireHTTPErrorCode(t, err, http.StatusNotFound)

	deleted := insertURL(t, engine, owner, entity, "deleted")
	if _, err := engine.ID(deleted.Id).Delete(new(databasemodels.Url)); err != nil {
		t.Fatal(err)
	}
	_, err = updateURLRequest(e, owner, fmt.Sprint(deleted.Id), `{"FullUrl":"https://example.com/deleted"}`)
	requireHTTPErrorCode(t, err, http.StatusNotFound)
}

func TestUpdateRejectsInvalidID(t *testing.T) {
	e, _, owner, _, _ := setupURLTest(t)

	_, err := updateURLRequest(e, owner, "invalid", `{"FullUrl":"https://example.com/updated"}`)
	requireHTTPErrorCode(t, err, http.StatusBadRequest)
}

func TestGetReturnsOwnedURL(t *testing.T) {
	e, engine, owner, _, _ := setupURLTest(t)
	entity := insertEntity(t, engine, owner, "first")
	stored := insertURL(t, engine, owner, entity, "details")

	recorder, err := getURLRequest(e, owner, fmt.Sprint(stored.Id))
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var response responseschemas.Url
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Id != stored.Id || response.FullUrl != stored.FullUrl || response.ShortCode != stored.ShortCode {
		t.Fatalf("unexpected URL response: %+v", response)
	}
	if response.ShortUrl != "http://example.com/details" {
		t.Fatalf("ShortUrl = %q, want %q", response.ShortUrl, "http://example.com/details")
	}
}

func TestGetEnforcesOwnershipAndExistence(t *testing.T) {
	e, engine, owner, other, admin := setupURLTest(t)
	entity := insertEntity(t, engine, owner, "first")
	stored := insertURL(t, engine, owner, entity, "protected-get")

	_, err := getURLRequest(e, other, fmt.Sprint(stored.Id))
	requireHTTPErrorCode(t, err, http.StatusNotFound)

	recorder, err := getURLRequest(e, admin, fmt.Sprint(stored.Id))
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("admin get status = %d, want %d", recorder.Code, http.StatusOK)
	}

	_, err = getURLRequest(e, owner, fmt.Sprint(stored.Id+1))
	requireHTTPErrorCode(t, err, http.StatusNotFound)

	deleted := insertURL(t, engine, owner, entity, "deleted-get")
	if _, err := engine.ID(deleted.Id).Delete(new(databasemodels.Url)); err != nil {
		t.Fatal(err)
	}
	_, err = getURLRequest(e, owner, fmt.Sprint(deleted.Id))
	requireHTTPErrorCode(t, err, http.StatusNotFound)

	_, err = getURLRequest(e, owner, "invalid")
	requireHTTPErrorCode(t, err, http.StatusBadRequest)
}

func setupURLTest(t *testing.T) (*echo.Echo, *xorm.Engine, databasemodels.User, databasemodels.User, databasemodels.User) {
	t.Helper()

	engine, err := xorm.NewEngine("sqlite", filepath.Join(t.TempDir(), "url-test.db"))
	if err != nil {
		t.Fatal(err)
	}
	engine.SetMapper(names.GonicMapper{})
	if err := engine.Sync(new(databasemodels.User), new(databasemodels.Entity), new(databasemodels.Url)); err != nil {
		t.Fatal(err)
	}

	previousEngine := database.Engine
	previousConfig := configuration.CurrentConfig
	previousCache := redirectcache.Default
	database.Engine = engine
	configuration.CurrentConfig = &configuration.Config{}
	redirectcache.Default = nil
	t.Cleanup(func() {
		database.Engine = previousEngine
		configuration.CurrentConfig = previousConfig
		redirectcache.Default = previousCache
		_ = engine.Close()
	})

	owner := databasemodels.User{Username: "owner", Password: "password"}
	other := databasemodels.User{Username: "other", Password: "password"}
	admin := databasemodels.User{Username: "admin", Password: "password", Admin: true}
	if _, err := engine.Insert(&owner, &other, &admin); err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	e.Validator = ivalidator.EchoValidator
	return e, engine, owner, other, admin
}

func createURLRequest(e *echo.Echo, user databasemodels.User, body string) (*httptest.ResponseRecorder, error) {
	request := httptest.NewRequest(http.MethodPost, "http://example.com/api/url/", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.Set(constrains.UserInfoContextVar, user)
	return recorder, Create(context)
}

func updateURLRequest(e *echo.Echo, user databasemodels.User, id, body string) (*httptest.ResponseRecorder, error) {
	request := httptest.NewRequest(http.MethodPut, "http://example.com/api/url/"+id+"/", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.SetPath("/api/url/:" + constrains.IdParamName + "/")
	context.SetParamNames(constrains.IdParamName)
	context.SetParamValues(id)
	context.Set(constrains.UserInfoContextVar, user)
	return recorder, Update(context)
}

func getURLRequest(e *echo.Echo, user databasemodels.User, id string) (*httptest.ResponseRecorder, error) {
	request := httptest.NewRequest(http.MethodGet, "http://example.com/api/url/"+id+"/", nil)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.SetPath("/api/url/:" + constrains.IdParamName + "/")
	context.SetParamNames(constrains.IdParamName)
	context.SetParamValues(id)
	context.Set(constrains.UserInfoContextVar, user)
	return recorder, Get(context)
}

func insertEntity(t *testing.T, engine *xorm.Engine, creator databasemodels.User, name string) databasemodels.Entity {
	t.Helper()
	entity := databasemodels.Entity{Name: name, Creator: creator}
	if _, err := engine.Insert(&entity); err != nil {
		t.Fatal(err)
	}
	return entity
}

func insertURL(t *testing.T, engine *xorm.Engine, creator databasemodels.User, entity databasemodels.Entity, shortCode string) databasemodels.Url {
	t.Helper()
	lastVisitedAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	stored := databasemodels.Url{
		FullUrl:       "https://example.com/" + shortCode,
		ShortCode:     shortCode,
		Creator:       creator,
		Entity:        entity,
		VisitCount:    7,
		LastVisitedAt: &lastVisitedAt,
	}
	if _, err := engine.Insert(&stored); err != nil {
		t.Fatal(err)
	}

	oldUpdatedAt := time.Now().Add(-24 * time.Hour).UTC().Truncate(time.Second)
	affected, err := engine.Table(new(databasemodels.Url)).NoAutoTime().NoVersionCheck().Cols("updated_at").Update(
		&databasemodels.Url{UpdatedAt: oldUpdatedAt},
		&databasemodels.Url{Id: stored.Id},
	)
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatalf("setting fixture UpdatedAt affected %d rows, want 1", affected)
	}
	return stored
}

func getURL(t *testing.T, engine *xorm.Engine, id int64) databasemodels.Url {
	t.Helper()
	stored := databasemodels.Url{Id: id}
	has, err := engine.Get(&stored)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatalf("URL %d was not found", id)
	}
	return stored
}

func requireHTTPErrorCode(t *testing.T, err error, code int) {
	t.Helper()
	var httpError *echo.HTTPError
	if !errors.As(err, &httpError) {
		t.Fatalf("error = %v, want Echo HTTP error %d", err, code)
	}
	if httpError.Code != code {
		t.Fatalf("HTTP error code = %d, want %d", httpError.Code, code)
	}
}

func sameTime(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.Equal(*right)
}
