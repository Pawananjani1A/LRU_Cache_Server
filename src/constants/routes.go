/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package constants

// RouterBasePrefix BasePrefix added before all router paths
type RouterBasePrefix string

const (
	BasePrefix       RouterBasePrefix = "lruCache"
	ModuleBasePrefix RouterBasePrefix = "poc"
)

// RVersion ConfigVersion management for routes
type RVersion string

const (
	RV1 RVersion = "v1"
)

// RoutesGroup route groups
type RoutesGroup string

// RoutesPath route sub-paths
type RoutesPath string

const (
	RGDocs            RoutesGroup = "docs"
	RGInternalGroup   RoutesGroup = "testRoutes"
	RGCheckRouterPath RoutesPath  = "routerCheck"
)

const (
	RGGroup RoutesGroup = "groups"
	RGCache RoutesGroup = "cache"
)

const (
	RPGetMemberGroups RoutesPath = ":memberId"
	RPUpdateGroup     RoutesPath = ":groupId/update"
	RPGet             RoutesPath = "get"
	RPSet             RoutesPath = "set"
)
