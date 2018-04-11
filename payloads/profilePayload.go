package payloads

//The ProfilePayload
type ProfilePayload struct {
	Email      string `json:"email"       bson:"email"       validate:"omitempty,email,lowercase"`
	GivenName  string `json:"given_name"  bson:"given_name"  validate:"omitempty,alphaunicode,min=1,max=50"`
	FamilyName string `json:"family_name" bson:"family_name" validate:"omitempty,alphaunicode,min=1,max=50"`
	PictureURL string `json:"picture"     bson:"picture"     validate:"omitempty,url"`
	Nickname   string `json:"nickname"    bson:"nickname"    validate:"omitempty,alphanumunicode,min=2,max=50"`
}
