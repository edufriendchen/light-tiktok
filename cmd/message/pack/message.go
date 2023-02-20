package pack

import (
	"time"

	"github.com/edufriendchen/light-tiktok/cmd/message/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/message"
)

// User pack user info
func Message(u *dal.Message) *message.Message {
	if u == nil {
		return nil
	}
	return &message.Message{
		Id:         int64(u.ID),
		ToUserId:   u.ToUserId,
		FromUserId: u.FromUserId,
		Content:    u.Content,
		CreateTime: FormatTime(u.CreatedAt),
	}
}

// Messages pack list of message info
func Messages(us []*dal.Message) []*message.Message {
	messages := make([]*message.Message, 0)
	for _, u := range us {
		if temp := Message(u); temp != nil {
			messages = append(messages, temp)
		}
	}
	return messages
}

func FormatTime(t time.Time) *string {
	s := t.Format("2006-01-02 15:04:05")
	return &s
}
