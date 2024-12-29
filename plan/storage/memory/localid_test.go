package memory_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-mod.ewintr.nl/planner/plan/storage"
	"go-mod.ewintr.nl/planner/plan/storage/memory"
)

func TestLocalID(t *testing.T) {
	t.Parallel()

	repo := memory.NewLocalID()

	t.Log("start empty")
	actIDs, actErr := repo.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actIDs) != 0 {
		t.Errorf("exp nil, got %v", actErr)
	}

	t.Log("next id")
	actNext, actErr := repo.Next()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actNext != 1 {
		t.Errorf("exp 1, got %v", actNext)
	}

	t.Log("store")
	if actErr = repo.Store("test", 1); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}

	t.Log("retrieve known")
	actLid, actErr := repo.FindOrNext("test")
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actLid != 1 {
		t.Errorf("exp 1, git %v", actLid)
	}
	t.Log("retrieve unknown")
	actLid, actErr = repo.FindOrNext("new")
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actLid != 2 {
		t.Errorf("exp 2, got %v", actLid)
	}

	t.Log("find by local id")
	actID, actErr := repo.FindOne(1)
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if actID != "test" {
		t.Errorf("exp test, got %v", actID)
	}

	t.Log("unknown local id")
	actID, actErr = repo.FindOne(2)
	if !errors.Is(actErr, storage.ErrNotFound) {
		t.Errorf("exp ErrNotFound, got %v", actErr)
	}

	actIDs, actErr = repo.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	expIDs := map[string]int{
		"test": 1,
	}
	if diff := cmp.Diff(expIDs, actIDs); diff != "" {
		t.Errorf("(exp +, got -)\n%s", diff)
	}

	t.Log("delete")
	if actErr = repo.Delete("test"); actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	actIDs, actErr = repo.FindAll()
	if actErr != nil {
		t.Errorf("exp nil, got %v", actErr)
	}
	if len(actIDs) != 0 {
		t.Errorf("exp 0, got %v", actErr)
	}

	t.Log("delete non-existing")
	actErr = repo.Delete("non-existing")
	if !errors.Is(actErr, storage.ErrNotFound) {
		t.Errorf("exp %v, got %v", storage.ErrNotFound, actErr)
	}
}
