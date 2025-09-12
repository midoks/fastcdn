use crate::{db::pool, utils};

pub async fn count() -> Result<i64, Box<dyn std::error::Error + Send + Sync>> {
    let db = pool::Manager::instance().await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;
    let results =
        db.count("admin", None)
            .await
            .map_err(|e| -> Box<dyn std::error::Error + Send + Sync> {
                Box::new(std::io::Error::new(
                    std::io::ErrorKind::Other,
                    e.to_string(),
                ))
            })?;
    Ok(results)
}

pub async fn update_admin_password(
    id: u64,
    password: &str,
) -> Result<bool, Box<dyn std::error::Error + Send + Sync>> {
    let db = pool::Manager::instance().await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;
    let table_name = db.get_table_name("admin");

    let update = db
        .update_builder("admin")
        .set_str("password", password)
        .where_id(id);

    let affected = db.update_with_builder(update).await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;

    Ok(affected > 0)
}

pub async fn find_admin_id_with_username(
    username: &str,
) -> Result<Vec<serde_json::Value>, Box<dyn std::error::Error + Send + Sync>> {
    let db = pool::Manager::instance().await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;

    let table_name = db.get_table_name("admin");
    let query = db
        .query_builder(&table_name)
        .select(&["id"])
        .limit(1)
        .where_eq("username", username);
    let results = db.query_with_builder(query).await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;
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
) -> Result<u64, Box<dyn std::error::Error + Send + Sync>> {
    let db = pool::Manager::instance().await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;
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
    data.insert("is_on".to_string(), serde_json::Value::Bool(is_on));
    data.insert("is_super".to_string(), serde_json::Value::Bool(is_super));
    data.insert("can_login".to_string(), serde_json::Value::Bool(can_login));
    data.insert("state".to_string(), serde_json::Value::Bool(state));
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

    let id = db.insert("admin", &data).await.map_err(
        |e| -> Box<dyn std::error::Error + Send + Sync> {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        },
    )?;
    Ok(id)
}
