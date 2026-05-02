package models

type Config struct {
	Hertz HertzConfigEntry
	Nacos NacosConfigEntry
}

type HertzConfigEntry struct {
	ListenAddr string
	ListenPort string
}

type NacosConfigEntry struct {
	UserName  string
	Password  string
	Namespace string
	Host      string
	Port      string
}
