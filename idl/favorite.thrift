include "user.thrift"
include "feed.thrift"

namespace go favorite

struct ActionRequest {
  1:required string token (go.tag = 'json:"token" query:"token"')                                           // 凭证token
  2:required i64 video_id (go.tag = 'json:"video_id" query:"video_id"')                                     // 视频id
  3:required i8 action_type (vt.in = "1", vt.in = "2" go.tag = 'json:"action_type" query:"action_type"')    // 操作类型（1-点赞，2-取消点赞）
}

struct ActionResponse {
  1:required i32 status_code                                                // 状态值
  2:optional string status_msg                                              // 状态信息
}

struct FavoriteRequest {
  1:required i64 user_id (go.tag = 'json:"user_id" query:"user_id"')        // 用户id
  2:required string token (go.tag = 'json:"token" query:"token"')           // 凭证token
}

struct FavoriteResponse {
  1:required i32 status_code                                                // 状态值
  2:optional string status_msg                                              // 状态信息
  3:required list<feed.Video> video_list                                    // 用户点赞视频列表
}

service FavoriteService{
  ActionResponse ActionFavorite(1: ActionRequest req)                       // 点赞操作
  FavoriteResponse MGetFavoriteList(1: FavoriteRequest req)                 // 获取关注列表
}