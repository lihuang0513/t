package format

type Server struct {
	Host    string `ini:"host"`
	Port    int    `ini:"port"`
	LogDir  string `ini:"log_dir"`
	Debug   bool   `ini:"debug"`
	NodeUrl string `ini:"node_url"`
}

type PlDb struct {
	Host    string `ini:"host"`
	User    string `ini:"user"`
	Pwd     string `ini:"pwd"`
	Port    int    `ini:"port"`
	Charset string `ini:"charset"`
	Name    string `ini:"name"`
	Enable  bool   `ini:"enable"`
}

type PlBackendRedis struct {
	Host   string `ini:"host"`
	Pwd    string `ini:"pwd"`
	Port   int    `ini:"port"`
	Enable bool   `ini:"enable"`
}

type PlRedis struct {
	Host   string `ini:"host"`
	Pwd    string `ini:"pwd"`
	Port   int    `ini:"port"`
	Enable bool   `ini:"enable"`
}

type ThirdParams struct {
	WyTextBusinessId string `ini:"wy_text_businessId"`
	WyImgBusinessId  string `ini:"wy_img_businessId"`
	WySecretId       string `ini:"wy_secretId"`
	WySecretKey      string `ini:"wy_secretKey"`
	AliIpAppcode     string `ini:"ali_ip_appcode"`
}

type Config struct {
	Server         Server         `ini:"server"`
	PlDb           PlDb           `ini:"pl-db"`
	PlBackendRedis PlBackendRedis `ini:"pl-backend-redis"`
	PlRedis        PlRedis        `ini:"pl-redis"`
	ThirdParams    ThirdParams    `ini:"third-params"`
}
