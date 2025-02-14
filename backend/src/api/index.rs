use crate::logging::log_request;
use actix_web::{get, HttpRequest, HttpResponse, Responder};
use serde_json::Value;

#[get("/")]
pub async fn get_index(req: HttpRequest) -> impl Responder {
    log_request(&req);

    let response = serde_json::json!({
        "body": "",
        "etag": null,
        "success": true
    });

    HttpResponse::Ok().json(response)
}
