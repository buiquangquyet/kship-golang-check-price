package repo

func NewSettingShopRepo(base *baseRepo) *SettingShopRepo {
	return &SettingShopRepo{
		base,
	}
}

type SettingShopRepo struct {
	*baseRepo
}
