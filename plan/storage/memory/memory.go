package memory

import (
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Memories struct {
	localID *LocalID
	sync    *Sync
	task    *Task
}

func New() *Memories {
	return &Memories{
		localID: NewLocalID(),
		sync:    NewSync(),
		task:    NewTask(),
	}
}

func (mems *Memories) Begin() (*storage.Tx, error) {
	return &storage.Tx{}, nil
}

func (mems *Memories) LocalID(_ *storage.Tx) storage.LocalID {
	return mems.localID
}

func (mems *Memories) Sync(_ *storage.Tx) storage.Sync {
	return mems.sync
}

func (mems *Memories) Task(_ *storage.Tx) storage.Task {
	return mems.task
}
