package config

type Config struct {
	Name        string `json:"Name"`
	LogPath     string `json:"LogPath"`
	LogMode     string `json:"LogMode"`
	LogEncoding string `json:"LogEncoding"`

	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	MysqlDSN string `json:"MysqlDSN"`

	Redis struct {
		Addrs  []string `json:"Addrs"`
		Pass   string   `json:"Pass"`
		DB     int      `json:"DB"`
		Master string   `json:"Master"`
	} `json:"Redis"`
}
