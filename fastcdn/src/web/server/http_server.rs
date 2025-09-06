use super::static_handler::StaticHandler;
use crate::web::app;
use actix_cors::Cors;
use actix_web::{App, HttpServer, web};

/// HTTP服务器配置和启动
pub struct HttpServerManager;

impl HttpServerManager {
    /// 创建并启动HTTP服务器
    pub async fn start() -> std::io::Result<()> {
        let config_intance = fastcdn_common::config::ConfigServer::Server::instance()
            .map_err(|e| std::io::Error::new(std::io::ErrorKind::Other, e.to_string()))?;

        match config_intance.lock() {
            Ok(config) => {
                let http_listen_raw = config.get_http_addresses()[0];
                let http_listen = http_listen_raw.replace("\"", "");
                println!("web start: {:?}", http_listen);

                let service = HttpServer::new(|| {
                    let cors = Cors::default()
                        .allow_any_origin()
                        .allow_any_method()
                        .allow_any_header()
                        .max_age(3600);

                    App::new()
                        .wrap(cors)
                        .service(
                            web::resource("/static/{_:.*}")
                                .route(web::get().to(StaticHandler::handle_static)),
                        )
                        .service(web::scope("/api").service(app::api::hello))
                        .service(
                            web::scope("/setup")
                                .service(app::setup::db_test_post)
                                .service(app::setup::db_test_get),
                        )
                        .route("/", web::get().to(StaticHandler::index))
                })
                .bind(&http_listen)
                .map_err(|e| std::io::Error::new(std::io::ErrorKind::Other, e.to_string()))?;

                match service.run().await {
                    Ok(_) => Ok(()),
                    Err(e) => {
                        eprintln!("server startup failed: {}", e);
                        std::process::exit(1);
                    }
                }
            }
            Err(_e) => {
                return Err(std::io::Error::new(
                    std::io::ErrorKind::Other,
                    "Failed to acquire config lock",
                ));
            }
        }
    }
}
