package v1

// CreateEmployeeResponseWrapper wraps StandardResponse with CreateEmployeeResponse as data.
// swagger:model CreateEmployeeResponseWrapper
type CreateEmployeeResponseWrapper struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Data      CreateEmployeeResponse `json:"data"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"request_id"`
}

// GetEmployeeByIdResponseWrapper wraps StandardResponse with GetEmployeeByIdResponse as data.
// swagger:model GetEmployeeByIdResponseWrapper
type GetEmployeeByIdResponseWrapper struct {
	Success   bool                    `json:"success"`
	Message   string                  `json:"message"`
	Data      GetEmployeeByIdResponse `json:"data"`
	Timestamp string                  `json:"timestamp"`
	RequestID string                  `json:"request_id"`
}

// GetAllEmployeesResponseWrapper wraps StandardResponse with a list of employees.
// swagger:model GetAllEmployeesResponseWrapper
type GetAllEmployeesResponseWrapper struct {
	Success   bool                      `json:"success"`
	Message   string                    `json:"message"`
	Data      []GetAllEmployeesResponse `json:"data"`
	Timestamp string                    `json:"timestamp"`
	RequestID string                    `json:"request_id"`
}

// UpdateEmployeeResponseWrapper wraps StandardResponse with UpdateEmployeeResponse.
// swagger:model UpdateEmployeeResponseWrapper
type UpdateEmployeeResponseWrapper struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Data      UpdateEmployeeResponse `json:"data"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"request_id"`
}
