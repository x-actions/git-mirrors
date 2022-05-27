package gitee

func (u *UserHook) GetLogin() string {
	if u == nil {
		return ""
	}

	return u.Login
}