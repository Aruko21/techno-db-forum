package model

import "time"

type Post struct {
	ID            int64   `json:"id,omitempty"`
	Author        string  `json:"author"`
	Created       time.Time  `json:"created,omitempty"`
	Forum         string  `json:"forum,omitempty"`
	IsEdited      bool    `json:"isEdited"`
	Message       string  `json:"message"`
	Parent        int64   `json:"parent"`
	Thread        int32   `json:"thread,omitempty"`
	Path          []int64 `json:"-"`
	Posts         Posts   `json:"posts,omitempty"`
	ParentPointer *Post   `json:"-"`
}

type PostFull struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

type Posts []*Post