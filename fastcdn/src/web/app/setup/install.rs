use actix_web::{Responder, get, post, web};
use serde::{Deserialize, Serialize};

// 必须为所有需要序列化/反序列化的结构体添加derive
#[derive(Debug, Serialize, Deserialize)] // 添加Debug方便日志记录
pub struct InstallResponse {
    pub message: String,
    pub status: i16,
}

#[derive(Debug, Serialize, Deserialize, Clone)] // 接收JSON请求的结构体
pub struct InstallRequest {
    pub api_host: String,
    pub api_port: u32,
    pub api_protocol: String,
    pub api_type: String,

    pub node_id: String,
    pub secret: String,

    pub hostname: String,
    pub port: u16,
    pub dbname: String,
    pub username: String,
    pub password: String,

    pub admin_username: String,
    pub admin_password: String,
}

#[post("/install")]
pub async fn install_post(req: web::Json<InstallRequest>) -> impl Responder {
    println!("{:?}", req);

    let config = fastcdn_common::db::pool::DbConfig {
        hostname: req.hostname.clone(),
        port: req.port,
        dbname: req.dbname.clone(),
        username: req.username.clone(),
        password: req.password.clone(),
    };

    match fastcdn_common::db::pool::Manager::new().await {
        Ok(db) => match db.test_connection(&config).await {
            Ok(_) => web::Json(InstallResponse {
                message: "ok".to_string(),
                status: 0,
            }),
            Err(e) => web::Json(InstallResponse {
                message: format!("error: {}", e),
                status: -1,
            }),
        },
        Err(e) => web::Json(InstallResponse {
            message: format!("db error: {}", e),
            status: -1,
        }),
    }
}

#[get("/install")]
pub async fn install_get() -> impl Responder {
    web::Json(InstallResponse {
        message: "ok".to_string(),
        status: 0,
    })
}
