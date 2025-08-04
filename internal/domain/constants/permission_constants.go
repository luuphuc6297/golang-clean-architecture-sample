package constants

// ContextKey represents a key type used for context values to avoid collisions.
type ContextKey string

const (
	// ResourceUser represents the user resource type for permissions.
	ResourceUser = "user"
	// ResourceProduct represents the product resource type for permissions.
	ResourceProduct = "product"

	// ActionCreate represents the create action
	ActionCreate = "create"
	// ActionRead represents the read action
	ActionRead = "read"
	// ActionUpdate represents the update action
	ActionUpdate = "update"
	// ActionDelete represents the delete action
	ActionDelete = "delete"
	// ActionList represents the list action
	ActionList = "list"

	// PermissionUserCreate defines permission for creating new users
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
