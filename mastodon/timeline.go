package mastodon

import (
	"fmt"
	"net/http"

	"github.com/davecheney/pub/internal/algorithms"
	"github.com/davecheney/pub/internal/models"
	"github.com/davecheney/pub/internal/to"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Timelines struct {
	service *Service
}

func (t *Timelines) Home(w http.ResponseWriter, r *http.Request) {
	user, err := t.service.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var followingIDs []int64
	if err := t.service.db.Model(&models.Relationship{ActorID: user.Actor.ID}).Where("following = true").Pluck("target_id", &followingIDs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followingIDs = append(followingIDs, int64(user.ID))

	var statuses []*models.Status
	scope := t.service.db.Scopes(models.PaginateStatuses(r)).Where("(actor_id IN (?) AND in_reply_to_actor_id is null) or (actor_id in (?) and in_reply_to_actor_id IN (?))", followingIDs, followingIDs, followingIDs)
	scope = scope.Joins("Actor").Preload("Reblog").Preload("Reblog.Actor").Preload("Attachments").Preload("Reaction", "actor_id = ?", user.Actor.ID)
	if err := scope.Find(&statuses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(statuses) > 0 {
		w.Header().Set("Link", fmt.Sprintf("<https://%s/api/v1/timelines/home?max_id=%d>; rel=\"next\", <https://%s/api/v1/timelines/home?min_id=%d>; rel=\"prev\"", r.Host, statuses[len(statuses)-1].ID, r.Host, statuses[0].ID))
	}
	to.JSON(w, algorithms.Map(statuses, serialiseStatus))
}

func TimelinesPublic(w http.ResponseWriter, r *http.Request) {
	db, _ := r.Context().Value("DB").(*gorm.DB)
	var statuses []*models.Status
	scope := db.Scopes(models.PaginateStatuses(r)).Where("visibility = ? and reblog_id is null and in_reply_to_id is null", "public")
	switch r.URL.Query().Get("local") {
	case "true":
		scope = scope.Joins("Actor").Where("Actor.domain = ?", r.Host)
	default:
		scope = scope.Joins("Actor")
	}
	scope = scope.Preload("Attachments")
	if err := scope.Find(&statuses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(statuses) > 0 {
		w.Header().Set("Link", fmt.Sprintf("<https://%s/api/v1/timelines/public?max_id=%d>; rel=\"next\", <https://%s/api/v1/timelines/public?min_id=%d>; rel=\"prev\"", r.Host, statuses[len(statuses)-1].ID, r.Host, statuses[0].ID))
	}
	to.JSON(w, algorithms.Map(statuses, serialiseStatus))
}

func (t *Timelines) ListShow(w http.ResponseWriter, r *http.Request) {
	user, err := t.service.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var listMembers []int64
	if err := t.service.db.Model(&models.AccountListMember{}).Where("account_list_id = ? ", chi.URLParam(r, "id")).Pluck("member_id", &listMembers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var statuses []*models.Status
	scope := t.service.db.Scopes(models.PaginateStatuses(r)).Where("(actor_id IN (?) AND in_reply_to_actor_id is null) or (actor_id in (?) and in_reply_to_actor_id IN (?))", listMembers, listMembers, listMembers)
	scope = scope.Joins("Actor").Preload("Reblog").Preload("Reblog.Actor").Preload("Attachments").Preload("Reaction", "actor_id = ?", user.Actor.ID)
	if err := scope.Find(&statuses).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if len(statuses) > 0 {
	// 	w.Header().Set("Link", fmt.Sprintf("<https://%s/api/v1/timelines/home?max_id=%d>; rel=\"next\", <https://%s/api/v1/timelines/home?min_id=%d>; rel=\"prev\"", r.Host, statuses[len(statuses)-1].ID, r.Host, statuses[0].ID))
	// }
	to.JSON(w, algorithms.Map(statuses, serialiseStatus))
}
