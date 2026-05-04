package enumTransfer

import "seckill/internal/userSvr/kitex_gen/usersvr"

func EnumToRoleString(kenum usersvr.UserRole) string {
	return kenum.String()
}

func RoleStringToEnum(role string) usersvr.UserRole {
	if role == usersvr.UserRole_ADMIN.String() {
		return usersvr.UserRole_ADMIN
	}

	return usersvr.UserRole_SIMPLE_USER
}
