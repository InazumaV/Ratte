package trigger

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func (t *Trigger) addCronHandle(cron any, job cron.FuncJob) (cron.EntryID, error) {
	switch cron.(type) {
	case string:
		return t.c.AddJob(cron.(string), job)
	case int:
		return t.c.Schedule(newSchedule(cron.(int)), job), nil
	default:
		return 0, fmt.Errorf("unknown cron type: %T", cron)
	}
}

func (t *Trigger) hashEqualsOrStore(name, hash string) bool {
	if h, ok := t.hashs.Get(name); ok {
		if h == hash {
			return true
		}
		t.hashs.Set(name, hash)
	} else {
		t.hashs.Set(name, hash)
	}
	return false
}

func (t *Trigger) pullNodeHandle() {
	t.l.Info("Run pull node task...")
	defer t.l.Info("Run pull node task done.")
	// get node info
	nn := t.p.GetNodeInfo(t.remoteId)
	if nn.Err != nil {
		t.l.WithError(nn.Err).Error("Get node info failed")
		return
	}
	if t.hashEqualsOrStore("pullNode", nn.GetHash()) {
		t.l.Debug("Node is not changed, skip")
		return
	}

	t.l.Debug("Node is changed, triggering handler...")
	// update node handler
	err := t.h.PullNodeHandle(&nn.NodeInfo)
	if err != nil {
		t.l.WithError(err).Error("Pull node failed")
		return
	}
	// done
	t.l.Debug("trigger handler done.")
}

func (t *Trigger) pullUserHandle() {
	t.l.Info("Run pull user task...")
	defer t.l.Info("Run pull user task done.")
	// get user info
	nu := t.p.GetUserList(t.remoteId)
	if nu.Err != nil {
		t.l.WithError(nu.Err).Error("Get user list failed")
		return
	}
	if t.hashEqualsOrStore("pullUser", nu.GetHash()) {
		t.l.Debug("Node is not changed, skip")
		return
	}

	t.l.Debug("user list is changed, triggering handler...")
	// triggering update user list handler
	err := t.h.PullUserHandle(nu.Users)
	if err != nil {
		t.l.WithError(err).Error("Pull user handle failed")
		return
	}
	// done
	t.l.Debug("trigger handler done.")
}

func (t *Trigger) reportUserHandle() {
	t.l.Info("Run report user task...")
	defer t.l.Info("Run pull user task done.")
	// triggering report user handler
	err := t.h.ReportUserHandle(t.remoteId)
	if err != nil {
		t.l.WithError(err).Error("Report user handle failed")
		return
	}
	// done
}

func (t *Trigger) renewCertCron() {
	t.l.Info("Run renew cert task...")
	defer t.l.Info("Run renew cert task done.")
	// triggering renew cert handler
	err := t.h.RenewCertHandle()
	if err != nil {
		t.l.WithError(err).Error("Renew cert handle failed")
		return
	}
	// done
}
