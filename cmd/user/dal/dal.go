package dal

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/gofrs/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/crypto/bcrypt"
)

const (
	queryUserCYPHER       = `MATCH (u:User) WHERE id(u) = $user_id RETURN u`
	queryToUserCYPHER     = `MATCH (n:User) WHERE id(n) = $toUserId MATCH (n1:User) WHERE id(n1) = $userId RETURN n, CASE WHEN (n1)-[:follow]->(n) THEN true ELSE false END AS result`
	getUserPassWordCYPHER = `MATCH (n:User {username: $username}) RETURN n.password AS ps, id(n) AS i LIMIT 1`
	queryUserConutCYPHER  = `MATCH (n:User {username: $username}) RETURN count(*) AS count LIMIT 1`
	createUserCYPHER      = `CREATE (n:User { username: $username, password: $password, nickname: $nickname,avatar: $avatar, follow_count: $follow_count,
							follower_count: $follower_count, background_image: $background_image, signature: $signature, total_favorited: $total_favorited,
							work_count: $work_count, favorite_count: $favorite_count }) RETURN n`
)

func CheckUser(ctx context.Context, session neo4j.SessionWithContext, req *user.CheckUserRequest) (int64, error) {
	userid, err := neo4j.ExecuteRead[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, getUserPassWordCYPHER, map[string]any{
			"username": req.Username,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, errno.AuthorizationFailedErr
		}
		if !BcryptCheck(req.Password, record.Values[0].(string)) {
			return 0, errno.AuthorizationFailedErr
		}
		return record.Values[1].(int64), nil
	})
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func CreateUser(ctx context.Context, session neo4j.SessionWithContext, req *user.CreateUserRequest) (int64, error) {
	req.Password = BcryptHash(req.Password)
	nickname := DefaultNickName()
	avatar := DefaultAvatar()
	userid, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (userid int64, err error) {
		result, err := tx.Run(ctx, queryUserConutCYPHER, map[string]any{
			"username": req.Username,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		count, found := record.Get("count")
		op, _ := count.(int64)
		if !found {
			return 0, fmt.Errorf("could not find column")
		}
		if op != 0 {
			return 0, errno.UserAlreadyExistErr
		}
		result, err = tx.Run(ctx, createUserCYPHER, map[string]any{
			"username":         req.Username,
			"password":         req.Password,
			"nickname":         nickname,
			"avatar":           avatar,
			"follow_count":     0,
			"follower_count":   0,
			"background_image": "",
			"signature":        "",
			"total_favorited":  0,
			"work_count":       0,
			"favorite_count":   0,
		})
		if err != nil {
			return 0, err
		}
		record, err = result.Single(ctx)
		if err != nil {
			return 0, err
		}
		rawPerson, found := record.Get("n")
		if !found {
			return 0, fmt.Errorf("could not find column")
		}
		itemNode, ok := rawPerson.(neo4j.Node)
		if !ok {
			return 0, fmt.Errorf("expected result to be a map but was %T", rawPerson)
		}
		userid = itemNode.GetId()
		return userid, nil
	})
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func GetUser(ctx context.Context, session neo4j.SessionWithContext, req *user.MGetUserRequest, user_id int64) (*user.User, error) {
	var userItem user.User
	user, err := neo4j.ExecuteRead[*user.User](ctx, session, func(tx neo4j.ManagedTransaction) (*user.User, error) {
		result, err := tx.Run(ctx, queryUserCYPHER,
			map[string]any{
				"user_id": user_id,
			})
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		rawItemNode, found := record.Get("u")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		itemNode := rawItemNode.(neo4j.Node)
		if userItem.Name, err = neo4j.GetProperty[string](itemNode, "nickname"); err != nil {
			return nil, err
		}
		if userItem.Avatar, err = neo4j.GetProperty[string](itemNode, "avatar"); err != nil {
			return nil, err
		}
		follow_count, err := neo4j.GetProperty[int64](itemNode, "follow_count")
		if err != nil {
			return nil, err
		}
		follower_count, err := neo4j.GetProperty[int64](itemNode, "follower_count")
		if err != nil {
			return nil, err
		}
		userItem.BackgroundImage, err = neo4j.GetProperty[string](itemNode, "background_image")
		if err != nil {
			return nil, err
		}
		userItem.Signature, err = neo4j.GetProperty[string](itemNode, "signature")
		if err != nil {
			return nil, err
		}
		total_favorited, err := neo4j.GetProperty[int64](itemNode, "total_favorited")
		if err != nil {
			return nil, err
		}
		work_count, err := neo4j.GetProperty[int64](itemNode, "work_count")
		if err != nil {
			return nil, err
		}
		favorite_count, err := neo4j.GetProperty[int64](itemNode, "favorite_count")
		if err != nil {
			return nil, err
		}
		userItem.FollowCount = follow_count
		userItem.FollowerCount = follower_count
		userItem.TotalFavorited = total_favorited
		userItem.WorkCount = work_count
		userItem.FavoriteCount = favorite_count
		return &userItem, nil
	})
	return user, err
}

func GetToUser(ctx context.Context, session neo4j.SessionWithContext, req *user.MGetUserRequest, user_id int64) (*user.User, error) {
	var userItem *user.User
	user, err := neo4j.ExecuteRead[*user.User](ctx, session, func(tx neo4j.ManagedTransaction) (*user.User, error) {
		result, err := tx.Run(ctx, queryUserCYPHER,
			map[string]any{
				"userId": user_id,
			})
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		rawItemNode, found := record.Get("n")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		itemNode := rawItemNode.(neo4j.Node)
		userItem.Id = itemNode.GetId()
		if userItem.Name, err = neo4j.GetProperty[string](itemNode, "nickname"); err != nil {
			return nil, err
		}
		if userItem.Avatar, err = neo4j.GetProperty[string](itemNode, "avatar"); err != nil {
			return nil, err
		}
		follow_count, err := neo4j.GetProperty[int64](itemNode, "follow_count")
		if err != nil {
			return nil, err
		}
		follower_count, err := neo4j.GetProperty[int64](itemNode, "follower_count")
		if err != nil {
			return nil, err
		}
		userItem.BackgroundImage, err = neo4j.GetProperty[string](itemNode, "background_image")
		if err != nil {
			return nil, err
		}
		userItem.Signature, err = neo4j.GetProperty[string](itemNode, "signature")
		if err != nil {
			return nil, err
		}
		total_favorited, err := neo4j.GetProperty[int64](itemNode, "total_favorited")
		if err != nil {
			return nil, err
		}
		work_count, err := neo4j.GetProperty[int64](itemNode, "work_count")
		if err != nil {
			return nil, err
		}
		favorite_count, err := neo4j.GetProperty[int64](itemNode, "favorite_count")
		if err != nil {
			return nil, err
		}
		userItem.FollowCount = follow_count
		userItem.FollowerCount = follower_count
		userItem.TotalFavorited = total_favorited
		userItem.WorkCount = work_count
		userItem.FavoriteCount = favorite_count
		return userItem, nil
	})
	return user, err
}

func GetUserInNode(node dbtype.Node) (user *user.User, err error) {
	user.Id = node.GetId()
	if user.Name, err = neo4j.GetProperty[string](node, "nickname"); err != nil {
		return nil, err
	}
	if user.Avatar, err = neo4j.GetProperty[string](node, "avatar"); err != nil {
		return nil, err
	}
	follow_count, err := neo4j.GetProperty[int64](node, "follow_count")
	if err != nil {
		return nil, err
	}
	follower_count, err := neo4j.GetProperty[int64](node, "follower_count")
	if err != nil {
		return nil, err
	}
	user.BackgroundImage, err = neo4j.GetProperty[string](node, "background_image")
	if err != nil {
		return nil, err
	}
	user.Signature, err = neo4j.GetProperty[string](node, "signature")
	if err != nil {
		return nil, err
	}
	total_favorited, err := neo4j.GetProperty[int64](node, "total_favorited")
	if err != nil {
		return nil, err
	}
	work_count, err := neo4j.GetProperty[int64](node, "work_count")
	if err != nil {
		return nil, err
	}
	favorite_count, err := neo4j.GetProperty[int64](node, "favorite_count")
	if err != nil {
		return nil, err
	}
	user.FollowCount = follow_count
	user.FollowerCount = follower_count
	user.TotalFavorited = total_favorited
	user.WorkCount = work_count
	user.FavoriteCount = favorite_count
	return user, nil
}

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// DefaultNickName 获得一个默认的昵称
func DefaultNickName() string {
	uid, _ := uuid.NewV4()
	return "用户" + uid.String()[0:8]
}

// DefaultBackground 获得一个的背景
func DefaultBackground() string {
	AvatarList := [...]string{
		"https://img1.baidu.com/it/u=1459539381,1684299919&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=264b8cfbd62ce23ee0d0a557091cc72d",
		"https://img0.baidu.com/it/u=2155033989,3634097964&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=a6973ebee4d25a8611a4efd711ed52a3",
		"https://img1.baidu.com/it/u=256830766,4270545878&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=967010975a9f1ae6ef80e46fca09f713",
		"https://img2.baidu.com/it/u=324241668,3161137356&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=7087110ff8179531e16664396e414809",
		"https://img1.baidu.com/it/u=4102089746,1733025287&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=e33dd1e45318b5402ae9317742c91c10",
		"https://img2.baidu.com/it/u=1659974989,4260768333&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=a60b66a45ff9585df06c491ae799ea74",
		"https://img2.baidu.com/it/u=346152429,3164401706&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=b25ae1f4c313ef0010e95975bcc21c5d",
		"https://img1.baidu.com/it/u=1900416729,2440027599&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=85dfaac59b2febf3fa995735e00190b1",
	}
	rand.Seed(time.Now().UnixNano())
	return AvatarList[rand.Intn(8)]
}

// DefaultAvatar 获得一个默认的头像
func DefaultAvatar() string {
	AvatarList := [...]string{
		"https://img1.baidu.com/it/u=1459539381,1684299919&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=264b8cfbd62ce23ee0d0a557091cc72d",
		"https://img0.baidu.com/it/u=2155033989,3634097964&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=a6973ebee4d25a8611a4efd711ed52a3",
		"https://img1.baidu.com/it/u=256830766,4270545878&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=967010975a9f1ae6ef80e46fca09f713",
		"https://img2.baidu.com/it/u=324241668,3161137356&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=7087110ff8179531e16664396e414809",
		"https://img1.baidu.com/it/u=4102089746,1733025287&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=e33dd1e45318b5402ae9317742c91c10",
		"https://img2.baidu.com/it/u=1659974989,4260768333&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=a60b66a45ff9585df06c491ae799ea74",
		"https://img2.baidu.com/it/u=346152429,3164401706&fm=253&app=138&size=w931&n=0&f=JPG&fmt=auto?sec=1675530000&t=b25ae1f4c313ef0010e95975bcc21c5d",
		"https://img1.baidu.com/it/u=1900416729,2440027599&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1675530000&t=85dfaac59b2febf3fa995735e00190b1",
	}
	rand.Seed(time.Now().UnixNano())
	return AvatarList[rand.Intn(8)]
}
