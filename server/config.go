package server

type GRPCConfigServer interface {
	Host() string
	Port() int
}

type YAMLGRPCConfigServer struct {
	ValuePort int    `yaml:"port"`
	ValueHost string `yaml:"host"`
}

func (c YAMLGRPCConfigServer) Port() int {
	return c.ValuePort
}

func (c YAMLGRPCConfigServer) Host() string {
	return c.ValueHost
}
