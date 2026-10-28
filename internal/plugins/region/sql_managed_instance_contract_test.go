package region

import "testing"

// Reuse the preserved explicit JSON-tagged strict wire runner; its decoder rejects
// unknown fields/trailing JSON and checks complete collections/catalogs/ownership.
func TestSQLManagedInstanceFrozenManifest(t *testing.T) {
	sqlRunFrozenManifest(t, "sql-managed-instance-capability-contract-v1.json", "sql-managed-instance-capability-contract-v1", "microsoft.sql/managedinstances", 163, 134, NewSQLManagedInstanceSKUCollector)
}
