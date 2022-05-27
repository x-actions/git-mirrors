package gitee

type NoteEvent struct {
	Action        *string          `json:"action,omitempty"`
	Comment       *NoteHook        `json:"comment,omitempty"`
	Repository    *ProjectHook     `json:"repository,omitempty"`
	Project       *ProjectHook     `json:"project,omitempty"`
	Author        *UserHook        `json:"author,omitempty"`
	Sender        *UserHook        `json:"sender,omitempty"`
	URL           *string          `json:"url,omitempty"`
	Note          *string          `json:"note,omitempty"`
	NoteableType  *string          `json:"noteable_type,omitempty"`
	NoteableID    int64            `json:"noteable_id,omitempty"`
	Title         *string          `json:"title,omitempty"`
	PerIID        *string          `json:"per_iid,omitempty"`
	ShortCommitID *string          `json:"short_commit_id,omitempty"`
	Enterprise    *EnterpriseHook  `json:"enterprise,omitempty"`
	PullRequest   *PullRequestHook `json:"pull_request,omitempty"`
	Issue         *IssueHook       `json:"issue,omitempty"`
	HookName      *string          `json:"hook_name,omitempty"`
	Password      *string          `json:"password,omitempty"`
}

type PushEvent struct {
	Ref                *string         `json:"ref,omitempty"`
	Before             *string         `json:"before,omitempty"`
	After              *string         `json:"after,omitempty"`
	TotalCommitsCount  int64           `json:"total_commits_count,omitempty"`
	CommitsMoreThanTen *bool           `json:"commits_more_than_ten,omitempty"`
	Created            *bool           `json:"created,omitempty"`
	Deleted            *bool           `json:"deleted,omitempty"`
	Compare            *string         `json:"compare,omitempty"`
	Commits            []CommitHook    `json:"commits,omitempty"`
	HeadCommit         *CommitHook     `json:"head_commit,omitempty"`
	Repository         *ProjectHook    `json:"repository,omitempty"`
	Project            *ProjectHook    `json:"project,omitempty"`
	UserID             int64           `json:"user_id,omitempty"`
	UserName           *string         `json:"user_name,omitempty"`
	User               *UserHook       `json:"user,omitempty"`
	Pusher             *UserHook       `json:"pusher,omitempty"`
	Sender             *UserHook       `json:"sender,omitempty"`
	Enterprise         *EnterpriseHook `json:"enterprise,omitempty"`
	HookName           *string         `json:"hook_name,omitempty"`
	Password           *string         `json:"password,omitempty"`
}

type IssueEvent struct {
	Action      *string         `json:"action,omitempty"`
	Issue       *IssueHook      `json:"issue,omitempty"`
	Repository  *ProjectHook    `json:"repository,omitempty"`
	Project     *ProjectHook    `json:"project,omitempty"`
	Sender      *UserHook       `json:"sender,omitempty"`
	TargetUser  *UserHook       `json:"target_user,omitempty"`
	User        *UserHook       `json:"user,omitempty"`
	Assignee    *UserHook       `json:"assignee,omitempty"`
	UpdatedBy   *UserHook       `json:"updated_by,omitempty"`
	IID         string          `json:"iid,omitempty"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	State       *string         `json:"state,omitempty"`
	Milestone   *string         `json:"milestone,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Enterprise  *EnterpriseHook `json:"enterprise,omitempty"`
	HookName    *string         `json:"hook_name,omitempty"`
	Password    *string         `json:"password,omitempty"`
}

type RepoInfo struct {
	Project    *ProjectHook `json:"project,omitempty"`
	Repository *ProjectHook `json:"repository,omitempty"`
}

type PullRequestEvent struct {
	Action         *string          `json:"action,omitempty"`
	ActionDesc     *string          `json:"action_desc,omitempty"`
	PullRequest    *PullRequestHook `json:"pull_request,omitempty"`
	Number         int64            `json:"number,omitempty"`
	IID            int64            `json:"iid,omitempty"`
	Title          *string          `json:"title,omitempty"`
	Body           *string          `json:"body,omitempty"`
	State          *string          `json:"state,omitempty"`
	MergeStatus    *string          `json:"merge_status,omitempty"`
	MergeCommitSha *string          `json:"merge_commit_sha,omitempty"`
	URL            *string          `json:"url,omitempty"`
	SourceBranch   *string          `json:"source_branch,omitempty"`
	SourceRepo     *RepoInfo        `json:"source_repo,omitempty"`
	TargetBranch   *string          `json:"target_branch,omitempty"`
	TargetRepo     *RepoInfo        `json:"target_repo,omitempty"`
	Project        *ProjectHook     `json:"project,omitempty"`
	Repository     *ProjectHook     `json:"repository,omitempty"`
	Author         *UserHook        `json:"author,omitempty"`
	UpdatedBy      *UserHook        `json:"updated_by,omitempty"`
	Sender         *UserHook        `json:"sender,omitempty"`
	TargetUser     *UserHook        `json:"target_user,omitempty"`
	Enterprise     *EnterpriseHook  `json:"enterprise,omitempty"`
	HookName       *string          `json:"hook_name,omitempty"`
	Password       *string          `json:"password,omitempty"`
}

type TagPushEvent struct {
	Action *string `json:"action,omitempty"`
}
