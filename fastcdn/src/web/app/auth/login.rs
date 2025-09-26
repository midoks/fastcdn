use actix_web::{Responder, post, web};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

// 必须为所有需要序列化/反序列化的结构体添加derive
#[derive(Debug, Serialize, Deserialize)] // 添加Debug方便日志记录
pub struct LoginResponseData {
    pub token: String,
}

// 必须为所有需要序列化/反序列化的结构体添加derive
#[derive(Debug, Serialize, Deserialize)] // 添加Debug方便日志记录
pub struct LoginResponse {
    pub message: String,
    pub status: i16,
    pub token: String,
    // pub data: Option<LoginResponseData>,
}

#[derive(Debug, Serialize, Deserialize, Clone)] // 接收JSON请求的结构体
pub struct LoginRequest {
    pub username: String,
    pub password: String,
}

#[post("/login")]
pub async fn login_post(
    req: web::Json<LoginRequest>,
) -> Result<impl Responder, Box<dyn std::error::Error>> {
    // println!("{:?}", req);
    let mut admin_rpc = fastcdn_common::rpc::client::CommonRpc::admin_rpc().await?;

    let rpc_req = fastcdn_common::rpc::fastcdn::AdminLoginRequest {
        username: req.username.clone(),
        password: req.password.clone(),
    };

    let resp = Arc::get_mut(&mut admin_rpc)
        .ok_or("failed to get mutable reference to admin_rpc")?
        .login(rpc_req)
        .await?;

    println!("resp:{:?}", resp);
    if resp.id == 0 {
        return Ok(web::Json(LoginResponse {
            message: "登陆失败".to_string(),
            status: -1,
            token: Some(token?),
        }));
    }

    let token = fastcdn::utils::jwt::create(&resp.id.to_string());
    // println!("token:{:?}", token);

    let resp_data = LoginResponse {
        message: "ok".to_string(),
        status: 0,
        token: Some(token?),
    };

    println!("resp_data:{:?}", resp_data);
    Ok(web::Json(resp_data))
}
