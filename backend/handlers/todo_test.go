package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"todo-list/backend/database"
	"todo-list/backend/router"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setup() *gin.Engine {
	database.Init(":memory:")
	return router.Setup("http://localhost")
}

func TestCreateTodo(t *testing.T) {
	r := setup()
	body := `{"title":"Test todo","description":"desc","priority":"high","tags":"work,urgent"}`
	req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["title"] != "Test todo" {
		t.Errorf("title mismatch: %v", resp["title"])
	}
}

func TestListTodos(t *testing.T) {
	r := setup()
	for _, title := range []string{"Todo A", "Todo B"} {
		body := `{"title":"` + title + `","priority":"medium"}`
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	req := httptest.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var todos []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &todos)
	if len(todos) < 2 {
		t.Errorf("expected at least 2 todos, got %d", len(todos))
	}
}

func TestToggleTodo(t *testing.T) {
	r := setup()
	body := `{"title":"Toggle me"}`
	req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := int(created["id"].(float64))

	url := "/api/todos/" + strconv.Itoa(id) + "/toggle"
	req = httptest.NewRequest("PATCH", url, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["completed"] != true {
		t.Errorf("expected completed=true, got %v", resp["completed"])
	}
}

func TestDeleteTodo(t *testing.T) {
	r := setup()
	body := `{"title":"Delete me"}`
	req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := int(created["id"].(float64))

	req = httptest.NewRequest("DELETE", "/api/todos/"+strconv.Itoa(id), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestFilterByTag(t *testing.T) {
	r := setup()
	body := `{"title":"Tagged","tags":"test-tag","priority":"medium"}`
	req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/api/todos?tag=test-tag", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var todos []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &todos)
	if len(todos) == 0 {
		t.Error("expected at least 1 todo with tag")
	}
}

func TestExportJSON(t *testing.T) {
	r := setup()
	req := httptest.NewRequest("GET", "/api/todos/export?format=json", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestExportCSV(t *testing.T) {
	r := setup()
	req := httptest.NewRequest("GET", "/api/todos/export?format=csv", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestReorder(t *testing.T) {
	r := setup()
	// Create 3 todos
	for _, title := range []string{"A", "B", "C"} {
		body := `{"title":"` + title + `"}`
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// Get IDs
	req := httptest.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var todos []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &todos)
	ids := make([]int, 0, len(todos))
	for _, td := range todos {
		ids = append(ids, int(td["id"].(float64)))
	}

	// Reverse order
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}
	body := `{"ids":[` + strings.Join(idsStr, ",") + `]}`
	req = httptest.NewRequest("PUT", "/api/todos/reorder", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
