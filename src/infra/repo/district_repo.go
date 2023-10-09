package repo

func NewDistrictRepo(base *baseRepo) *DistrictRepo {
	return &DistrictRepo{
		base,
	}
}

type DistrictRepo struct {
	*baseRepo
}
