package ims_server

type Product struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"quantity"`
	Price       float64 `json:"price"`
}
