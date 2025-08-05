package constants

type ContextKey string

const (
	ResourceUser    = "user"
	ResourceProduct = "product"

	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"

	PermissionUserCreate = "user:create"
	PermissionUserRead   = "user:read"
	PermissionUserUpdate = "user:update"
	PermissionUserDelete = "user:delete"
	PermissionUserList   = "user:list"

	PermissionProductCreate = "product:create"
	PermissionProductRead   = "product:read"
	PermissionProductUpdate = "product:update"
	PermissionProductDelete = "product:delete"
	PermissionProductList   = "product:list"

	PolicyEffectAllow = "allow"
	PolicyEffectDeny  = "deny"

	ContextUserID    = ContextKey("user_id")
	ContextUserRole  = ContextKey("user_role")
	ContextUserEmail = ContextKey("user_email")
	ContextClientIP  = ContextKey("client_ip")
)
