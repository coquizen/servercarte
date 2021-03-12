package account

import (
	"context"
)

// Service describes the expected behavior in creating accounts and establishing roles for users for the purposes of
// ACL and RBAC
type Service interface {
	Create(context.Context, *createRequest) error
	ChangePassword(context.Context, passwordChangeRequest, string) error
	Delete(context.Context, string, string) error
}


// TODO: implement ACL and Role based access