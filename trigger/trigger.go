package trigger

import (
	"Ratte/conf"
	"Ratte/handler"
	"fmt"
	"github.com/InazumaV/Ratte-Interface/panel"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Trigger struct {
	l        *logrus.Entry
	c        *cron.Cron
	h        *handler.Handler
	p        panel.Panel
	remoteId int
	remoteC  *conf.Remote
	hashs    cmap.ConcurrentMap[string, string]
}

func NewTrigger(
	l *logrus.Entry,
	tc *conf.Trigger,
	h *handler.Handler,
	p panel.Panel,
	rm *conf.Remote,
) (*Trigger, error) {
	tr := &Trigger{
		l:       l,
		c:       cron.New(),
		h:       h,
		p:       p,
		remoteC: rm,
	}

	// add pull node cron task
	_, err := tr.addCronHandle(tc.PullNodeCron, tr.pullNodeHandle)
	if err != nil {
		return nil, err
	}

	// add pull user cron task
	_, err = tr.addCronHandle(tc.PullUserCron, tr.pullUserHandle)
	if err != nil {
		return nil, err
	}

	// add report user cron task
	_, err = tr.addCronHandle(tc.ReportUserCron, tr.reportUserHandle)
	if err != nil {
		return nil, err
	}

	// add renew cert cron task
	_, err = tr.addCronHandle(tc.RenewCertCron, tr.renewCertCron)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (t *Trigger) Start() error {
	r := t.remoteC
	rsp := t.p.AddRemote(&panel.AddRemoteParams{
		Baseurl:  r.APIHost,
		NodeId:   r.NodeID,
		NodeType: r.NodeType,
		Timeout:  r.Timeout,
	})
	if rsp.Err != nil {
		return rsp.Err
	}
	t.remoteId = rsp.RemoteId
	t.pullNodeHandle()
	t.pullUserHandle()
	t.c.Start()
	return nil
}

func (t *Trigger) Close() error {
	t.c.Stop()
	err := t.p.DelRemote(t.remoteId)
	if err != nil {
		return fmt.Errorf("del remote err: %w", err)
	}
	return nil
}
