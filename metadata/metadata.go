package metadata

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	authorization              = "authorization"
)

type Metadata struct {
	IP        string
	Bearer    string
	UserAgent string
}

func ExtractMetadataFromContext(ctx context.Context) *Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("Aucune métadonnée trouvée dans le contexte")
		return nil
	}

	m := &Metadata{}

	// Extraire l'adresse IP
	if ip := md.Get(xForwardedForHeader); len(ip) > 0 {
		m.IP = ip[0]
	}

	// Extraire le User-Agent
	if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
		m.UserAgent = userAgents[0]
	}

	// Extraire le Bearer token
	if authHeaders := md.Get(authorization); len(authHeaders) > 0 {
		for _, authHeader := range authHeaders {
			if strings.HasPrefix(authHeader, "Bearer ") {
				m.Bearer = strings.TrimPrefix(authHeader, "Bearer ")
				break
			}
		}
	}

	return m
}
