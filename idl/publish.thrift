include "feed.thrift"

namespace go publish

struct ActionRequest {
    1:binary data (go.tag = 'json:"data" from:"data"')                                   // 视频数据
    2:required string token (go.tag = 'json:"token" from:"token"')                       // 用户凭证                        // 用户鉴权token
    3:required string title (go.tag = 'json:"title" from:"title"')                       // 视频标题
}

struct ActionResponse {
    1:required i32 status_code               // 状态值
    2:optional string status_msg             // 状态信息
}

struct PublishRequest {
    1:required i64 user_id (go.tag = 'json:"user_id" query:"user_id"')                // 用户id
    2:required string token (go.tag = 'json:"token" query:"token"')                   // 凭证token
}

struct PublishResponse {
    1:required i32 status_code               // 状态值
    2:optional string status_msg             // 状态信息
    3:required list<feed.Video> video_list   // 用户发布的视频列表
}

service PublishService{
  ActionResponse ActionPulish(1: ActionRequest req)               // 视频上传操作
  PublishResponse MGetPublishList(1: PublishRequest req)          // 获取发布的视频列表
}