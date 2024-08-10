package executor

import (
	"fmt"

	"github.com/DiamondDrakeVentures/patchver/common"
	"github.com/DiamondDrakeVentures/patchver/job"
	"github.com/DiamondDrakeVentures/patchver/task"
	"github.com/gammazero/deque"
)

type Executor struct {
	logger common.Logger

	current *job.Job
	q       *deque.Deque[*job.Job]

	jobCount  int
	taskCount int
}

func (e *Executor) Register(job *job.Job) {
	e.q.PushBack(job)
}

func (e *Executor) Execute() error {
	for tsk := e.nextTask(); tsk != nil; tsk = e.nextTask() {
		e.logger.Printf("Executing %s (%s)\n", tsk.Name(), tsk.Type())
		err := e.executeTask(tsk)
		if err != nil {
			return err
		}
		e.logger.Printf("Done executing %s\n", tsk.Name())
	}

	return nil
}

func (e *Executor) executeTask(tsk task.Task) error {
	prefix := e.logger.Prefix()

	e.logger.SetPrefix(fmt.Sprintf("%s/%s ", e.current.ID(), tsk.ID()))
	defer e.logger.SetPrefix(prefix)

	return tsk.Execute(e.logger)
}

func (e *Executor) nextTask() task.Task {
	if e.current == nil {
		e.nextJob()

		if e.current != nil {
			return e.nextTask()
		} else {
			return nil
		}
	}

	if tsk := e.current.NextTask(); tsk != nil {
		e.taskCount++
		return tsk
	} else {
		e.nextJob()
		return e.nextTask()
	}
}

func (e *Executor) nextJob() {
	if e.q.Len() > 0 {
		e.current = e.q.PopFront()

		e.logger.SetPrefix(fmt.Sprintf("%s ", e.current.ID()))

		e.taskCount = 0
		e.jobCount++
	} else {
		e.current = nil
	}
}

func New(logger common.Logger) *Executor {
	return &Executor{
		logger: logger,

		q: deque.New[*job.Job](),
	}
}
