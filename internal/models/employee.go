package models

type Employee struct {
	Name      string `json:"displayName"`
	Email     string `json:"email"`
	WorkPhone string `json:"workPhone"`
	Id        int    `json:"id"`
}
