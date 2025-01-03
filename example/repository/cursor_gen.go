// Code generated by go generate; DO NOT EDIT.
package repository

import (
	"context"
	"fmt"
	"errors"
	common "github.com/AugustineAurelius/eos/example/common" 


	sq "github.com/Masterminds/squirrel"
)

// builderParams represents optional query parameters.
type builderParams struct {
	OrderBy   *string // Column and direction for sorting (e.g., "id ASC")
	SearchCol *string // Column to search in (e.g., "name")
	SearchTerm *string // Term to search for (e.g., "alice")
	Offset    *int    // Initial offset for pagination
	Limit     *int    // Number of rows to fetch per batch
}

type Cursor struct {
	pool      common.Querier
	builder   sq.SelectBuilder
	rows      []UserModel
	index     int
	offset    int
	limit     int
	params    builderParams
	ctx       context.Context
}

func (r *repository) NewCursor(ctx context.Context, f UserFilter, params builderParams) *Cursor {
	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	b := sq.Select(
		ColumnUserID,
		ColumnUserName,
		ColumnUserEmail,
	).From(TableUser).PlaceholderFormat(sq.Question)

	b = ApplyWhere(b, f)

	return &Cursor{
		pool:      r.db,
		builder:   b,
		rows:      make([]UserModel, 0, limit),
		index:    -1,
		offset:   offset,
		limit:    limit,
		params:   params,
		ctx:      ctx,
	}
}

func (c *Cursor) fetchRows() error {
	builder := c.builder

	if c.params.SearchCol != nil && c.params.SearchTerm != nil {
		builder = builder.Where(sq.ILike{*c.params.SearchCol: fmt.Sprintf("%%%s%%", *c.params.SearchTerm)})
	}

	if c.params.OrderBy != nil {
		builder = builder.OrderBy(*c.params.OrderBy)
	}

	builder = builder.Limit(uint64(c.limit)).Offset(uint64(c.offset))

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	rows, err := c.pool.Query(c.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	c.rows = make([]UserModel, 0,c.limit)
	for rows.Next() {
		var item UserModel
		if err := rows.Scan(&item); err != nil {
			return err
		}
		c.rows = append(c.rows, item)
	}

	if len(c.rows) > 0 {
		c.offset += c.limit
	}
	return nil
}

func (c *Cursor) Next() bool {
	if c.index >= len(c.rows)-1 {
		if err := c.fetchRows(); err != nil || len(c.rows) == 0 {
			return false
		}
		c.index = 0
		return true
	}

	c.index++
	return true
}

func (c *Cursor) Prev() bool {
	if c.index > 0 {
		c.index--
		return true
	}
	return false
}

func (c *Cursor) Current() (UserModel, error) {
	if c.index >= 0 && c.index < len(c.rows) {
		return c.rows[c.index], nil
	}
	return UserModel{} , errors.New("not found")
}

func (c *Cursor) Reset() {
	c.index = -1
	c.offset = 0
	c.rows = make([]UserModel, 0)
}
