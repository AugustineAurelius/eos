// Code generated by generator, DO NOT EDIT.
package example_builder

import uuid "github.com/google/uuid"

type userbuilder struct {
	inner User
}

func NewUserBuilder() *userbuilder {
	return &userbuilder{}
}
func (b *userbuilder) Name(name string) *userbuilder {
	b.inner.name = name
	return b
}
func (b *userbuilder) Surname(surname string) *userbuilder {
	b.inner.surname = surname
	return b
}
func (b *userbuilder) SetId(id uuid.UUID) *userbuilder {
	b.inner.id = id
	return b
}
func (b *userbuilder) SetAddesses(addesses []string) *userbuilder {
	b.inner.addesses = addesses
	return b
}
func (b *userbuilder) AddOneToAddesses(one string) *userbuilder {
	b.inner.addesses = append(b.inner.addesses, one)
	return b
}
func (b *userbuilder) SetAddressesID(addressesID []uuid.UUID) *userbuilder {
	b.inner.addressesID = addressesID
	return b
}
func (b *userbuilder) AddOneToAddressesID(one uuid.UUID) *userbuilder {
	b.inner.addressesID = append(b.inner.addressesID, one)
	return b
}
func (b *userbuilder) SetInner(inner InnerUser) *userbuilder {
	b.inner.inner = inner
	return b
}
func (b *userbuilder) Pointer(pointer *string) *userbuilder {
	b.inner.pointer = pointer
	return b
}
func (b *userbuilder) SetPointerID(pointerID *uuid.UUID) *userbuilder {
	b.inner.pointerID = pointerID
	return b
}
func (b *userbuilder) SetPointerSliceID(pointerSliceID *[]uuid.UUID) *userbuilder {
	b.inner.pointerSliceID = pointerSliceID
	return b
}
func (b *userbuilder) AddOneToPointerSliceID(one uuid.UUID) *userbuilder {
	*b.inner.pointerSliceID = append(*b.inner.pointerSliceID, one)
	return b
}
func (b *userbuilder) SetPointerSlicePointer(pointerSlicePointer *[]string) *userbuilder {
	b.inner.pointerSlicePointer = pointerSlicePointer
	return b
}
func (b *userbuilder) AddOneToPointerSlicePointer(one string) *userbuilder {
	*b.inner.pointerSlicePointer = append(*b.inner.pointerSlicePointer, one)
	return b
}
func (b *userbuilder) Build() User {
	return b.inner
}
