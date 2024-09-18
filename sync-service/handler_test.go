package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestServerServeHTTP(t *testing.T) {
	t.Parallel()

	apiKey := "test"
	srv := NewServer(NewMemory(), apiKey, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	for _, tc := range []struct {
		name      string
		key       string
		url       string
		method    string
		expStatus int
	}{
		{
			name:      "index always visible",
			url:       "/",
			method:    http.MethodGet,
			expStatus: http.StatusOK,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			res := httptest.NewRecorder()
			srv.ServeHTTP(res, req)
			if res.Result().StatusCode != tc.expStatus {
				t.Errorf("exp %v, got %v", tc.expStatus, res.Result().StatusCode)
			}
		})
	}
}

func TestSyncGet(t *testing.T) {
	t.Parallel()

	now := time.Now()
	mem := NewMemory()

	items := []Item{
		{ID: "id-0", Kind: KindEvent, Updated: now.Add(-10 * time.Minute)},
		{ID: "id-1", Kind: KindEvent, Updated: now.Add(-5 * time.Minute)},
		{ID: "id-2", Kind: KindTask, Updated: now.Add(time.Minute)},
	}

	for _, item := range items {
		if err := mem.Update(item); err != nil {
			t.Errorf("exp nil, got %v", err)
		}
	}

	apiKey := "test"
	srv := NewServer(mem, apiKey, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	for _, tc := range []struct {
		name      string
		ts        time.Time
		ks        []string
		expStatus int
		expItems  []Item
	}{
		{
			name:      "full",
			expStatus: http.StatusOK,
			expItems:  items,
		},
		{
			name:      "new",
			ts:        now.Add(-6 * time.Minute),
			expStatus: http.StatusOK,
			expItems:  []Item{items[1], items[2]},
		},
		{
			name:      "kind",
			ks:        []string{string(KindTask)},
			expStatus: http.StatusOK,
			expItems:  []Item{items[2]},
		},
		{
			name:      "unknown kind",
			ks:        []string{"test"},
			expStatus: http.StatusBadRequest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/sync?ts=%s", url.QueryEscape(tc.ts.Format(time.RFC3339)))
			if len(tc.ks) > 0 {
				url = fmt.Sprintf("%s&ks=%s", url, strings.Join(tc.ks, ","))
			}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
			res := httptest.NewRecorder()
			srv.ServeHTTP(res, req)

			if res.Result().StatusCode != tc.expStatus {
				t.Errorf("exp %v, got %v", tc.expStatus, res.Result().StatusCode)
			}
			if tc.expStatus != http.StatusOK {
				return
			}

			var actItems []Item
			actBody, err := io.ReadAll(res.Result().Body)
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			defer res.Result().Body.Close()

			if err := json.Unmarshal(actBody, &actItems); err != nil {
				t.Errorf("exp nil, got %v", err)
			}

			if len(actItems) != len(tc.expItems) {
				t.Errorf("exp %d, got %d", len(tc.expItems), len(actItems))
			}
			sort.Slice(actItems, func(i, j int) bool {
				return actItems[i].ID < actItems[j].ID
			})
			for i := range actItems {
				if actItems[i].ID != tc.expItems[i].ID {
					t.Errorf("exp %v, got %v", tc.expItems[i].ID, actItems[i].ID)
				}
			}
		})
	}
}

func TestSyncPost(t *testing.T) {
	t.Parallel()

	apiKey := "test"
	for _, tc := range []struct {
		name      string
		reqBody   []byte
		expStatus int
		expItems  []Item
	}{
		{
			name:      "empty",
			expStatus: http.StatusBadRequest,
		},
		{
			name:      "invalid json",
			reqBody:   []byte(`{"fail}`),
			expStatus: http.StatusBadRequest,
		},
		{
			name: "invalid item",
			reqBody: []byte(`[
  {"id":"id-1","kind":"event","updated":"2024-09-06T08:00:00Z"},
]`),
			expStatus: http.StatusBadRequest,
		},
		{
			name: "normal",
			reqBody: []byte(`[
  {"id":"id-1","kind":"event","updated":"2024-09-06T08:00:00Z","deleted":false,"body":"item"},
  {"id":"id-2","kind":"event","updated":"2024-09-06T08:12:00Z","deleted":false,"body":"item2"}
]`),
			expStatus: http.StatusNoContent,
			expItems: []Item{
				{ID: "id-1", Kind: KindEvent, Updated: time.Date(2024, 9, 6, 8, 0, 0, 0, time.UTC)},
				{ID: "id-2", Kind: KindEvent, Updated: time.Date(2024, 9, 6, 12, 0, 0, 0, time.UTC)},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mem := NewMemory()
			srv := NewServer(mem, apiKey, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
			req, err := http.NewRequest(http.MethodPost, "/sync", bytes.NewBuffer(tc.reqBody))
			if err != nil {
				t.Errorf("exp nil, got %v", err)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
			res := httptest.NewRecorder()
			srv.ServeHTTP(res, req)

			if res.Result().StatusCode != tc.expStatus {
				t.Errorf("exp %v, got %v", tc.expStatus, res.Result().StatusCode)
			}

			actItems, err := mem.Updated([]Kind{}, time.Time{})
			if err != nil {
				t.Errorf("exp nil, git %v", err)
			}
			if len(actItems) != len(tc.expItems) {
				t.Errorf("exp %d, got %d", len(tc.expItems), len(actItems))
			}
			sort.Slice(actItems, func(i, j int) bool {
				return actItems[i].ID < actItems[j].ID
			})
			for i := range actItems {
				if actItems[i].ID != tc.expItems[i].ID {
					t.Errorf("exp %v, got %v", tc.expItems[i].ID, actItems[i].ID)
				}
			}
		})
	}
}
