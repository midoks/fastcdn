use actix_web::{Responder, get, post, web};
use serde::{Deserialize, Serialize};
use std::process::Command;

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
    };

    let db_cfg = fastcdn_common::config::db::Db {
        user: req.username.clone(),
        password: req.password.clone(),
        database: req.dbname.clone(),
        host: req.hostname.clone() + ":" + &req.port.to_string(),
    };
    let _ = db_cfg.write();
    let _ = db_cfg.write_api();

    let mut result_map: serde_json::Value = serde_json::json!({});

    if req.api_type == "new" {
        // 安装API节点
        let output = Command::new("bin/fastcdn-api")
            .current_dir("fastcdn-api")
            .arg("setup")
            .arg("--protocol=http")
            .arg("--host=127.0.0.1")
            .arg("--port=10001")
            .output()
            .expect("Failed to execute command");

        if !output.status.success() {
            return web::Json(InstallResponse {
                message: format!("install error: {}", String::from_utf8_lossy(&output.stderr)),
                status: -1,
            });
        }

        result_map = serde_json::from_slice(&output.stdout).unwrap_or_default();

        println!("output: {}", String::from_utf8_lossy(&output.stdout));
        println!("result_map:{}", result_map);
        println!("result_map:{:?}", result_map.get("node_id"));

        // 关闭正在运行的API节点，防止冲突
        let _ = Command::new("bin/fastcdn-api")
            .current_dir("fastcdn-api")
            .arg("stop")
            .output()
            .expect("Failed to execute command");

        // 启动API节点
        let _ = Command::new("bin/fastcdn-api")
            .current_dir("fastcdn-api")
            .arg("start")
            .arg("-d")
            .output()
            .expect("Failed to execute command");

        println!("{}", req.api_type);
    } else if req.api_type == "old" {
        println!("{}", req.api_type);
    }

    web::Json(InstallResponse {
        message: "ok".to_string(),
        status: 0,
    })
}

#[get("/install")]
pub async fn install_get() -> impl Responder {
    web::Json(InstallResponse {
        message: "ok".to_string(),
        status: 0,
    })
}
