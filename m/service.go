package m

import (
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// Service represents the m web service.
type Service struct {
	db *gorm.DB
}

// NewService returns a new Service.
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// authenticate authenticates the bearer token attached to the request and, if
// successful, returns the account associated with the token.
func (s *Service) authenticate(r *http.Request) (*Account, error) {
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	var token Token
	if err := s.db.Where("access_token = ?", bearer).Joins("Account").First(&token).Error; err != nil {
		return nil, err
	}
	return token.Account, nil
}

func (s *Service) API() *API {
	return &API{
		service: s,
	}
}

// NodeInfo returns a NodeInfo REST resource.
func (s *Service) NodeInfo() *NodeInfo {
	return &NodeInfo{
		db: s.db,
	}
}

func (s *Service) Users() *Users {
	return &Users{
		db:      s.db,
		service: s,
	}
}

func (s *Service) WellKnown() *WellKnown {
	return &WellKnown{
		db: s.db,
	}
}

func (s *Service) Statuses() *statuses {
	return &statuses{
		db:      s.db,
		service: s,
	}
}

func (s *Service) conversations() *conversations {
	return &conversations{
		db:      s.db,
		service: s,
	}
}

func (s *Service) Accounts() *accounts {
	return &accounts{
		db:      s.db,
		service: s,
	}
}

func (s *Service) instances() *instances {
	return &instances{
		db: s.db,
	}
}

func (s *Service) ActivityPub() *ActivityPub {
	return &ActivityPub{
		service: s,
	}
}

func (a *ActivityPub) Inboxes() *Inboxes {
	return &Inboxes{
		service: a.service,
	}
}

func (a *ActivityPub) Outboxes() *Outbox {
	return &Outbox{
		service: a.service,
	}
}

// API rerpesents the root of a Mastodon capable REST API.
type API struct {
	service *Service
}

func (a *API) Accounts() *Accounts {
	return &Accounts{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Applications() *Applications {
	return &Applications{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Contexts() *Contexts {
	return &Contexts{
		service: a.service,
	}
}

func (a *API) Conversations() *Conversations {
	return &Conversations{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Emojis() *Emojis {
	return &Emojis{
		db: a.service.db,
	}
}

func (a *API) Favourites() *Favourites {
	return &Favourites{
		service: a.service,
	}
}

func (a *API) Filters() *Filters {
	return &Filters{
		db: a.service.db,
	}
}

func (a *API) Instances() *Instances {
	return &Instances{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Lists() *Lists {
	return &Lists{
		db: a.service.db,
	}
}

func (a *API) Markers() *Markers {
	return &Markers{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Notifications() *Notifications {
	return &Notifications{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Relationships() *Relationships {
	return &Relationships{
		service: a.service,
	}
}

func (a *API) Statuses() *Statuses {
	return &Statuses{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Timelines() *Timelines {
	return &Timelines{
		db:      a.service.db,
		service: a.service,
	}
}

func (a *API) Search() *Search {
	return &Search{
		db:      a.service.db,
		service: a.service,
	}
}
