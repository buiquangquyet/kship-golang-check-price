package repo

func NewWardRepo(base *baseRepo) *WardRepo {
	return &WardRepo{
		base,
	}
}

type WardRepo struct {
	*baseRepo
}
