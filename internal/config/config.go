package config

type Config struct {
	Name string `yaml:"Name"`
	Port int    `yaml:"Port"`
	Log  struct {
		FileName string `yaml:"FileName"`
		Mode     string `yaml:"Mode"`
		Path     string `yaml:"Path"`
		Level    string `yaml:"Level"`
		Compress bool   `yaml:"Compress"`
		KeepDays int    `yaml:"KeepDays"`
	} `yaml:"Log"`

	UploadConf struct {
		MaxImageSize    int64  `yaml:"MaxImageSize"`
		MaxVideoSize    int64  `yaml:"MaxVideoSize"`
		MaxAudioSize    int64  `yaml:"MaxAudioSize"`
		MaxOtherSize    int64  `yaml:"MaxOtherSize"`
		PublicStorePath string `yaml:"PublicStorePath"`
		ServerURL       string `yaml:"ServerURL"`
	} `yaml:"UploadConf"`
}
