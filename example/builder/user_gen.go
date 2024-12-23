// Code generated by generator, DO NOT EDIT.
package example_builder

import uuid "github.com/google/uuid"

type userbuilder struct {
	inner *User
}

func NewUserBuilder() *userbuilder {
	return &userbuilder{}
}
func (b *userbuilder) Name(name string) {
	b.inner.name = name
}
func (b *userbuilder) Surname(surname string) {
	b.inner.surname = surname
}
func (b *userbuilder) SetId(id uuid.UUID) {
	b.inner.id = id
}
func (b *userbuilder) SetAddesses(addesses []string) {
	b.inner.addesses = addesses
}
func (b *userbuilder) AddOneToAddesses(one string) {
	b.inner.addesses = append(b.inner.addesses, one)
}
func (b *userbuilder) SetAddressesID(addressesID []uuid.UUID) {
	b.inner.addressesID = addressesID
}
func (b *userbuilder) AddOneToAddressesID(one uuid.UUID) {
	b.inner.addressesID = append(b.inner.addressesID, one)
}
func (b *userbuilder) SetInner(inner InnerUser) {
	b.inner.inner = inner
}
func (b *userbuilder) Build() User {
	return *b.inner
}
