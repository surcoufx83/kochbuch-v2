use crate::logging::log_request;
use crate::types::{PictureInfo, Recipe, RecipeListResponse};
use actix_web::{get, web, HttpRequest, HttpResponse, Responder};
use sqlx::mysql::MySqlPool;
use sqlx::Row;

#[get("/recipes")]
pub async fn get_recipes(req: HttpRequest, db_pool: web::Data<MySqlPool>) -> impl Responder {
    log_request(&req);

    let query = "SELECT * FROM allrecipes_nouser";

    let rows = sqlx::query(query).fetch_all(db_pool.get_ref()).await;

    match rows {
        Ok(records) => {
            let mut recipes_map = std::collections::HashMap::new();

            for row in records.iter() {
                let recipe_id: u32 = row.try_get("recipe_id").unwrap();
                let categories: String =
                    row.try_get("categories").unwrap_or_else(|_| "".to_string());
                let categories_vec: Vec<String> = if categories.is_empty() {
                    vec![]
                } else {
                    categories
                        .split(',')
                        .map(|s| s.trim().to_string())
                        .collect()
                };

                let picture = PictureInfo {
                    picture_id: row.try_get("picture_id").unwrap(),
                    picture_name: row.try_get("picture_name").unwrap(),
                    picture_description: row.try_get("picture_description").unwrap(),
                    picture_filename: row.try_get("picture_filename").unwrap(),
                    picture_full_path: row.try_get("picture_full_path").unwrap(),
                    picture_uploaded: row.try_get("picture_uploaded").unwrap(),
                    picture_width: row.try_get("picture_width").unwrap(),
                    picture_height: row.try_get("picture_height").unwrap(),
                };

                let recipe = recipes_map.entry(recipe_id).or_insert_with(|| Recipe {
                    recipe_id,
                    user_id: row.try_get("user_id").ok(),
                    edit_user_id: row.try_get("edit_user_id").ok(),
                    aigenerated: row.try_get("aigenerated").unwrap(),
                    localized: row.try_get("localized").unwrap(),
                    recipe_placeholder: row.try_get("recipe_placeholder").unwrap(),
                    recipe_public_internal: row.try_get("recipe_public_internal").unwrap(),
                    recipe_public_external: row.try_get("recipe_public_external").unwrap(),
                    recipe_name: row.try_get("recipe_name").unwrap(),
                    recipe_description: row.try_get("recipe_description").unwrap(),
                    recipe_eater: row.try_get("recipe_eater").unwrap(),
                    recipe_source_desc: row.try_get("recipe_source_desc").unwrap(),
                    recipe_source_url: row.try_get("recipe_source_url").unwrap(),
                    recipe_created: row.try_get("recipe_created").unwrap(),
                    recipe_modified: row.try_get("recipe_modified").unwrap(),
                    stepscount: row.try_get("stepscount").unwrap(),
                    preparationtime: row.try_get("preparationtime").unwrap(),
                    cookingtime: row.try_get("cookingtime").unwrap(),
                    chilltime: row.try_get("chilltime").unwrap(),
                    votesum: row.try_get("votesum").ok(),
                    votes: row.try_get("votes").ok(),
                    avgvotes: row.try_get("avgvotes").ok(),
                    ratesum: row.try_get("ratesum").ok(),
                    ratings: row.try_get("ratings").ok(),
                    avgratings: row.try_get("avgratings").ok(),
                    categories: categories_vec,
                    pictures: vec![],
                });

                recipe.pictures.push(picture);
            }

            let response = RecipeListResponse {
                count: recipes_map.len(),
                limit: 25,
                list: recipes_map.into_values().collect(),
                cache: true,
            };

            HttpResponse::Ok().json(response)
        }
        Err(e) => {
            eprintln!("Database query failed: {:?}", e);
            HttpResponse::InternalServerError().json("Failed to fetch recipes")
        }
    }
}
