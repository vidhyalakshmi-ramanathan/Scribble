package enums

type Roles string

const (
	Owner  Roles = "owner"
	Admin  Roles = "admin"
	Editor Roles = "editor"
	Viewer Roles = "viewer"
)

type Invites string

const (
	InvitePending  Invites = "pending"
	InviteExpires  Invites = "expires"
	InviteAccepted Invites = "accepted"
	InviteRejected Invites = "rejected"
)
