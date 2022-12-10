package wellknown

import (
	"fmt"
	"net/http"

	"github.com/davecheney/m/m"
)

type NodeInfo struct {
	service *Service
}

func (ni *NodeInfo) Index(rw http.ResponseWriter, r *http.Request) {
	toJSON(rw, map[string]any{
		"links": []map[string]any{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				"href": fmt.Sprintf("https://%s/api/nodeinfo/2.0", r.Host),
			},
		},
	})
}

func (ni *NodeInfo) Show(w http.ResponseWriter, r *http.Request) {
	var instance m.Instance
	if err := ni.service.DB().Where("domain = ?", r.Host).First(&instance).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	toJSON(w, serializeNodeInfo(&instance))
}

func serializeNodeInfo(i *m.Instance) map[string]any {
	return map[string]any{
		"version": "2.0", // https://github.com/jhass/nodeinfo/blob/main/schemas/2.0/schema.json
		"software": map[string]any{
			"name":    "https://github.com/davecheney/m",
			"version": "0.0.0-devel",
		},
		"protocols": []string{
			"activitypub",
		},
		"services": map[string]any{
			"outbound": []any{},
			"inbound":  []any{},
		},
		"usage": map[string]any{
			"users": map[string]any{
				"total":          i.AccountsCount,
				"activeMonth":    0,
				"activeHalfyear": 0,
			},
			"localPosts": i.StatusesCount,
		},
		"openRegistrations": false,
		"metadata":          map[string]any{},
	}
}
