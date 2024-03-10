package moddel

import "go.mongodb.org/mongo-driver/bson/primitive"

type USER struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FIRST_NAME string             `json:"first_name,omitempty" validate:"required,min=2"`
	LAST_NAME  string             `json:"last_name,omitempty" validate:"required,min=2"`
	Email      string             `json:"email,omitempty"  validate:"required"`
	Phone      string             `json:"phone,omitempty" validate:"required,max=10"`
	Password   string             `json:"password" validate:"requried,min=6"`
}
