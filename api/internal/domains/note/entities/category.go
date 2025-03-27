package entities

type Category struct {
	ID                 string     `json:"key" db:"category_id"`
	OwnerID            int64      `json:"ownerID" db:"owner_id"`
	Name               string     `json:"label" db:"name"`
	ParentCategoryID   *string    `json:"parentCategoryID" db:"parent_category_id"`
	ChildrenCategories []Category `json:"children" db:"-"`
}
