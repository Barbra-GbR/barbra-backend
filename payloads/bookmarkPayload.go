package payloads

type BookmarkPayload struct {
	SuggestionId string `json:"suggestion_id" bson:"suggestion_id" validate:"hexadecimal,max30" binding:"required"`
}
