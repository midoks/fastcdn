use crate::rpc::client::rpc::{self, RequestAuth};
use crate::rpc::fastcdn::{
    AdminCreateRequest, AdminCreateResponse, AdminLoginRequest, AdminLoginResponse,
};

impl rpc::CommonRpc {
    pub async fn login(
        &mut self,
        req: AdminLoginRequest,
    ) -> Result<AdminLoginResponse, Box<dyn std::error::Error>> {
        let request = self.prepare_request_with_metadata(req, RequestAuth::ADMIN)?;
        self.make_grpc_call(request, "/fastcdn.Admin/login").await
    }
    pub async fn create(
        &mut self,
        req: AdminCreateRequest,
    ) -> Result<AdminCreateResponse, Box<dyn std::error::Error>> {
        let request = self.prepare_request_with_metadata(req, RequestAuth::ADMIN)?;
        self.make_grpc_call(request, "/fastcdn.Admin/create").await
    }
}
