include "user.thrift"

namespace go feed

struct Video {
    1:required i64 id                           // 状态值
    2:required user.User author                 // 视频作者信息
    3:required string play_url                  // 视频播放地址
    4:required string cover_url                 // 视频封面地址
    5:required i64 favorite_count               // 视频的点赞总数
    6:required i64 comment_count                // 视频的评论总数
    7:required bool is_favorite                 // true-已点赞，false-未点赞
    8:required string title                     // 视频标题
}

struct FeedRequest {
    1:i64 latest_time (go.tag = 'json:"latest_time" query:"latest_time"')         // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
    2:string token (go.tag = 'json:"token" query:"token"')                        // 用户登录状态下设置
}

struct FeedResponse {
    1:required i32 status_code                 // 状态值
    2:optional string status_msg               // 状态信息
    3:required list<Video> video_list          // 视频列表
    4:optional i64 next_time                   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}


service FeedService{
    FeedResponse MGetFeedList(1: FeedRequest req)       // 获取按投稿时间倒序的视频列表
}