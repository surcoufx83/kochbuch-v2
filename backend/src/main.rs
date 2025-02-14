use actix_web::{get, App, HttpResponse, HttpServer, Responder};
use serde::Serialize;

#[derive(Serialize)]
struct GenericApiReply {
    body: serde_json::Value, // Use serde_json::Value for dynamic "any" type
    #[serde(skip_serializing_if = "Option::is_none")] // Skip "etag" if None
    etag: Option<String>,
    success: bool,
}

#[get("/")]
async fn api_handler_index() -> impl Responder {
    let response = GenericApiReply {
        body: serde_json::Value::String("".to_string()), // Empty string for body
        etag: None,                                      // No etag provided
        success: true,
    };

    HttpResponse::Ok().json(response) // Return as JSON
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(api_handler_index))
        .bind("0.0.0.0:8080")?
        .run()
        .await
}
