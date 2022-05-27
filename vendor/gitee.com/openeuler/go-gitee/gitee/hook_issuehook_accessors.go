package gitee

func (ih *IssueHook) GetUser() *UserHook {
	if ih == nil {
		return nil
	}

	return ih.User
}

func (ih *IssueHook) GetNumber() string {
	if ih == nil {
		return ""
	}

	return ih.Number
}

func (ih *IssueHook) GetState() string {
	if ih == nil {
		return ""
	}

	return ih.State
}
