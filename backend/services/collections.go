package services

import (
	"database/sql"
	"fmt"
	"kochbuch-v2-backend/types"
	"log"
	"time"
)

func createCollection(user *types.UserProfile, title string, description string) (*types.Collection, error) {
	fn := "createCollection()"
	// log.Printf("%v", fn)

	query := "INSERT INTO `user_collections`(`user_id`, `title`, `description`) VALUES(?, ?, ?)"
	stmt, err := dbPrepareStmt("createCollection", query)
	if err != nil {
		log.Printf("%v: Failed preparing stmt: %v", fn, stmt)
		return nil, err
	}

	result, err := stmt.Exec(user.Id, title, description)
	if err != nil {
		log.Printf("%v: Failed executing stmt: %v", fn, stmt)
		return nil, err
	}

	newid, err := result.LastInsertId()
	if err != nil {
		log.Printf("%v: Failed retrieving insertId: %v", fn, stmt)
		return nil, err
	}

	coll := &types.Collection{
		Created:     time.Now(),
		Deleted:     types.NullTime{Valid: false},
		Description: description,
		Id:          uint32(newid),
		Modified:    time.Now(),
		Name:        title,
		Published:   types.NullTime{Valid: false},
		UserId:      user.Id,
		Items:       []*types.CollectionItem{},
	}

	userMutex.Lock()
	defer userMutex.Unlock()
	if user.Collections == nil {
		user.Collections = make(map[uint32]*types.Collection)
	}
	user.Collections[coll.Id] = coll

	if coll.Modified.After(user.Modified) {
		defer setUserModified(user, coll.Modified)
	}

	log.Printf("%v: Created new collection %s for user %s", fn, coll.Name, user.DisplayName)

	return coll, nil
}

func loadCollections() {
	fn := "loadCollections()"
	//log.Printf("%v: Loading user collections", fn)

	query := "SELECT * FROM `user_collections`"
	var colls []types.Collection

	err := Db.Select(&colls, query)
	if err != nil {
		log.Fatalf("%v: Failed: %v", fn, err)
	}

	for _, coll := range colls {
		user, found := userCache[coll.UserId]
		if !found {
			log.Fatalf("%v: Loaded collection %d with unknown user %d", fn, coll.Id, coll.UserId)
		}

		if user.Collections == nil {
			user.Collections = make(map[uint32]*types.Collection)
		}

		user.Collections[coll.Id] = &coll

		if coll.Modified.After(user.Modified) {
			setUserModified(user, coll.Modified)
		}
	}
	log.Printf("Loaded %d user collections into cache", len(colls))
}

func loadCollectionItems() {
	fn := "loadCollectionItems()"
	//log.Printf("%v: Loading user collection items", fn)

	query := "SELECT `cr`.*, `c`.`user_id` FROM `user_collection_recipes` `cr` JOIN `user_collections` `c` ON `c`.`collection_id` = `cr`.`collection_id`"
	var items []types.CollectionItem

	err := Db.Select(&items, query)
	if err != nil {
		log.Fatalf("%v: Failed: %v", fn, err)
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	for _, item := range items {
		user, found := userCache[item.UserId]
		if !found {
			log.Fatalf("%v: Loaded collection item (Collection %d, Recipe %d) with unknown user %d", fn, item.CollectionId, item.RecipeId, item.UserId)
		}

		if user.Collections == nil {
			user.Collections = make(map[uint32]*types.Collection)
		}

		coll, found := user.Collections[item.CollectionId]
		if !found {
			log.Fatalf("%v: Loaded collection item (Collection %d, Recipe %d, User %d) but collection is not loaded", fn, item.CollectionId, item.RecipeId, item.UserId)
		}

		if coll.Items == nil {
			coll.Items = []*types.CollectionItem{}
		}

		coll.Items = append(coll.Items, &item)

		if item.Modified.After(coll.Modified) {
			setUserCollectionModified(nil, coll, item.Modified)

			if item.Modified.After(user.Modified) {
				setUserModified(user, item.Modified)
			}
		}
	}
	log.Printf("Loaded %d collection items into cache", len(items))
}

func setUserCollectionModified(tx *sql.Tx, coll *types.Collection, time time.Time) {
	fn := fmt.Sprintf("setUserCollectionModified(%d, %s)", coll.Id, time)
	log.Print(fn)
	query := "UPDATE `user_collections` SET `modified` = ? WHERE `collection_id` = ?"

	stmt, err := dbPrepareStmt("setUserCollectionModified", query)
	if err != nil {
		log.Printf("%s: Failed to prepare stmt: %v", fn, err)
		return
	}

	if tx != nil {
		stmt = tx.Stmt(stmt)
	}

	_, err = stmt.Exec(time, coll.Id)
	if err != nil {
		log.Printf("%s: Failed to exec stmt: %v", fn, err)
		return
	}

	coll.Modified = time
}
