use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;

#[derive(Serialize)]
pub struct GenericApiReply {
    pub body: serde_json::Value, // Use serde_json::Value for dynamic "any" type
    #[serde(skip_serializing_if = "Option::is_none")] // Skip "etag" if None
    pub etag: Option<String>,
    pub success: bool,
}

/// Structure to represent a category item from `categoryitemsview`
#[derive(Serialize, Deserialize, FromRow)]
pub struct CategoryItem {
    pub itemid: u16,
    pub itemname: String,
    pub itemicon: String,
    pub itemmodified: DateTime<Utc>,
    pub catid: u8,
    pub catname: String,
    pub caticon: String,
    pub catmodified: DateTime<Utc>,
}

#[derive(Serialize)]
pub struct RecipeListResponse {
    pub count: usize,
    pub limit: usize,
    pub list: Vec<Recipe>,
    pub cache: bool,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize, FromRow)]
pub struct Recipe {
    pub recipe_id: u32,
    pub user_id: Option<u32>,
    pub edit_user_id: Option<u32>,
    pub aigenerated: u8,
    pub localized: u8,
    pub recipe_placeholder: u8,
    pub recipe_public_internal: u8,
    pub recipe_public_external: u8,
    pub recipe_name: String,
    pub recipe_description: String,
    pub recipe_eater: u8,
    pub recipe_source_desc: String,
    pub recipe_source_url: String,
    pub recipe_created: chrono::DateTime<Utc>,
    pub recipe_edited: Option<chrono::DateTime<Utc>>,
    pub recipe_modified: chrono::DateTime<Utc>,
    pub recipe_published: chrono::DateTime<Utc>,
    pub ingredientsGroupByStep: u8,
    pub stepscount: i32,
    pub preparationtime: i32,
    pub cookingtime: i32,
    pub chilltime: i32,
    pub votesum: Option<i32>,
    pub votes: Option<i32>,
    pub avgvotes: Option<f32>,
    pub ratesum: Option<i32>,
    pub ratings: Option<i32>,
    pub avgratings: Option<f32>,
    pub categories: Vec<String>,
    pub pictures: Vec<PictureInfo>,
}

#[derive(Serialize, Deserialize)]
pub struct EditInfo {
    pub user: UserInfo,
    pub when: String,
}

#[derive(Serialize, Deserialize)]
pub struct UserInfo {
    pub id: i32,
    pub name: String,
    pub meta: UserMeta,
    pub statistics: UserStatistics,
}

#[derive(Serialize, Deserialize)]
pub struct UserMeta {
    pub email: Option<String>,
    #[serde(rename = "fn")]
    pub fn_: String,
    pub ln: String,
    pub un: String,
    pub initials: String,
}

#[derive(Serialize, Deserialize)]
pub struct UserStatistics {
    pub recipes: RecipeStats,
}

#[derive(Serialize, Deserialize)]
pub struct RecipeStats {
    pub created: i32,
    pub createdext: i32,
}

#[derive(Serialize, Deserialize)]
pub struct SourceInfo {
    pub description: String,
    pub url: String,
}

#[derive(Serialize, Deserialize, FromRow)]
pub struct PictureInfo {
    pub picture_id: i32,
    pub picture_name: String,
    pub picture_description: String,
    pub picture_filename: String,
    pub picture_full_path: String,
    pub picture_uploaded: String,
    pub picture_width: i32,
    pub picture_height: i32,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize)]
pub struct PreparationInfo {
    pub ingredients: Vec<String>,
    pub ingredientsGrouping: String,
    pub steps: Vec<String>,
    pub stepscount: i32,
    pub timeConsumed: TimeInfo,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize)]
pub struct TimeInfo {
    pub cooking: i32,
    pub preparing: i32,
    pub rest: i32,
    pub total: i32,
    pub unit: String,
}

#[derive(Serialize, Deserialize)]
pub struct SocialInfo {
    pub cooked: i32,
    pub views: i32,
    pub sharing: SharingInfo,
    pub rating: RatingInfo,
    pub voting: VotingInfo,
    pub myvotes: Option<String>,
}

#[derive(Serialize, Deserialize)]
pub struct SharingInfo {
    pub links: Vec<String>,
    pub publication: PublicationInfo,
}

#[allow(non_snake_case)]
#[derive(Serialize, Deserialize)]
pub struct PublicationInfo {
    pub isPublished: PublishStatus,
    pub when: String,
}

#[derive(Serialize, Deserialize)]
pub struct PublishStatus {
    pub internal: bool,
    pub external: bool,
}

#[derive(Serialize, Deserialize)]
pub struct RatingInfo {
    pub ratings: i32,
    pub sum: i32,
    pub avg: Option<f32>,
}

#[derive(Serialize, Deserialize)]
pub struct VotingInfo {
    pub votes: i32,
    pub sum: i32,
    pub avg: Option<f32>,
}
