package gitee

func (ne *NoteEvent) GetAction() string {
	if ne == nil || ne.Action == nil {
		return ""
	}

	return *ne.Action
}

func (ne *NoteEvent) GetComment() *NoteHook {
	if ne == nil {
		return nil
	}

	return ne.Comment
}

func (ne *NoteEvent) GetRepository() *ProjectHook {
	if ne == nil {
		return nil
	}

	return ne.Repository
}

func (ne *NoteEvent) GetProject() *ProjectHook {
	if ne == nil {
		return nil
	}

	return ne.Project
}

func (ne *NoteEvent) GetAuthor() *UserHook {
	if ne == nil {
		return nil
	}

	return ne.Author
}

func (ne *NoteEvent) GetSender() *UserHook {
	if ne == nil {
		return nil
	}

	return ne.Sender
}

func (ne *NoteEvent) GetURL() string {
	if ne == nil || ne.URL == nil {
		return ""
	}

	return *ne.URL
}

func (ne *NoteEvent) GetNote() string {
	if ne == nil || ne.Note == nil {
		return ""
	}

	return *ne.Note
}

func (ne *NoteEvent) GetNoteableType() string {
	if ne == nil || ne.NoteableType == nil {
		return ""
	}

	return *ne.NoteableType
}

func (ne *NoteEvent) GetTitle() string {
	if ne == nil || ne.Title == nil {
		return ""
	}

	return *ne.Title
}

func (ne *NoteEvent) GetPerIID() string {
	if ne == nil || ne.PerIID == nil {
		return ""
	}

	return *ne.PerIID
}

func (ne *NoteEvent) GetShortCommitID() string {
	if ne == nil || ne.ShortCommitID == nil {
		return ""
	}

	return *ne.ShortCommitID
}

func (ne *NoteEvent) GetEnterprise() *EnterpriseHook {
	if ne == nil {
		return nil
	}

	return ne.Enterprise
}

func (ne *NoteEvent) GetPullRequest() *PullRequestHook {
	if ne == nil {
		return nil
	}

	return ne.PullRequest
}

func (ne *NoteEvent) GetIssue() *IssueHook {
	if ne == nil {
		return nil
	}

	return ne.Issue
}

func (ne *NoteEvent) GetHookName() string {
	if ne == nil || ne.HookName == nil {
		return ""
	}

	return *ne.HookName
}

func (ne *NoteEvent) GetPassword() string {
	if ne == nil || ne.Password == nil {
		return ""
	}

	return *ne.Password
}
