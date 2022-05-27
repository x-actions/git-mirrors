package gitee

func (eh *EnterpriseHook) GetName() string {
	if eh == nil {
		return ""
	}

	return eh.Name
}

func (eh *EnterpriseHook) GetUrl() string {
	if eh == nil {
		return ""
	}

	return eh.Url
}