package family

type Family struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	BaseCurrency string `json:"baseCurrency"`
	OwnerUserId  int64  `json:"ownerUserId"`
	Role         string `json:"role"`
}

type CreateFamilyRequest struct {
	Name         string `json:"name"`
	BaseCurrency string `json:"baseCurrency"`
}
