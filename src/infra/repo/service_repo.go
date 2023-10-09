package repo

func NewServiceRepo(base *baseRepo) *ServiceRepo {
	return &ServiceRepo{
		base,
	}
}

type ServiceRepo struct {
	*baseRepo
}
