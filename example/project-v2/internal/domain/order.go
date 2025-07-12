package domain

import (
    "github.com/google/uuid"
)

type Order struct {
    ID uuid.UUID
    Product string
    Quantity int
}


