// Code generated by generator, DO NOT EDIT.
package example_builder

import uuid "github.com/google/uuid"

type userbuilder struct {
	inner *User
}

func NewUserBuilder() *userbuilder {
	return &userbuilder{}
}
func (b *userbuilder) SetName(Name string) {
	b.inner.Name = Name
}
func (b *userbuilder) SetSurname(Surname string) {
	b.inner.Surname = Surname
}
func (b *userbuilder) SetID(ID uuid.UUID) {
	b.inner.ID = ID
}
func (b *userbuilder) SetAddesses(Addesses []string) {
	b.inner.Addesses = Addesses
}
func (b *userbuilder) AddOneToAddesses(one string) {
	b.inner.Addesses = append(b.inner.Addesses, one)
}
func (b *userbuilder) SetAddressesID(AddressesID []uuid.UUID) {
	b.inner.AddressesID = AddressesID
}
func (b *userbuilder) AddOneToAddressesID(one uuid.UUID) {
	b.inner.AddressesID = append(b.inner.AddressesID, one)
}
func (b *userbuilder) Build() User {
	return *b.inner
}