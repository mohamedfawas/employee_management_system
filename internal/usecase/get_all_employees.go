package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mohamedfawas/employee_management_system/internal/domain/entity"
)

const (
	employeesListKey    = "employees:list"
	employeesListTTL    = 5 * time.Minute
	maxEmployeesToCache = 100
)

func (u *employeeUsecaseImpl) GetAllEmployees(ctx context.Context) ([]*entity.Employee, error) {
	// 1) Try cache
	cached, err := u.cache.Get(ctx, employeesListKey)
	if err == nil && cached != "" {
		var fromCache []*entity.Employee

		// json.Unmarshal converts the JSON string from Redis into Go structs.
		//
		// Example:
		//   cached = `[{"id":1,"name":"Alice"}, {"id":2,"name":"Bob"}]`
		// After Unmarshal:
		//   fromCache = []*Employee{
		//       {ID:1, Name:"Alice"},
		//       {ID:2, Name:"Bob"},
		//   }
		if umErr := json.Unmarshal([]byte(cached), &fromCache); umErr == nil {
			log.Println("[CACHE HIT] employees:list returned from Redis")
			return fromCache, nil
		} else {
			// unmarshalling failed — log and fall back to DB
			log.Printf("[CACHE ERROR] failed to unmarshal employees:list: %v", umErr)
		}
	} else {
		// Key not found or redis error → cache miss
		log.Println("[CACHE MISS] employees:list returning from DB")
	}

	// 2) Cache miss or error , then fetch from DB
	employees, err := u.employeeRepository.GetAllEmployees(ctx)
	if err != nil {
		return nil, err
	}

	// 3) Limit cached list size
	toCache := employees
	if len(toCache) > maxEmployeesToCache {
		toCache = toCache[:maxEmployeesToCache]
	}

	// 4) SAVE FRESH DATA INTO CACHE
	if b, mErr := json.Marshal(toCache); mErr == nil {
		// After this call, Redis will store:
		//   "employees:list" → "[{...employee1...}, {...employee2...}, ...]"
		//
		if cerr := u.cache.Set(ctx, employeesListKey, string(b), employeesListTTL); cerr != nil {
			log.Printf("cache: failed to set employees list: %v", cerr)
		}
	} else {
		log.Printf("cache: marshal error: %v", mErr)
	}

	return employees, nil
}
