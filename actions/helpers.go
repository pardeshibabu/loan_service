package actions

// SearchParams represents common search parameters
// type SearchParams struct {
// 	Query      string   `form:"q"`
// 	SortBy     string   `form:"sort_by"`
// 	SortOrder  string   `form:"sort_order"`
// 	Filters    []string `form:"filters[]"`
// 	PageSize   int      `form:"page_size"`
// 	PageNumber int      `form:"page"`
// }

// // ApplySearch applies search parameters to a query
// func ApplySearch(q *pop.Query, search SearchParams) *pop.Query {
// 	// Apply text search if provided
// 	if search.Query != "" {
// 		q = q.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?",
// 			"%"+search.Query+"%", "%"+search.Query+"%", "%"+search.Query+"%")
// 	}

// 	// Apply sorting
// 	if search.SortBy != "" {
// 		order := "asc"
// 		if search.SortOrder == "desc" {
// 			order = "desc"
// 		}
// 		q = q.Order(search.SortBy + " " + order)
// 	}

// 	return q
// }

func strPtr(s string) *string {
	return &s
}
