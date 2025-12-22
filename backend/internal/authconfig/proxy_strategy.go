package authconfig

import (
	"context"
	"strings"

	"connectrpc.com/connect"
)

type ProxyStrategy struct {
	baseStrategy
	GroupHeader     string
	UserIdHeader    string
	SuperUserGroups []string
	AdminGroups     []string
	GroupSeparator  *string
}

func NewProxyStrategy(
	groupHeader, userIdHeader string,
	superUserGroups, adminGroups []string,
	groupSeparator *string,
) AuthStrategy {
	return &ProxyStrategy{
		GroupHeader:     groupHeader,
		UserIdHeader:    userIdHeader,
		SuperUserGroups: superUserGroups,
		AdminGroups:     adminGroups,
		GroupSeparator:  groupSeparator,
	}
}

func (s *ProxyStrategy) Type() AuthStrategyType { return AuthStrategyProxy }

func (s *ProxyStrategy) Authenticate(ctx context.Context, req connect.AnyRequest) (context.Context, error) {
	separator := "|"
	if s.GroupSeparator != nil {
		separator = *s.GroupSeparator
	}
	rawGroups := strings.Split(req.Header().Get(s.GroupHeader), separator)
	contextValue := authContext{}
	if includesGroup(s.AdminGroups, rawGroups) {
		contextValue.UserPermission = Admin
	} else if includesGroup(s.SuperUserGroups, rawGroups) {
		contextValue.UserPermission = SuperUser
	} else {
		contextValue.UserPermission = User
	}
	userId := req.Header().Get(s.UserIdHeader)
	if userId == "" {
		contextValue.UserId = nil
	} else {
		contextValue.UserId = &userId
	}
	ctx = s.setContext(ctx, contextValue.UserPermission, contextValue.UserId)
	return ctx, nil
}

func includesGroup(rawGroups, Groups []string) bool {
	set := make(map[string]struct{})
	result := make([]string, 0)
	for _, v := range rawGroups {
		set[v] = struct{}{}
	}
	for _, v := range Groups {
		if _, ok := set[v]; ok {
			result = append(result, v)
			delete(set, v)
		}
	}
	return len(result) > 0
}
