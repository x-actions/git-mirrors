package gitee

func (l *LabelHook) GetName() string {
	if l == nil {
		return ""
	}

	return l.Name
}

func (l *LabelHook) GetColor() string {
	if l == nil {
		return ""
	}

	return l.Color
}
