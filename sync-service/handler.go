package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"slices"
	"strings"
	"time"
)

type Server struct {
	syncer Syncer
	apiKey string
	logger *slog.Logger
}

func NewServer(syncer Syncer, apiKey string, logger *slog.Logger) *Server {
	return &Server{
		syncer: syncer,
		apiKey: apiKey,
		logger: logger,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/" {
		Index(w, r)
		return
	}

	if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", s.apiKey) {
		msg := "not authorized"
		http.Error(w, fmtError(msg), http.StatusUnauthorized)
		s.logger.Info(msg)
		return
	}

	head, tail := ShiftPath(r.URL.Path)
	switch {
	case head == "sync" && tail != "/":
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	case head == "sync" && r.Method == http.MethodGet:
		s.SyncGet(w, r)
	case head == "sync" && r.Method == http.MethodPost:
		s.SyncPost(w, r)
	default:
		msg := "not found"
		http.Error(w, fmtError(msg), http.StatusNotFound)
		s.logger.Info(msg)
	}
}

func (s *Server) SyncGet(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Time{}
	tsStr := r.URL.Query().Get("ts")
	if tsStr != "" {
		var err error
		if timestamp, err = time.Parse(time.RFC3339, tsStr); err != nil {
			msg := err.Error()
			http.Error(w, fmtError(msg), http.StatusBadRequest)
			s.logger.Info(msg)
			return
		}
	}
	ks := make([]Kind, 0)
	ksStr := r.URL.Query().Get("ks")
	if ksStr != "" {
		for _, k := range strings.Split(ksStr, ",") {
			if !slices.Contains(KnownKinds, Kind(k)) {
				msg := fmt.Sprintf("unknown kind: %s", k)
				http.Error(w, fmtError(msg), http.StatusBadRequest)
				s.logger.Info(msg)
				return
			}
			ks = append(ks, Kind(k))
		}
	}

	items, err := s.syncer.Updated(ks, timestamp)
	if err != nil {
		msg := err.Error()
		http.Error(w, fmtError(msg), http.StatusInternalServerError)
		s.logger.Error(msg)
		return
	}

	body, err := json.Marshal(items)
	if err != nil {
		msg := err.Error()
		http.Error(w, fmtError(msg), http.StatusInternalServerError)
		s.logger.Error(msg)
		return
	}

	fmt.Fprint(w, string(body))
	s.logger.Info("served get sync")
}

func (s *Server) SyncPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg := err.Error()
		http.Error(w, fmtError(msg), http.StatusBadRequest)
		s.logger.Info(msg)
		return
	}
	defer r.Body.Close()

	var items []Item
	if err := json.Unmarshal(body, &items); err != nil {
		msg := err.Error()
		http.Error(w, fmtError(msg), http.StatusBadRequest)
		s.logger.Info(msg)
		return
	}

	for _, item := range items {
		if item.ID == "" {
			msg := "item without an id"
			http.Error(w, fmtError(msg), http.StatusBadRequest)
			s.logger.Info(msg)
			return
		}
		if item.Kind == "" {
			msg := fmt.Sprintf("item %s does not have a kind", item.ID)
			http.Error(w, fmtError(msg), http.StatusBadRequest)
			s.logger.Info(msg)
			return
		}
		if !slices.Contains(KnownKinds, item.Kind) {
			msg := fmt.Sprintf("items %s does not have a know kind", item.ID)
			http.Error(w, fmtError(msg), http.StatusBadRequest)
			s.logger.Info(msg)
			return
		}
		if item.Body == "" {
			msg := fmt.Sprintf(`{"error":"item %s does not have a body"}`, item.ID)
			http.Error(w, msg, http.StatusBadRequest)
			s.logger.Info(msg)
			return
		}
		item.Updated = time.Now()
		if err := s.syncer.Update(item); err != nil {
			msg := err.Error()
			http.Error(w, fmtError(msg), http.StatusInternalServerError)
			s.logger.Error(msg)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)

	s.logger.Info("served get sync")
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
// See https://blog.merovius.de/posts/2017-06-18-how-not-to-use-an-http-router/
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"ok"}`)
}

func fmtError(msg string) string {
	return fmt.Sprintf(`{"error":%q}`, msg)
}
