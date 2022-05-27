package gitee

func (nh *NoteHook) GetBody() string {
	if nh == nil {
		return ""
	}

	return nh.Body
}
