namespace go user

struct User {
  1:required i64 id                        // 用户id
  2:required string name                   // 用户名称
  3:required string avatar                 // 用户头像
  4:required i64 follow_count              // 关注总数
  5:required i64 follower_count            // 粉丝总数
  6:required bool is_follow                // true-已关注，false-未关注
  7:required string background_image       // 用户个人页顶部大图
  8:required string signature              // 个人简介
  9:required i64 total_favorited           // 获赞数量
  10:required i64 work_count               // 作品数
  11:required i64 favorite_count           // 喜欢数
}

struct CreateUserRequest {
  1:required string username (vt.min_size = "1", vt.max_size = "32" go.tag = 'json:"username" query:"username"')           // 注册用户名，最长32个字符
  2:required string password (vt.min_size = "1", vt.max_size = "32" go.tag = 'json:"password" query:"password"')           // 密码，最长32个字符
}

struct CreateUserResponse {
  1:required i32 status_code               // 状态值
  2:optional string status_msg             // 状态信息
  3:required i64 user_id                   // 用户id
  4:required string token                  // 用户鉴权token
}

struct CheckUserRequest {
  1:string username (vt.min_size = "1", vt.max_size = "32" go.tag = 'json:"username" query:"username"')           // 登录用户名，最长32个字符
  2:string password (vt.min_size = "1", vt.max_size = "32" go.tag = 'json:"password" query:"password"')           // 密码，最长32个字符
}

struct CheckUserResponse {
  1:required i32 status_code               // 状态值
  2:optional string status_msg             // 状态信息
  3:required i64 user_id                   // 用户id
  4:required string token                  // 用户鉴权token
}

struct MGetUserRequest {
  1:required i64 user_id      (go.tag = 'json:"user_id" query:"user_id"')       // 用户id
  2:required string token     (go.tag = 'json:"token" query:"token"')           // 用户鉴权token
}

struct MGetUserResponse {
  1:required i32 status_code               // 状态值
  2:optional string status_msg             // 状态信息
  3:required User user                     // 用户信息
}

service UserService {
  CreateUserResponse CreateUser(1: CreateUserRequest req)   // 创建用户信息
  MGetUserResponse MGetUser(1: MGetUserRequest req)         // 获取用户信息
  CheckUserResponse CheckUser(1: CheckUserRequest req)      // 验证用户
}