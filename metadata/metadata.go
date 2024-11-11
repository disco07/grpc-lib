package metadata

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	if p, ok := peer.FromContext(ctx); ok {
		m.IP, _, _ = net.SplitHostPort(p.Addr.String())
	}

	// Extraire le User-Agent
	if userAgents, ok := md["user-agent"]; ok && len(userAgents) > 0 {
		m.UserAgent = userAgents[0]
	}

	// Extraire le Bearer token
	if authHeaders, ok := md["authorization"]; ok && len(authHeaders) > 0 {
		for _, authHeader := range authHeaders {
			if strings.HasPrefix(authHeader, "Bearer ") {
				m.Bearer = strings.TrimPrefix(authHeader, "Bearer ")
				break
			}
		}
	}

	return m
}
