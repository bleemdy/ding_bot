package schedule

import (
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	cron  *cron.Cron
	tasks map[string]cron.EntryID
}

type BaseTask struct {
	F func()
}

func (p *BaseTask) Run() {
	p.F()
}

func (m *Schedule) AddJob(taskName, spec string, f func()) {
	id, _ := m.cron.AddJob(spec, cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&BaseTask{f}))
	m.tasks[taskName] = id
}

func (m Schedule) Run() {
	m.cron.Run()
}

func New() *Schedule {
	return &Schedule{
		cron:  cron.New(),
		tasks: make(map[string]cron.EntryID, 0),
	}
}
