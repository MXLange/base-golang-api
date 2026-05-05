package dto

type CreateIN struct {
	Name *string `json:"name" example:"My User"`
	NestedField *NestedField `json:"nestedField"`
}

type NestedField struct {
	Data *string `json:"data" example:"Some data"`
}

var CreateInSchema = `{
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"example": "My User"
		},
		"nestedField": {
			"type": "object",
			"properties": {
				"data": {
					"type": "string",
					"example": "Some data"
				}
			},
			"required": ["data"]
		}
	},
	"required": ["name", "nestedField"]
}`

type CreateOUT struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"My User"`
}
