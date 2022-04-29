package mode

type ActiveMode int

const (
	UserListMode ActiveMode = iota
	UserDetailsMode
	RoleListMode
	RoleDetailsMode
)
