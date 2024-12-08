package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func AuthAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	var (
		auth string
	)

	if v := r.Header.Get("Authorization"); v != "" {
		auth = v
	} else if v := r.CookiesNamed("authorization"); len(v) > 0 {
		auth = v[0].Value
	} else {
		return nil
	}

	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}

	jwt, err := Parse(auth)
	if err == nil {
		return metadata.Pairs("username", jwt.Payload.Username, "role", jwt.Payload.Role.String())
	}

	return nil
}

func ParseJwtPayload(r *http.Request) (JwtPayload, error) {
	var (
		auth  string
		empty JwtPayload
	)

	if v := r.Header.Get("Authorization"); v != "" {
		auth = v
	} else if v := r.CookiesNamed("authorization"); len(v) > 0 {
		auth = v[0].Value
	} else {
		return empty, nil
	}

	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}

	jwt, err := Parse(auth)
	if err == nil {
		return jwt.Payload, nil
	}

	return empty, err
}

func ForwardResponseOption(ctx context.Context, w http.ResponseWriter, response proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if vs := md.HeaderMD["set-cookie"]; len(vs) > 0 {
		for _, v := range vs {
			w.Header().Add("Set-Cookie", v)
		}
		delete(md.HeaderMD, "set-cookie")
	}
	return nil
}

func UsernameFromContext(ctx context.Context) string {
	if v := metadata.ValueFromIncomingContext(ctx, "username"); len(v) > 0 {
		return v[0]
	}

	if v, ok := ctx.Value("username").(string); ok {
		return v
	}

	return ""
}

func RoleFromContext(ctx context.Context) movie_reservation.User_Role {
	var roleStr string
	if v := metadata.ValueFromIncomingContext(ctx, "role"); len(v) > 0 {
		roleStr = v[0]
	}

	if v, ok := ctx.Value("role").(string); ok {
		roleStr = v
	}

	return movie_reservation.User_Role(movie_reservation.User_Role_value[roleStr])
}

func AssertRole(ctx context.Context, role movie_reservation.User_Role) error {
	if RoleFromContext(ctx) < role {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	return nil
}
