package enums

type ModelTypeSettingShop string

const (
	ModelTypeClientDisableShop       ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ClientDisableShop"
	ModelTypeClientSettingShop       ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ClientSettingShop"
	ModelTypeServiceExtraSettingShop ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ServiceExtraSettingShop"
	ModelTypeServiceExtraDisableShop ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ServiceDisableExtraSettingShop"
	ModelTypeServiceExtraSettingUser ModelTypeSettingShop = "App\\Core\\WidgetSetting\\ServiceExtraSettingUser"
)

func (m ModelTypeSettingShop) ToString() string {
	return string(m)
}
