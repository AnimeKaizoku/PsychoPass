package tokenHandlers

import sv "github.com/AnimeKaizoku/PsychoPass/sibyl/core/sibylValues"

type ChangePermResult struct {
	PreviousPerm sv.UserPermission `json:"previous_perm"`
	CurrentPerm  sv.UserPermission `json:"current_perm"`
}
