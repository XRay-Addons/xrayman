package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

type UserID = int

type User struct {
	ID        UserID
	Name      string
	VlessUUID string
}

func (u User) VlessEmail() string {
	return fmt.Sprintf("%d-%s", u.ID, u.Name)
}

func ParseVlessEmail(email string) (id UserID, name string, err error) {
	defer func() {
		if err != nil {
			err = xerr.WrapWithf(err, "email: %s", email)
		}
	}()

	i := strings.IndexByte(email, '-')
	if i <= 0 || i == len(email)-1 {
		return 0, "", xerr.New("invalid format")
	}

	if id, err = strconv.Atoi(email[:i]); err != nil {
		return 0, "", xerr.WrapWithStack(err)
	}

	return UserID(id), email[i+1:], nil
}
