package user

type JiraUser interface {
	//GetLogin() string
	//GetPassword() string
	//GetJiraUserKey() (string, error)
}

type user struct {
	root        string
	login       string
	password    string
	jiraUserKey string
}

func NewUser(root, login, password string) JiraUser {
	return &user{root, login, password, ""}
}

func (u *user) GetRoot() string {
	return u.root
}

func (u *user) GetLogin() string {
	return u.login
}

func (u *user) GetPassword() string {
	return u.password
}

func (u *user) GetJiraUserKey() string {
	return u.jiraUserKey
}

func (u *user) SetJiraUserKey(jiraUserKey string) {
	u.jiraUserKey = jiraUserKey
}
