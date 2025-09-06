use actix_web::{Responder, get, post, web};
use serde::{Deserialize, Serialize};

// 必须为所有需要序列化/反序列化的结构体添加derive
#[derive(Debug, Serialize, Deserialize)] // 添加Debug方便日志记录
pub struct DbTestResponse {
    pub message: String,
    pub status: u16,
}

#[derive(Debug, Deserialize)] // 接收JSON请求的结构体
pub struct DbTestRequest {
    pub host: String,
}

#[post("/db_test")]
pub async fn db_test_post(req: web::Json<DbTestRequest>) -> impl Responder {
    println!("{:?}", req);
    web::Json(DbTestResponse {
        message: "ok".to_string(),
        status: 0,
    })
}

#[get("/db_test")]
pub async fn db_test_get() -> impl Responder {
    web::Json(DbTestResponse {
        message: "ok".to_string(),
        status: 0,
    })
}
