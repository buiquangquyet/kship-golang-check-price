package repo

func NewService(base *baseRepo) *Service {
	return &Service{
		base,
	}
}

type Service struct {
	*baseRepo
}
