package dal

import (
	"context"

	"github.com/edufriendchen/light-tiktok/pkg/global"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ToUserId   int64  `gorm:"type:varchar(32);not null" json:"to_user_id"`
	FromUserId int64  `gorm:"type:varchar(32);not null" json:"from_user_id"`
	Content    string `gorm:"type:varchar(256);not null" json:"content"`
}

func (Message) TableName() string {
	return "message"
}

// CreateMessage
func CreateMessage(ctx context.Context, messages []*Message) error {
	return global.DB.WithContext(ctx).Create(messages).Error
}

// MGetChatMsg
func MGetChatMsg(ctx context.Context, to_user_id int64, from_user_id int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := global.DB.WithContext(ctx).Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)", to_user_id, from_user_id, from_user_id, to_user_id).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
