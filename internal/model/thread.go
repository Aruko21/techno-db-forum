package model

type Thread struct {
	ID      int32  `json:"id,omitempty"`
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int32  `json:"votes"`
}

type NewThread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Threads []Thread
