// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package entv2

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/entc/integration/migrate/entv2/predicate"
	"github.com/facebookincubator/ent/entc/integration/migrate/entv2/user"
	"github.com/facebookincubator/ent/schema/field"
)

// UserDelete is the builder for deleting a User entity.
type UserDelete struct {
	config
	predicates []predicate.User
}

// Where adds a new predicate to the delete builder.
func (ud *UserDelete) Where(ps ...predicate.User) *UserDelete {
	ud.predicates = append(ud.predicates, ps...)
	return ud
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ud *UserDelete) Exec(ctx context.Context) (int, error) {
	return ud.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (ud *UserDelete) ExecX(ctx context.Context) int {
	n, err := ud.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ud *UserDelete) sqlExec(ctx context.Context) (int, error) {
	spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: user.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	if ps := ud.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, ud.driver, spec)
}

// UserDeleteOne is the builder for deleting a single User entity.
type UserDeleteOne struct {
	ud *UserDelete
}

// Exec executes the deletion query.
func (udo *UserDeleteOne) Exec(ctx context.Context) error {
	n, err := udo.ud.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{user.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (udo *UserDeleteOne) ExecX(ctx context.Context) {
	udo.ud.ExecX(ctx)
}
