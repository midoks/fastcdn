use serde::{Deserialize, Serialize};

pub mod api;
pub mod auth;
pub mod setup;

// 必须为所有需要序列化/反序列化的结构体添加derive
#[derive(Debug, Serialize, Deserialize)] // 添加Debug方便日志记录
pub struct DataResponse {
    pub message: String,
    pub status: i16,
}
