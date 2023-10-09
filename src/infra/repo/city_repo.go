package repo

func NewCityRepo(base *baseRepo) *CityRepo {
	return &CityRepo{
		base,
	}
}

type CityRepo struct {
	*baseRepo
}
