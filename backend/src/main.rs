mod api;
mod logging;
mod types;

use actix_web::{web, App, HttpServer};
use api::index::get_index;
use api::recipes::get_recipes;
use dotenv::dotenv;
use sqlx::mysql::MySqlPool;
use sqlx::Row;
use std::collections::HashMap;
use std::env;
use std::sync::{Arc, RwLock};
use types::CategoryItem;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();

    let db_host = env::var("DB_Server").unwrap_or_else(|_| "localhost".to_string());
    let db_user = env::var("DB_User").expect("DB_User is required");
    let db_password = env::var("DB_Password").expect("DB_Password is required");
    let db_name = env::var("DB_Name").expect("DB_Name is required");
    let db_port = env::var("DB_Port").unwrap_or_else(|_| "3306".to_string());

    let database_url = format!(
        "mysql://{}:{}@{}:{}/{}",
        db_user, db_password, db_host, db_port, db_name
    );

    println!("Connecting to the database");

    let pool = MySqlPool::connect(&database_url).await.map_err(|e| {
        eprintln!("Failed to connect to database: {:?}", e);
        std::io::Error::new(std::io::ErrorKind::Other, "Database connection failed")
    })?;

    println!("Connected to database successfully!");

    let category_cache = Arc::new(RwLock::new(
        load_categories(&pool).await.unwrap_or_default(),
    ));

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .app_data(web::Data::new(category_cache.clone()))
            .service(get_index)
            .service(get_recipes)
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}

async fn load_categories(pool: &MySqlPool) -> Result<Vec<CategoryItem>, sqlx::Error> {
    let query = "SELECT * FROM categoryitemsview";

    let categories = sqlx::query_as::<_, CategoryItem>(query)
        .fetch_all(pool)
        .await
        .map_err(|e| {
            eprintln!("Failed to load categories: {:?}", e);
            std::io::Error::new(std::io::ErrorKind::Other, "Database connection failed")
        })?;

    println!("Loaded {} categories into cache", categories.len());

    Ok(categories)
}
