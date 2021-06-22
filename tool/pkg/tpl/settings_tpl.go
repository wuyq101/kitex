package tpl

// SettingsTpl is the template for generating setting.py.
var SettingsTpl = `PRODUCT="{{.Product}}"
SUBSYSTEM="{{.Subsystem}}"
MODULE="{{.Module}}"
APP_TYPE="binary"
`
