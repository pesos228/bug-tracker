package dto

type IdTokenClaims struct {
	Email       string `json:"email"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}
