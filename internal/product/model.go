package product

type Product struct {
	ID              string          `json:"id" bson:"_id,omitempty"`
	Name            string          `json:"name" bson:"name"`
	Description     string          `json:"description,omitempty" bson:"description,omitempty"`
	Type            string          `json:"type" bson:"type"`
	InStock         int             `json:"in_stock" bson:"in_stock"`
	Characteristics Characteristics `json:"characteristics,omitempty" bson:"characteristics,omitempty"`
}

type Characteristics struct {
	Color    string `json:"color,omitempty" bson:"color,omitempty"`
	Material string `json:"material,omitempty" bson:"material,omitempty"`
}

type CreateProductDTO struct {
	Name            string          `json:"name" bson:"name" binding:"required"`
	Description     string          `json:"description,omitempty" bson:"description,omitempty"`
	Type            string          `json:"type" bson:"type" binding:"required"`
	InStock         int             `json:"in_stock" bson:"in_stock" binding:"required"`
	Characteristics Characteristics `json:"characteristics,omitempty" bson:"characteristics,omitempty"`
}

func NewProduct(dto CreateProductDTO) Product {
	return Product{
		Name:            dto.Name,
		Description:     dto.Description,
		Type:            dto.Type,
		InStock:         dto.InStock,
		Characteristics: dto.Characteristics,
	}
}

type UpdateProductDTO struct {
	Name            string          `json:"name,omitempty" bson:"name,omitempty"`
	Description     string          `json:"description,omitempty" bson:"description,omitempty"`
	Type            string          `json:"type,omitempty" bson:"type,omitempty"`
	InStock         int             `json:"in_stock,omitempty" bson:"in_stock,omitempty"`
	Characteristics Characteristics `json:"characteristics,omitempty" bson:"characteristics,omitempty"`
}

func UpdateProduct(dto UpdateProductDTO) Product {
	return Product{
		Name:            dto.Name,
		Description:     dto.Description,
		Type:            dto.Type,
		InStock:         dto.InStock,
		Characteristics: dto.Characteristics,
	}
}
