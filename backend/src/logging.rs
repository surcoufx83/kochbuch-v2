use actix_web::HttpRequest;
use chrono::Utc;

pub fn log_request(req: &HttpRequest) {
    let method = req.method();
    let path = req.path();
    let timestamp = Utc::now().format("%Y-%m-%d %H:%M:%S%.3f");
    let ip = req
        .peer_addr()
        .map(|addr| addr.ip().to_string())
        .unwrap_or_else(|| "Unknown".to_string());

    println!("[{}] {} {} from {}", timestamp, method, path, ip);
}
