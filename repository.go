package storage

// Repository embeddable struct for database repositories (embedd this in your repository struct)
type Repository struct{}

// Insert data entry to database
func (m *Repository) Insert(tx *Transaction, entry interface{}) (err error) {
	return tx.gormTx.Create(entry).Error
}

// Update some existing data entry
func (m *Repository) Update(tx *Transaction, entry interface{}) (err error) {
	return tx.gormTx.Save(entry).Error
}

// Delete some existing data entry
func (m *Repository) Delete(tx *Transaction, entry interface{}) (err error) {
	return tx.gormTx.Delete(entry).Error
}

// Remove (completely) some existing data entry
func (m *Repository) Remove(tx *Transaction, entry interface{}) (err error) {
	return tx.gormTx.Unscoped().Delete(entry).Error
}
