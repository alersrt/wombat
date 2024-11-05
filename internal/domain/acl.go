package domain

type Acl struct {
	AuthorId  string `json:"author_id"`
	IsAllowed bool   `json:"is_allowed"`
}
