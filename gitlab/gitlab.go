package gitlab

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chenlujjj/hook/weixin"
	"github.com/gin-gonic/gin"
)

const (
	MRActionOpen       = "open"
	MRActionClose      = "close"
	MRActionReopen     = "reopen"
	MRActionUpdate     = "update"
	MRActionApproved   = "approved"
	MRActionUnapproved = "unapproved"
	MRActionMerge      = "merge"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type Project struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	WebURL            string      `json:"web_url"`
	AvatarURL         interface{} `json:"avatar_url"`
	GitSSHURL         string      `json:"git_ssh_url"`
	GitHTTPURL        string      `json:"git_http_url"`
	Namespace         string      `json:"namespace"`
	VisibilityLevel   int         `json:"visibility_level"`
	PathWithNamespace string      `json:"path_with_namespace"`
	DefaultBranch     string      `json:"default_branch"`
	Homepage          string      `json:"homepage"`
	URL               string      `json:"url"`
	SSHURL            string      `json:"ssh_url"`
	HTTPURL           string      `json:"http_url"`
}

type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url"`
	Author    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

type Label struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Color       string    `json:"color"`
	ProjectID   int       `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Template    bool      `json:"template"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	GroupID     int       `json:"group_id"`
}
type Attributes struct {
	ID              int         `json:"id"`
	TargetBranch    string      `json:"target_branch"`
	SourceBranch    string      `json:"source_branch"`
	SourceProjectID int         `json:"source_project_id"`
	AuthorID        int         `json:"author_id"`
	AssigneeID      int         `json:"assignee_id"`
	Title           string      `json:"title"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	MilestoneID     interface{} `json:"milestone_id"`
	State           string      `json:"state"`
	MergeStatus     string      `json:"merge_status"`
	TargetProjectID int         `json:"target_project_id"`
	Iid             int         `json:"iid"`
	Description     string      `json:"description"`
	Source          Project     `json:"source"` // without project id
	Target          Project     `json:"target"`
	LastCommit      Commit      `json:"last_commit"`
	WorkInProgress  bool        `json:"work_in_progress"`
	URL             string      `json:"url"`
	Action          string      `json:"action"`
	Assignee        struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"assignee"`
}

// See https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#merge-request-events.
type MergeRequestEvent struct {
	ObjectKind       string     `json:"object_kind"`
	User             User       `json:"user"`
	Project          Project    `json:"project"`
	Repository       Repository `json:"repository"`
	ObjectAttributes Attributes `json:"object_attributes"`
	Labels           []Label    `json:"labels"`
	Changes          struct {
		UpdatedByID struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"updated_by_id"`
		UpdatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"updated_at"`
		Labels struct {
			Previous []Label `json:"previous"`
			Current  []Label `json:"current"`
		} `json:"labels"`
	} `json:"changes"`
}

func NewMRHandler(wc *weixin.WechatClient) func(*gin.Context) {
	return func(c *gin.Context) {
		var mrEvent MergeRequestEvent
		if err := c.BindJSON(&mrEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := handleMergeRequestEvent(mrEvent, wc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func handleMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	switch event.ObjectAttributes.Action {
	case MRActionOpen:
		return handleOpenMergeRequestEvent(event, wc)
	case MRActionApproved:
		return handleApprovedMergeRequestEvent(event, wc)
	case MRActionMerge:
		return handleMergeMergeRequestEvent(event, wc)
	case MRActionClose:
		return handleCloseMergeRequestEvent(event, wc)
	case MRActionUpdate:
		return handleUpdateMergeRequestEvent(event, wc)
	default:
		// 暂不处理其他类型的MR事件
		return nil
	}
}

var openMergeRequestMessageFmt = `✨ %s发起了新的Merge Request: %s
项目：%s
源分支：%s
目标分支：%s
分配给：%s
描述：%s
链接：%s
`

var approvedMergeRequestMessageFmt = `🍻 %s批准了Merge Request: %s
项目：%s
源分支：%s
目标分支：%s
分配给：%s
链接：%s
`

var mergeMergeRequestMessageFmt = `🚀 %s合并了Merge Request: %s
项目：%s
源分支：%s
目标分支：%s
分配给：%s
链接：%s
`

func handleOpenMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	user := event.User.Name
	title := event.ObjectAttributes.Title
	desc := event.ObjectAttributes.Description
	project := event.Project.Name
	sourceBranch := event.ObjectAttributes.SourceBranch
	targetBranch := event.ObjectAttributes.TargetBranch
	url := event.ObjectAttributes.URL
	assignee := event.ObjectAttributes.Assignee.Name
	return wc.SendText(fmt.Sprintf(openMergeRequestMessageFmt, user, title, project, sourceBranch, targetBranch, assignee, desc, url))
}

func handleApprovedMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	user := event.User.Name
	title := event.ObjectAttributes.Title
	project := event.Project.Name
	sourceBranch := event.ObjectAttributes.SourceBranch
	targetBranch := event.ObjectAttributes.TargetBranch
	url := event.ObjectAttributes.URL
	assignee := event.ObjectAttributes.Assignee.Name
	return wc.SendText(fmt.Sprintf(approvedMergeRequestMessageFmt, user, title, project, sourceBranch, targetBranch, assignee, url))

}

func handleMergeMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	user := event.User.Name
	title := event.ObjectAttributes.Title
	project := event.Project.Name
	sourceBranch := event.ObjectAttributes.SourceBranch
	targetBranch := event.ObjectAttributes.TargetBranch
	url := event.ObjectAttributes.URL
	assignee := event.ObjectAttributes.Assignee.Name
	return wc.SendText(fmt.Sprintf(mergeMergeRequestMessageFmt, user, title, project, sourceBranch, targetBranch, assignee, url))

}

func handleUpdateMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	// TODO
	return nil
}

func handleCloseMergeRequestEvent(event MergeRequestEvent, wc *weixin.WechatClient) error {
	// TODO
	return nil
}
