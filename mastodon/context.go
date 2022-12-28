package mastodon

import (
	"net/http"
	"strconv"

	"github.com/davecheney/m/m"
	"github.com/go-chi/chi/v5"
	"github.com/go-json-experiment/json"
)

type Contexts struct {
	service *Service
}

func (c *Contexts) Show(w http.ResponseWriter, r *http.Request) {
	_, err := c.service.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	statusID, _ := strconv.ParseUint(id, 10, 64)

	conv, err := c.service.Service.Conversations().FindConversationByStatusID(statusID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// load conversation statuses
	var statuses []m.Status
	if err := c.service.DB().Where("conversation_id = ?", conv.ID).Preload("Actor").Find(&statuses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ancestors, decendants := thread(statusID, statuses)
	w.Header().Set("Content-Type", "application/json")
	json.MarshalFull(w, map[string]interface{}{
		"ancestors": func() []map[string]interface{} {
			var a []map[string]interface{}
			for _, s := range ancestors {
				a = append(a, serializeStatus(s))
			}
			return a
		}(),
		"descendants": func() []map[string]interface{} {
			var a []map[string]interface{}
			for _, s := range decendants {
				a = append(a, serializeStatus(s))
			}
			return a
		}(),
	})
}

// thread sorts statuses into a tree, it returns the statuses
// preceding id, and statuses following id.
func thread(id uint64, statuses []m.Status) ([]*m.Status, []*m.Status) {
	type link struct {
		parent   *link
		status   *m.Status
		children []*link
	}
	ids := make(map[uint64]*link)
	for i := range statuses {
		ids[statuses[i].ID] = &link{status: &statuses[i]}
	}

	for _, l := range ids {
		if l.status.InReplyToID != nil {
			parent, ok := ids[*l.status.InReplyToID]
			if ok {
				// watch out for deleted toots
				l.parent = parent
				parent.children = append(parent.children, l)
			}
		}
	}

	var ancestors []*m.Status
	var l = ids[id].parent
	for l != nil {
		ancestors = append(ancestors, l.status)
		l = l.parent
	}
	reverse(ancestors)

	var descendants []*m.Status
	var walk func(*link)
	walk = func(l *link) {
		for _, c := range l.children {
			descendants = append(descendants, c.status)
			walk(c)
		}
	}
	walk(ids[id])
	return ancestors, descendants
}

func reverse[T any](a []T) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
