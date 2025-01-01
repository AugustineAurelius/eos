//Code generated by generator, DO NOT EDIT.
package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	common "github.com/AugustineAurelius/eos/example/common"
)


type repository struct{
	db common.Database 
}

func New(db common.Database) *repository {
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
  email *string
  emails []string
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
func (f *UserFilter) Email (email string)  *UserFilter {
  f.email = &email
  return f
}

func (f *UserFilter) AddOneToEmails (email string)  *UserFilter {
  f.emails = append(f.emails, email)
  return f
}

func (f *UserFilter) Emails (emails []string)  *UserFilter {
  f.emails =  emails
  return f
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
  email *string
}
func (f *UserUpdate) ID (id uuid.UUID)  *UserUpdate {
  f.id = &id
  return f
}
func (f *UserUpdate) Name (name string)  *UserUpdate {
  f.name = &name
  return f
}
func (f *UserUpdate) Email (email string)  *UserUpdate {
  f.email = &email
  return f
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
