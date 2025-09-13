use tonic::{Request, Response, Status};

use crate::db::pool;
use crate::orm;
use crate::rpc::auth::AuthMiddleware;
use crate::rpc::fastcdn::admin_server::Admin;
use crate::rpc::fastcdn::{
    AdminCreateRequest, AdminCreateResponse, AdminLoginRequest, AdminLoginResponse,
    CreateOrUpdateAdminRequest, CreateOrUpdateAdminResponse,
};

/// Admin 实现
#[derive(Debug, Default)]
pub struct FcAdmin {}

#[tonic::async_trait]
impl Admin for FcAdmin {
    async fn login(
        &self,
        request: Request<AdminLoginRequest>,
    ) -> Result<Response<AdminLoginResponse>, Status> {
        // 验证请求头认证
        AuthMiddleware::verify_admin_request(&request).await?;

        let login_req = request.into_inner();
        println!("admin login username: {:?}", login_req.username);
        println!("admin login password: {:?}", login_req.password);

        let reply = AdminLoginResponse {
            id: -1,       // 用户ID
            is_ok: false, // 登录是否成功
            message: "登陆失败".to_string(),
        };

        match pool::Manager::instance().await {
            Ok(db) => {
                println!("db: {:?}", db);
                println!("addr: {:p}", &db);
            }
            Err(e) => {
                println!("db manager fail: {:?}", e);
            }
        }

        Ok(Response::new(reply))
    }

    // 创建或修改管理员
    async fn create_or_update_admin(
        &self,
        request: Request<CreateOrUpdateAdminRequest>,
    ) -> Result<Response<CreateOrUpdateAdminResponse>, Status> {
        AuthMiddleware::verify_admin_request(&request).await?;

        let inner_request = request.get_ref();

        let adminids = orm::admin::find_admin_id_with_username(&inner_request.username)
            .await
            .map_err(|e| Status::internal(e.to_string()))?;

        let mut resp = CreateOrUpdateAdminResponse { id: 0 };
        if adminids.len() > 0 {
            let admin_id = adminids[0]["id"].as_u64().unwrap_or(0);
            orm::admin::update_admin_password(admin_id, &inner_request.password)
                .await
                .map_err(|e| Status::internal(e.to_string()))?;

            resp.id = admin_id as i64;
        } else {
            let admin_id = orm::admin::add(
                &inner_request.username,
                &inner_request.password,
                &inner_request.username,
                true,
                true,
                true,
                true,
                "zh",
                true,
            )
            .await
            .map_err(|e| Status::internal(e.to_string()))?;
            resp.id = admin_id as i64;
        }

        // println!("resp:{:?}", resp);
        Ok(Response::new(resp))
    }

    async fn create(
        &self,
        request: Request<AdminCreateRequest>,
    ) -> Result<Response<AdminCreateResponse>, Status> {
        // 验证请求头认证
        AuthMiddleware::verify_request(&request)?;

        let reply = AdminCreateResponse {
            id: 1, // 示例ID，实际应该从数据库生成
        };

        match pool::Manager::instance().await {
            Ok(manager) => println!("数据库管理器实例: {:?}", manager),
            Err(e) => println!("获取数据库管理器失败: {:?}", e),
        }

        Ok(Response::new(reply))
    }
}
