package repo

func NewClientRepo(base *baseRepo) *ClientRepo {
	return &ClientRepo{
		base,
	}
}

type ClientRepo struct {
	*baseRepo
}
