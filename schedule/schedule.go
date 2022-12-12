package schedule

import (
	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	AddJob(taskName, spec string, f func())
	Run()
}

type Schedule struct {
	cron *cron.Cron
	jobs map[string]cron.EntryID
}

func (m *Schedule) Jobs() map[string]cron.EntryID {
	return m.jobs
}

func (m *Schedule) AddJob(taskName, spec string, f func()) {
	id, _ := m.cron.AddJob(spec, cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&BaseJob{f}))
	m.jobs[taskName] = id
}

func (m *Schedule) RemoveJob(taskName string) {
	if id, ok := m.jobs[taskName]; ok {
		m.cron.Remove(id)
	}
}

func (m *Schedule) Run() {
	m.cron.Run()
}

type BaseJob struct {
	F func()
}

func (p *BaseJob) Run() {
	p.F()
}

func New() *Schedule {
	return &Schedule{
		cron: cron.New(),
		jobs: make(map[string]cron.EntryID, 0),
	}
}
