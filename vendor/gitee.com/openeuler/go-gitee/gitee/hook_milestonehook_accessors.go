package gitee

func (m *MilestoneHook) GetTitle() string {
	if m == nil {
		return ""
	}

	return m.Title
}

func (m *MilestoneHook) GetState() string {
	if m == nil {
		return ""
	}

	return m.State
}
