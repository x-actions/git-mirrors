package gitee

func (pe *PushEvent) GetAction() string {
	if pe == nil || pe.Ref == nil {
		return ""
	}

	return *pe.Ref
}

func (pe *PushEvent) GetBefore() string {
	if pe == nil || pe.Before == nil {
		return ""
	}

	return *pe.Before
}

func (pe *PushEvent) GetAfter() string {
	if pe == nil || pe.After == nil {
		return ""
	}

	return *pe.After
}

func (pe *PushEvent) GetCommitsMoreThanTen() bool {
	if pe == nil || pe.CommitsMoreThanTen == nil {
		return false
	}

	return *pe.CommitsMoreThanTen
}

func (pe *PushEvent) GetCreated() bool {
	if pe == nil || pe.Created == nil {
		return false
	}

	return *pe.Created
}

func (pe *PushEvent) GetDeleted() bool {
	if pe == nil || pe.Deleted == nil {
		return false
	}

	return *pe.Deleted
}

func (pe *PushEvent) GetCompare() string {
	if pe == nil || pe.Compare == nil {
		return ""
	}

	return *pe.Compare
}

func (pe *PushEvent) GetHeadCommit() *CommitHook {
	if pe == nil {
		return nil
	}

	return pe.HeadCommit
}

func (pe *PushEvent) GetRepository() *ProjectHook {
	if pe == nil {
		return nil
	}

	return pe.Repository
}

func (pe *PushEvent) GetProject() *ProjectHook {
	if pe == nil {
		return nil
	}

	return pe.Project
}

func (pe *PushEvent) GetUserName() string {
	if pe == nil || pe.UserName == nil {
		return ""
	}

	return *pe.UserName
}

func (pe *PushEvent) GetUser() *UserHook {
	if pe == nil {
		return nil
	}

	return pe.User
}

func (pe *PushEvent) GetPusher() *UserHook {
	if pe == nil {
		return nil
	}

	return pe.Pusher
}

func (pe *PushEvent) GetSender() *UserHook {
	if pe == nil {
		return nil
	}

	return pe.Sender
}

func (pe *PushEvent) GetEnterprise() *EnterpriseHook {
	if pe == nil {
		return nil
	}

	return pe.Enterprise
}

func (pe *PushEvent) GetHookName() string {
	if pe == nil || pe.HookName == nil {
		return ""
	}

	return *pe.HookName
}

func (pe *PushEvent) GetPassword() string {
	if pe == nil || pe.Password == nil {
		return ""
	}

	return *pe.Password
}
