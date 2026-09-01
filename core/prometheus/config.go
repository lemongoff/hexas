package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:",optional"`
	Port int
	Path string
}
