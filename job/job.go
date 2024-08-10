package job

import (
	"github.com/DiamondDrakeVentures/patchver/task"
	"github.com/gammazero/deque"
)

type Job struct {
	id    string
	tasks *deque.Deque[task.Task]
}

func (j Job) ID() string {
	return j.id
}

func (j *Job) NextTask() task.Task {
	if j.tasks.Len() > 0 {
		return j.tasks.PopFront()
	}

	return nil
}

func NewJob(id string, tasks []task.Task) *Job {
	q := deque.New[task.Task]()

	for _, task := range tasks {
		q.PushBack(task)
	}

	return &Job{
		id:    id,
		tasks: q,
	}
}
