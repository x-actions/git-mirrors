package gitee

func (c *CommitHook) GetAuthor() *UserHook {
	if c == nil {
		return nil
	}

	return c.Author
}

func (c *CommitHook) GetCommitter() *UserHook {
	if c == nil {
		return nil
	}

	return c.Committer
}