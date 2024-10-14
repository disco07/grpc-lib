package server

type GRPCConfigServer interface {
	Port() int
}

type YAMLGRPCConfigServer struct {
	ValuePort int `yaml:"port"`
}

func (c YAMLGRPCConfigServer) Port() int {
	return c.ValuePort
}
