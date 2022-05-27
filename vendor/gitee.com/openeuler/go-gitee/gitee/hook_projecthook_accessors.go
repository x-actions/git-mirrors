package gitee

func (pj *ProjectHook) GetNameSpace() string {
	if pj == nil {
		return ""
	}

	return pj.Namespace
}

func (pj *ProjectHook) GetPath() string {
	if pj == nil {
		return ""
	}

	return pj.Path
}

func (pj *ProjectHook) GetOwner() *UserHook {
	if pj == nil {
		return nil
	}

	return pj.Owner
}

func (pj *ProjectHook) GetOwnerAndRepo() (string, string) {
	if pj == nil {
		return "", ""
	}

	return pj.Namespace, pj.Path
}