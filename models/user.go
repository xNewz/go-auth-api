package models

type User struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Username   string `json:"username,omitempty" bson:"username,omitempty"`
	Password   string `json:"password,omitempty" bson:"password,omitempty"`
	Email      string `json:"email,omitempty" bson:"email,omitempty"`
	Role       string `json:"role,omitempty" bson:"role,omitempty"`
	ProfileURL string `json:"profile_url,omitempty" bson:"profile_url,omitempty"`
	Provider   string `json:"provider,omitempty" bson:"provider,omitempty"`
	ProviderID string `json:"provider_id,omitempty" bson:"provider_id,omitempty"`
}
