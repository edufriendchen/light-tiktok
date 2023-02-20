include "user.thrift"

namespace go comment

struct Comment {
    1:required i64 id                          // 状态值
    2:required user.User user                  // 用户信息
    3:required string content                  // 评论内容
    4:required string  create_date             // 评论发布日期，格式 mm-dd
}

struct ActionRequest {
    1:required string token (go.tag = 'json:"token" query:"token"')                                          // 用户鉴权token
    2:required i64 video_id (go.tag = 'json:"video_id" query:"video_id"')                                    // 视频id
    3:required i8 action_type (vt.in = "1", vt.in = "2" go.tag = 'json:"action_type" query:"action_type"')   // 操作类型（1-发布评论，2-删除评论）
    4:string comment_text (go.tag = 'json:"comment_text" query:"comment_text"')                              // 用户填写的评论内容，在action_type=1的时候使用
    5:i64 comment_id (go.tag = 'json:"comment_id" query:"comment_id"')                                       // 要删除的评论id，在action_type=2的时候使用
}

struct ActionResponse {
    1:required i32 status_code               // 状态值
    2:optional string status_msg             // 状态信息
    3:optional Comment comment               // 评论成功返回评论内容，不需要重新拉取整个列表
}

struct CommentRequest {
    1:required i64 user_id (go.tag = 'json:"user_id" query:"user_id"')                  // 用户鉴权token
    2:required i64 video_id (go.tag = 'json:"video_id" query:"video_id"')               // 视频id
}

struct CommentResponse {
    1:required i32 status_code               // 状态值
    2:optional string status_msg             // 状态信息
    3:required list<Comment> comment_list    // 评论成功返回评论内容，不需要重新拉取整个列表
}

service CommentService {
    ActionResponse ActionComment(1: ActionRequest req)              // 评论操作
    CommentResponse MGetCommentList(1: CommentRequest req)          // 获取评论列表
}