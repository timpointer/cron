package cron

import (
	"sync"

	c "gopkg.in/robfig/cron.v2"
)

// Cron wrap cron.Cron and add key remove function
type Cron struct {
	*c.Cron
	sm sync.Mutex
	m  map[string]c.EntryID
}

// EntryID wrap cron.EntryID
type EntryID struct {
	c.EntryID
}

// Entry wrap cron.Entry
type Entry struct {
	c.Entry
}

// Job wrap cron.Job and identify a job by ID
type Job interface {
	c.Job
	ID() string
}

// New create a new cron wrapper
func New() *Cron {
	return &Cron{
		Cron: c.New(),
		m:    make(map[string]c.EntryID),
	}
}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (n *Cron) AddJob(spec string, cmd Job) error {
	n.sm.Lock()
	defer n.sm.Unlock()
	id, err := n.Cron.AddJob(spec, cmd)
	n.m[cmd.ID()] = id
	return err
}

// Clean clean all cron job
func (n *Cron) Clean() {
	n.sm.Lock()
	defer n.sm.Unlock()
	for _, id := range n.m {
		n.Cron.Remove(id)
	}
}

func (n *Cron) SetJob(ID string, spec string, cmd Job) error {
	n.sm.Lock()
	defer n.sm.Unlock()
	id, ok := n.m[ID]
	if ok == true {
		n.remove(ID)
	}
	id, err := n.Cron.AddJob(spec, cmd)
	n.m[cmd.ID()] = id
	return err
}

// Remove an entry from being run in the future.
func (n *Cron) Remove(ID string) {
	n.sm.Lock()
	defer n.sm.Unlock()
	n.remove(ID)
}

func (n *Cron) remove(ID string) {
	id, ok := n.m[ID]
	if ok == true {
		n.Cron.Remove(id)
		delete(n.m, ID)
	}
}

// Entries wrap cron.Entries function
func (n *Cron) Entries() []Entry {
	var list []Entry
	for _, v := range n.Cron.Entries() {
		list = append(list, Entry{v})
	}
	return list
}

// Entries wrap cron.Entries function
func (n *Cron) EntryMap() map[string]Entry {
	var list []Entry
	for _, v := range n.Cron.Entries() {
		list = append(list, Entry{v})
	}
	mp := make(map[string]Entry)
	n.sm.Lock()
	for key, val := range n.m {
		for _, j := range list {
			if val == j.ID {
				mp[key] = j
			}
		}
	}
	n.sm.Unlock()

	return mp
}
