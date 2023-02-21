package entities

type Search struct {
	Column string `json:"column"`
	Action string `json:"action"`
	Query  string `json:"query"`
}

type Pagination struct {
	Limit         int         `json:"limit"`
	Page          int         `json:"page"`
	SortByCrated  string      `json:"sort_created"`
	SortByComment string      `json:"sort_comment"`
	SortByLiked   string      `json:"sort_like"`
	TotalRows     int         `json:"total_rows"`
	FirstPage     string      `json:"first_page"`
	PreviousPage  string      `json:"previous_page"`
	NextPage      string      `json:"next_page"`
	LastPage      string      `json:"last_page"`
	FromRow       int         `json:"from_row"`
	ToRow         int         `json:"to_row"`
	Rows          interface{} `json:"rows"`
	Searchs       []Search    `json:"searchs"`
}
