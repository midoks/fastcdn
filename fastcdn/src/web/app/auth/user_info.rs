use crate::web::middleware::UserId;
use actix_web::{HttpMessage, HttpRequest, Responder, web};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct UserInfoResponse {
    pub message: String,
    pub code: i16,
    pub user_id: i64,
}

// 使用路由配置而不是宏
pub async fn user_info(req: HttpRequest) -> impl Responder {
    // 从请求扩展中获取用户ID
    let user_id = req.extensions().get::<UserId>().map(|id| id.0).unwrap_or(0);

    let response = UserInfoResponse {
        message: "获取用户信息成功".to_string(),
        code: 0,
        user_id,
    };

    web::Json(response)
}