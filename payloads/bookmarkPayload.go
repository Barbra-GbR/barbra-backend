package payloads

//The BookmarkPayload
type BookmarkPayload struct {
	SuggestionId string `json:"suggestion_id" bson:"suggestion_id" validate:"hexadecimal,max30" binding:"required"`
}
