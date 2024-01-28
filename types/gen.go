package types

type Person struct {
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Age         int      `json:"age"`
	Generation  int      `json:"generation"`
	ParentID    int      `json:"parentID"`
	ChildrenIDs []int    `json:"childrenIDs"`
	Death       bool     `json:"death"`
	Marriage    Marriage `json:"marriage"`
	BornedYear  int      `json:"bornedYear"`
	Sex         bool     `json:"sex"`
	Nationality string   `json:"nationality"`
	ID          int      `json:"id"`
}

type Marriage struct {
	ManID     int  `json:"manID"`
	WomanID   int  `json:"womanID"`
	Status    bool `json:"status"`
	MarryYear int  `json:"marryYear"`
	ID        int  `json:"id"`
}
