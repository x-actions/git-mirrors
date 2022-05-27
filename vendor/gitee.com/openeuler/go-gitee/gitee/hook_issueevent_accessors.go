package gitee

func (ie *IssueEvent) GetAction() string {
	if ie == nil || ie.Action == nil {
		return ""
	}

	return *ie.Action
}

func (ie *IssueEvent) GetIssue() *IssueHook {
	if ie == nil {
		return nil
	}

	return ie.Issue
}

func (ie *IssueEvent) GetRepository() *ProjectHook {
	if ie == nil {
		return nil
	}

	return ie.Repository
}

func (ie *IssueEvent) GetProject() *ProjectHook {
	if ie == nil {
		return nil
	}

	return ie.Project
}

func (ie *IssueEvent) GetSender() *UserHook {
	if ie == nil {
		return nil
	}

	return ie.Sender
}

func (ie *IssueEvent) GetTargetUser() *UserHook {
	if ie == nil {
		return nil
	}

	return ie.TargetUser
}

func (ie *IssueEvent) GetUser() *UserHook {
	if ie == nil {
		return nil
	}

	return ie.User
}

func (ie *IssueEvent) GetAssignee() *UserHook {
	if ie == nil {
		return nil
	}

	return ie.Assignee
}

func (ie *IssueEvent) GetUpdatedBy() *UserHook {
	if ie == nil {
		return nil
	}

	return ie.UpdatedBy
}

func (ie *IssueEvent) GetTitle() string {
	if ie == nil || ie.Title == nil {
		return ""
	}

	return *ie.Title
}

func (ie *IssueEvent) GetDescription() string {
	if ie == nil || ie.Description == nil {
		return ""
	}

	return *ie.Description
}

func (ie *IssueEvent) GetState() string {
	if ie == nil || ie.State == nil {
		return ""
	}

	return *ie.State
}

func (ie *IssueEvent) GetMilestone() string {
	if ie == nil || ie.Milestone == nil {
		return ""
	}

	return *ie.Milestone
}

func (ie *IssueEvent) GetURL() string {
	if ie == nil || ie.URL == nil {
		return ""
	}

	return *ie.URL
}

func (ie *IssueEvent) GetEnterprise() *EnterpriseHook {
	if ie == nil {
		return nil
	}

	return ie.Enterprise
}

func (ie *IssueEvent) GetHookName() string {
	if ie == nil || ie.HookName == nil {
		return ""
	}

	return *ie.HookName
}

func (ie *IssueEvent) GetPassword() string {
	if ie == nil || ie.Password == nil {
		return ""
	}

	return *ie.Password
}
