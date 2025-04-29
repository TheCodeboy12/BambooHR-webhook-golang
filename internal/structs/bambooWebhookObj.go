package structs

type Employee struct {
	// ChangedFields is a slice of strings listing which fields changed.
	ChangedFields []string `json:"changedFields"`

	// Fields is a map containing all employee attributes as key-value pairs.
	// Both keys (like "Employee #") and values are strings.
	Fields map[string]string `json:"fields"`

	// ID is the unique identifier for the employee, represented as a string.
	ID string `json:"id"`
}

// EmployeeData represents the top-level JSON structure containing a list of employees.
type EmployeeData struct {
	// Employees is a slice containing one or more Employee objects.
	Employees []Employee `json:"employees"`
}
