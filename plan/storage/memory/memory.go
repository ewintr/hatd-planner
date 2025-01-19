package memory

import (
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Memories struct {
	localID  *LocalID
	sync     *Sync
	task     *Task
	schedule *Schedule
}

func New() *Memories {
	return &Memories{
		localID:  NewLocalID(),
		sync:     NewSync(),
		task:     NewTask(),
		schedule: NewSchedule(),
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

func (mems *Memories) Schedule(_ *storage.Tx) storage.Schedule {
	return mems.schedule
}
