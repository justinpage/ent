---
id: schema-edges
title: Edges
---

## Quick Summary

Edges are the relations (or associations) of entities. For example, user's pets, or group's users.


![er-group-users](https://entgo.io/assets/er_user_pets_groups.png)

In the example above, you can see 2 relations declared using edges. Let's go over them.

1\. `pets` / `owner` edges; user's pets and pet's owner - 

`ent/schema/user.go`
```go
package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
)

// User schema.
type User struct {
	ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// ...
	}
}

// Edges of the user.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}
```


`ent/schema/pet.go`
```go
package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
)

// User schema.
type Pet struct {
	ent.Schema
}

// Fields of the user.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		// ...
	}
}

// Edges of the user.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique(),
	}
}
```

As you can see, a `User` entity can have many pets, but a `Pet` entity can have only one owner.  
In relationship definition, the `pets` edge is a *O2M* (one-to-many) relationship, and the `owner` edge
is a *M2O* (many-to-one) relationship.

The `User` schema **owns** the `pets/owner` relationship because it uses `edge.To`, and the `Pet` schema
just have a back-reference to it, declared using `edge.From` with the `Ref` method.

The `Ref` method describes which edge of the `User` schema we're referencing to, because, there can be multiple
references from one schema to other. 

The cardinality of the edge/relationship can be controlled using the `Unique` method, and it's explained
more widely below. 

2\. `users` / `groups` edges; group's users and user's groups - 

`ent/schema/group.go`
```go
package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
)

// Group schema.
type Group struct {
	ent.Schema
}

// Fields of the group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		// ...
	}
}

// Edges of the group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
```

`ent/schema/user.go`
```go
package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
)

// User schema.
type User struct {
	ent.Schema
}

// Fields of the user.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// ...
	}
}

// Edges of the user.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).
			Ref("users"),
		// "pets" declared in the example above.
		edge.To("pets", Pet.Type),
	}
}
```

As you can see, a Group entity can have many users, and a User entity can have have many groups.  
In relationship definition, the `users` edge is a *M2M* (many-to-many) relationship, and the `groups`
edge is also a *M2M* (many-to-many) relationship.

## To and From

`edge.To` and `edge.From` are the 2 builders for creating edges/relations.

A schema that defines an edge using the `edge.To` builder is owning the relation,
unlike using the `edge.From` builder that gives only a reference for the relation (with different name).

Let's go over a few examples, that show how to define different relation types using edges.

## Relationship

The following examples:

- [O2O Two Types](#o2o-two-types)
- [O2O Same Type](#o2o-same-type)
- [O2O Bidirectional](#o2o-bidirectional)
- [O2M Two Types](#o2m-two-types)
- [O2M Same Type](#o2m-same-type)

## O2O Two Types

![er-user-card](https://entgo.io/assets/er_user_card.png)

In this example, a user **has only one** credit-card, and a card **has only one** owner.

The `User` schema defines an `edge.To` card named `card`, and the `Card` schema
defines a reference to this edge using `edge.From` named `owner`. 


`ent/schema/user.go`
```go
// Edges of the user.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("card", Card.Type).
			Unique(),
	}
}
```

`ent/schema/card.go`
```go
// Edges of the user.
func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("card").
			Unique().
			// We add the "Required" method to the builder
			// to make this edge required on entity creation.
			// i.e. Card cannot be created without its owner.
			Required(),
	}
}
```

The API for interacting with these edges is as follows:
```go
func Do(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Mashraki").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}
	log.Println("user:", a8m)
	card1, err := client.Card.
		Create().
		SetOwner(a8m).
		SetNumber("1020").
		SetExpired(time.Now().Add(time.Minute)).
		Save(ctx)
	if err != nil {
    	return fmt.Errorf("creating card: %v", err)
    }
	log.Println("card:", card1)
	// Only returns the card of the user,
	// and expects that there's only one.
	card2, err := a8m.QueryCard().Only(ctx)
	if err != nil {
		return fmt.Errorf("querying card: %v", err)
    }
	log.Println("card:", card2)
	// The Card entity is able to query its owner using
	// its back-reference.
	owner, err := card2.QueryOwner().Only(ctx)
	if err != nil {
		return fmt.Errorf("querying owner: %v", err)
    }
	log.Println("owner:", owner)
	return nil
}
```

The full example exists in [GitHub](https://github.com/facebookincubator/ent/tree/master/examples/o2o2types).

## O2O Same Type

![er-linked-list](https://entgo.io/assets/er_linked_list.png)

In this linked-list example, we have a **recursive relation** named `next`/`prev`. Each node in the list can
have only of `next`. If a node A points (using `next`) to a node B, B can get its pointer using `prev`.   

`ent/schema/node.go`
```go
// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("next", Node.Type).
			Unique().
			From("prev").
			Unique(),
	}
}
```

As you can see, in cases of relations of the same type, you can declare the edge and its
reference in the same builder.

```diff
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
+		edge.To("next", Node.Type).
+			Unique().
+			From("prev").
+			Unique(),

-		edge.To("next", Node.Type).
-			Unique(),
-		edge.From("prev", Node.Type).
-			Ref("next).
-			Unique(),
	}
}
```

The API for interacting with these edges is as follows:

```go
func Do(ctx context.Context, client *ent.Client) error {
	head, err := client.Node.
		Create().
		SetValue(1).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating the head: %v", err)
	}
	curr := head
	// Generate the following linked-list: 1<->2<->3<->4<->5.
	for i := 0; i < 4; i++ {
		curr, err = client.Node.
			Create().
			SetValue(curr.Value + 1).
			SetPrev(curr).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	// Loop over the list and print it. `FirstX` panics if an error occur.
	for curr = head; curr != nil; curr = curr.QueryNext().FirstX(ctx) {
		fmt.Printf("%d ", curr.Value)
	}
	// Output: 1 2 3 4 5

	// Make the linked-list circular:
	// The tail of the list, has no "next".
	tail, err := client.Node.
		Query().
		Where(node.Not(node.HasNext())).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("getting the tail of the list: %v", tail)
	}
	tail, err = tail.Update().SetNext(head).Save(ctx)
	if err != nil {
		return err
	}
	// Check that the change actually applied:
	prev, err := head.QueryPrev().Only(ctx)
	if err != nil {
		return fmt.Errorf("getting head's prev: %v", err)
	}
	fmt.Printf("\n%v", prev.Value == tail.Value)
	// Output: true
	return nil
}
```

The full example exists in [GitHub](https://github.com/facebookincubator/ent/tree/master/examples/o2o2recur).

## O2O Bidirectional

![er-user-spouse](https://entgo.io/assets/er_user_spouse.png)

In this user-spouse example, we have a **reflexive relation** named `spouse`. Each user can have only one spouse.
If user A sets its spouse (using `spouse`) to B, B can get its spouse using the `spouse` edge.

Note that, there's no owner/inverse terms in cases of bidirectional edges.

`ent/schema/user.go`
```go
// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("spouse", User.Type).
			Unique(),
	}
}
```

The API for interacting with this edge is as follows:

```go
func Do(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}
	nati, err := client.User.
		Create().
		SetAge(28).
		SetName("nati").
		SetSpouse(a8m).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}

	// Query the spouse edge.
	// Unlike `Only`, `OnlyX` panics if an error occurs.
	spouse := nati.QuerySpouse().OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: a8m

	spouse = a8m.QuerySpouse().OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: nati

	// Query how many users have a spouse.
	// Unlike `Count`, `CountX` panics if an error occurs.
	count := client.User.
		Query().
		Where(user.HasSpouse()).
		CountX(ctx)
	fmt.Println(count)
	// Output: 2

	// Get the user, that has a spouse with name="a8m".
	spouse = client.User.
		Query().
		Where(user.HasSpouseWith(user.Name("a8m"))).
		OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: nati
	return nil
}
```

The full example exists in [GitHub](https://github.com/facebookincubator/ent/tree/master/examples/o2obidi).

## O2M Two Types

![er-user-pets](https://entgo.io/assets/er_user_pets.png)

In this user-pets example, we have a O2M relation between user and its pets.
Each user **has many** pets, and a pet **has one** owner.
If user A adds a pet B using the `pets` edge, B can get its owner using the `owner` edge.

Note that, this relation is also a M2O (many-to-one) from the point of view of the `Pet` schema. 

`ent/schema/user.go`
```go
// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}
```

`ent/schema/pet.go`
```go
// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique(),
	}
}
```

The API for interacting with these edges is as follows:

```go
func Do(ctx context.Context, client *ent.Client) error {
	// Create the 2 pets.
	pedro, err := client.Pet.
		Create().
		SetName("pedro").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating pet: %v", err)
	}
	lola, err := client.Pet.
		Create().
		SetName("lola").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating pet: %v", err)
	}
	// Create the user, and add its pets on the creation.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddPets(pedro, lola).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}
	fmt.Println("User created:", a8m)
	// Output: User(id=1, age=30, name=a8m)

	// Query the owner. Unlike `Only`, `OnlyX` panics if an error occurs.
	owner := pedro.QueryOwner().OnlyX(ctx)
	fmt.Println(owner.Name)
	// Output: a8m

	// Traverse the sub-graph. Unlike `Count`, `CountX` panics if an error occurs.
	count := pedro.
		QueryOwner(). // a8m
		QueryPets().  // pedro, lola
		CountX(ctx)   // count
	fmt.Println(count)
	// Output: 2
	return nil
}
```
The full example exists in [GitHub](https://github.com/facebookincubator/ent/tree/master/examples/o2m2types).

## O2M Same Type

![er-tree](https://entgo.io/assets/er_tree.png)

In this example, we have a recursive O2M relation between tree's nodes and their children (or their parent).  
Each node in the tree **has many** children, and **has one** parent. If node A adds B to its children,
B can get its owner using the `owner` edge.


`ent/schema/node.go`
```go
// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Node.Type).
			From("parent").
			Unique(),
	}
}
```

As you can see, in cases of relations of the same type, you can declare the edge and its
reference in the same builder.

```diff
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
+		edge.To("children", Node.Type).
+			From("parent").
+			Unique(),

-		edge.To("children", Node.Type),
-		edge.From("parent", Node.Type).
-			Ref("children).
-			Unique(),
	}
}
```

The API for interacting with these edges is as follows:

```go
func Do(ctx context.Context, client *ent.Client) error {
	root, err := client.Node.
		Create().
		SetValue(2).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating the root: %v", err)
	}
	// Add additional nodes to the tree:
	//
	//       2
	//     /   \
	//    1     4
	//        /   \
	//       3     5
	//
	// Unlike `Create`, `CreateX` panics if an error occurs.
	n1 := client.Node.
		Create().
		SetValue(1).
		SetParent(root).
		SaveX(ctx)
	n4 := client.Node.
		Create().
		SetValue(4).
		SetParent(root).
		SaveX(ctx)
	n3 := client.Node.
		Create().
		SetValue(3).
		SetParent(n4).
		SaveX(ctx)
	n5 := client.Node.
		Create().
		SetValue(5).
		SetParent(n4).
		SaveX(ctx)

	fmt.Println("Tree leafs", []int{n1.Value, n3.Value, n5.Value})
	// Output: Tree leafs [1 3 5]

	// Get all leafs (nodes without children).
	// Unlike `Int`, `IntX` panics if an error occurs.
	ints := client.Node.
		Query().                             // All nodes.
		Where(node.Not(node.HasChildren())). // Only leafs.
		Order(ent.Asc(node.FieldValue)).     // Order by their `value` field.
		GroupBy(node.FieldValue).            // Extract only the `value` field.
		IntsX(ctx)
	fmt.Println(ints)
	// Output: [1 3 5]

	// Get orphan nodes (nodes without parent).
	// Unlike `Only`, `OnlyX` panics if an error occurs.
	orphan := client.Node.
		Query().
		Where(node.Not(node.HasParent())).
		OnlyX(ctx)
	fmt.Println(orphan)
	// Output: Node(id=1, value=2)

	return nil
}
```

The full example exists in [GitHub](https://github.com/facebookincubator/ent/tree/master/examples/o2mrecur).


## M2M Two Types

## M2M Same Type

## M2M Bidirectional

## Required

## Indexes