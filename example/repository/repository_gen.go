//Code generated by generator, DO NOT EDIT.
package repository

import (
	sq "github.com/Masterminds/squirrel"
	common "github.com/AugustineAurelius/eos/example/common"
    "github.com/google/uuid"
)



type repository struct{
	db common.Querier
}

func New(db common.Querier) *repository {
	return &repository{
		db: db,
	}
}

// UserFiler represents the User filter.
type UserFilter struct {
  id *uuid.UUID
  ids []uuid.UUID
  name *string
  names []string
  email **string
  emails []*string
}

func NewFilter() *UserFilter{
	return &UserFilter{}
}
func (f *UserFilter) ID (id uuid.UUID)  *UserFilter {
  f.id = &id
  return f
}

func (f *UserFilter) AddOneToIDs (id uuid.UUID)  *UserFilter {
  f.ids = append(f.ids, id)
  return f
}

func (f *UserFilter) IDs (ids []uuid.UUID)  *UserFilter {
  f.ids =  ids
  return f
}
func (f *UserFilter) Name (name string)  *UserFilter {
  f.name = &name
  return f
}

func (f *UserFilter) AddOneToNames (name string)  *UserFilter {
  f.names = append(f.names, name)
  return f
}

func (f *UserFilter) Names (names []string)  *UserFilter {
  f.names =  names
  return f
}
func (f *UserFilter) Email (email *string)  *UserFilter {
  f.email = &email
  return f
}

func (f *UserFilter) AddOneToEmails (email *string)  *UserFilter {
  f.emails = append(f.emails, email)
  return f
}

func (f *UserFilter) Emails (emails []*string)  *UserFilter {
  f.emails =  emails
  return f
}

func (f *UserFilter) Build()  UserFilter {
   return *f
}

func ApplyWhere[B interface {
    Where(pred interface{}, args ...interface{}) B
}](b B,f UserFilter) B {
	if f.id != nil {
      b = b.Where(sq.Eq{ColumnUserID: *f.id})
    }
	if f.ids != nil {
      b = b.Where(sq.Eq{ColumnUserID: f.ids})
    }
	if f.name != nil {
      b = b.Where(sq.Eq{ColumnUserName: *f.name})
    }
	if f.names != nil {
      b = b.Where(sq.Eq{ColumnUserName: f.names})
    }
	if f.email != nil {
      b = b.Where(sq.Eq{ColumnUserEmail: *f.email})
    }
	if f.emails != nil {
      b = b.Where(sq.Eq{ColumnUserEmail: f.emails})
    }
  return b
}

// UserUpdate represents the User update struct.
type UserUpdate struct {
  id *uuid.UUID
  name *string
  email **string
}

func NewUpdate() *UserUpdate{
	return &UserUpdate{}
}
func (u *UserUpdate) ID (id uuid.UUID)  *UserUpdate {
  u.id = &id
  return u
}
func (u *UserUpdate) Name (name string)  *UserUpdate {
  u.name = &name
  return u
}
func (u *UserUpdate) Email (email *string)  *UserUpdate {
  u.email = &email
  return u
}

func (u *UserUpdate) Build()  UserUpdate {
  return *u
}

func ApplySet[B interface {
    Set(column string, value interface{}) B
}] (b B, f UserUpdate) B {
	if f.id != nil {
      b = b.Set(ColumnUserID, *f.id)
    }
	if f.name != nil {
      b = b.Set(ColumnUserName, *f.name)
    }
	if f.email != nil {
      b = b.Set(ColumnUserEmail, *f.email)
    }

  return b
}

type Users []User
func (s Users) ToIDs ()  []uuid.UUID {
	output := make([]uuid.UUID, 0, len(s))
	for _, item := range s {
		output = append(output, item.ID)
	}
	return output
}
func (s Users) ToNames ()  []string {
	output := make([]string, 0, len(s))
	for _, item := range s {
		output = append(output, item.Name)
	}
	return output
}
func (s Users) ToEmails ()  []*string {
	output := make([]*string, 0, len(s))
	for _, item := range s {
		output = append(output, item.Email)
	}
	return output
}

func (s Users) FilterUsers(f func(i User) bool)  Users {
	output := make(Users, 0, len(s))
	for _, item := range s {
		if f(item) {
			output = append(output, item)
		}
	}
	return output
}
