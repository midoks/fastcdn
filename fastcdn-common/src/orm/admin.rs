use crate::{db::pool, utils};

pub async fn count() -> Result<i64, Box<dyn std::error::Error>> {
    let db = pool::Manager::instance().await?;
    let results = db.count("admin", None).await?;
    Ok(results)
}

pub async fn add(
    username: &str,
    password: &str,
    fullname: &str,
    is_on: bool,
    is_super: bool,
    can_login: bool,
    theme: bool,
    lang: &str,
    state: bool,
) -> Result<u64, Box<dyn std::error::Error>> {
    let db = pool::Manager::instance().await?;
    let time_unix = utils::time::now_unix();
    let mut data = std::collections::HashMap::new();
    data.insert(
        "username".to_string(),
        serde_json::Value::String(username.to_string()),
    );
    data.insert(
        "password".to_string(),
        serde_json::Value::String(password.to_string()),
    );
    data.insert(
        "fullname".to_string(),
        serde_json::Value::String(fullname.to_string()),
    );
    data.insert(
        "is_on".to_string(),
        serde_json::Value::Bool(is_on),
    );
    data.insert(
        "is_super".to_string(),
        serde_json::Value::Bool(is_super),
    );
    data.insert(
        "state".to_string(),
        serde_json::Value::Bool(state),
    );
    data.insert(
        "theme".to_string(),
        serde_json::Value::String(theme.to_string()),
    );
    data.insert(
        "lang".to_string(),
        serde_json::Value::String(lang.to_string()),
    );
    data.insert(
        "created_at".to_string(),
        serde_json::Value::String(time_unix.clone()),
    );
    data.insert(
        "updated_at".to_string(),
        serde_json::Value::String(time_unix),
    );

    let id = db.insert("admin", &data).await?;
    Ok(id)
}
