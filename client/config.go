package client

type GRPCConfigClient interface {
	Port() int
}

type YAMLGRPCConfigClient struct {
	ValuePort int `yaml:"port"`
}

func (c YAMLGRPCConfigClient) Port() int {
	return c.ValuePort
}
