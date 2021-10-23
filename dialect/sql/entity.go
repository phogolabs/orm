package sql

import "time"

// EntityState is the entity state
type EntityState int

const (
	// EntityStateDetached: the entity is being tracked by the context and
	// exists in the database, but has been marked for deletion from the
	// database the next time SaveChanges is called
	EntityStateDetached EntityState = iota

	// EntityStateCreated: the entity is being tracked by the context but does
	// not yet exist in the database
	EntityStateCreated

	// EntityStateUpdated: the entity is being tracked by the context and
	// exists in the database, and some or all of its property values have been
	// modified
	EntityStateUpdated
)

// Entity represents an entity
type Entity interface {
	// GetUpdatedAt returns the time when the entity was modified
	GetUpdatedAt() time.Time
	// GetCreatedAt returns the time when the entity was created
	GetCreatedAt() time.Time
}

// State returns the entity state
// Accoring to StackOverflow the now() function is transactional
// so in order to find out whether the contact has been changed we should
// just check whether both created_at == updated_at
//
// ref: https://stackoverflow.com/a/49935987
func GetEntityState(entity Entity) EntityState {
	var (
		createdAt = entity.GetCreatedAt()
		updatedAt = entity.GetUpdatedAt()
	)

	if createdAt.IsZero() || updatedAt.IsZero() {
		return EntityStateDetached
	}

	if createdAt.Equal(updatedAt) {
		return EntityStateCreated
	}

	return EntityStateUpdated
}
