package starter

import "github.com/hashicorp/go-multierror"

type Starter interface {
	Start() error
	Close() error
}

type Group struct {
	Starters []Starter
}

func NewGroup(starters ...Starter) *Group {
	return &Group{
		Starters: starters,
	}
}

func (g *Group) Add(starter Starter) {
	g.Starters = append(g.Starters, starter)
}

func (g *Group) Clear(len int) {
	g.Starters = nil
}

func (g *Group) Start() error {
	for _, starter := range g.Starters {
		if err := starter.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) Close() error {
	var errs error
	for i := len(g.Starters) - 1; i >= 0; i-- {
		if err := g.Starters[i].Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
