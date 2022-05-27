package gitee

func (b *BranchHook) GetLabel() string {
	if b == nil {
		return ""
	}

	return b.Label
}

func (b *BranchHook) GetRef() string {
	if b == nil {
		return ""
	}

	return b.Ref
}

func (b *BranchHook) GetSha() string {
	if b == nil {
		return ""
	}

	return b.Sha
}

func (b *BranchHook) GetUser() *UserHook {
	if b == nil {
		return nil
	}

	return b.User
}

func (b *BranchHook) GetRepo() *ProjectHook {
	if b == nil {
		return nil
	}

	return b.Repo
}