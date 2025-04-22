package handler

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/common/slices"
)

func compareUserList(old, new []panel.UserInfo) (deleted []string, added []panel.UserInfo) {
	tmp := map[string]struct{}{}
	tmp2 := map[string]struct{}{}
	for i := range old {
		tmp[old[i].GetHashOrKey()] = struct{}{}
	}
	l := len(tmp)
	for i := range new {
		e := new[i].GetHashOrKey()
		tmp[e] = struct{}{}
		tmp2[e] = struct{}{}
		if l != len(tmp) {
			added = append(added, new[i])
			l++
		}
	}
	tmp = nil
	l = len(tmp2)
	for i := range old {
		tmp2[old[i].GetHashOrKey()] = struct{}{}
		if l != len(tmp2) {
			deleted = append(deleted, old[i].Name)
			l++
		}
	}
	return deleted, added
}

func (h *Handler) PullUserHandle(users []panel.UserInfo) error {
	if !h.nodeAdded.Load() {
		return nil
	}
	del, add := compareUserList(h.userList, users)
	cas := slices.RangeToNew[panel.UserInfo, core.UserInfo](add, func(_ int, v panel.UserInfo) core.UserInfo {
		return core.UserInfo(v.UserInfo)
	})
	err := h.c.AddUsers(&core.AddUsersParams{
		NodeName: h.nodeName,
		Users:    cas,
	})
	if err != nil {
		return fmt.Errorf("add user error: %w", err)
	}
	h.l.Infof("Added %d users", len(users))
	err = h.c.DelUsers(&core.DelUsersParams{
		NodeName: h.nodeName,
		Users:    del,
	})
	if err != nil {
		return fmt.Errorf("del user error: %w", err)
	}
	h.l.Infof("Deleted %d users", len(users))
	h.userList = users
	return nil
}

func (h *Handler) ReportUserHandle(id int) error {
	if !h.nodeAdded.Load() {
		return nil
	}
	var err error
	req := &core.GetUserTrafficParams{NodeName: h.nodeName}
	var users []panel.UserTrafficInfo
	slices.Range(h.userList, func(_ int, v panel.UserInfo) bool {
		req.Username = v.Name
		rsp := h.c.GetUserTraffic(req)
		if rsp.Err != nil {
			err = rsp.Err
			return true
		}
		if rsp.Up == 0 && rsp.Down == 0 {
			return false
		}
		users = append(users, panel.UserTrafficInfo{
			Id:       v.Id,
			Name:     v.Name,
			Upload:   rsp.Up,
			Download: rsp.Down,
		})
		return false
	})
	if err != nil {
		return fmt.Errorf("get user traffic error: %w", err)
	}
	if len(users) == 0 {
		return nil
	}
	err = h.p.ReportUserTraffic(&panel.ReportUserTrafficParams{
		Id:    id,
		Users: users,
	})
	if err != nil {
		return fmt.Errorf("report user error: %w", err)
	}
	return nil
}
