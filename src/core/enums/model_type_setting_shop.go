package enums

type ModelTypeSettingShop string

const (
	ModelTypeClientDisableShop       ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ClientDisableShop"
	ModelTypeClientSettingShop       ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ServiceSettingShop"
	ModelTypeServiceExtraSettingShop ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ServiceExtraSettingShop"
)

func (m ModelTypeSettingShop) ToString() string {
	return string(m)
}
