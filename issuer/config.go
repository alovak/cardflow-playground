package issuer

type Config struct {
	HTTPAddr    string
	ISO8583Addr string
}

func DefaultConfig() *Config {
	return &Config{
		HTTPAddr:    "localhost:9090",
		ISO8583Addr: "localhost:8583",
	}
}
