package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Employee struct represents the employee details
type Employee struct {
	ID        int     `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	PhoneNo   string  `json:"phoneNo"`
	Role      string  `json:"role"`
	Salary    float64 `json:"salary"`
	Birthdate string  `json:"birthdate"`
}

var csvFilename = "employeesData.csv"

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// containsSubstring checks if a string contains a given substring (case-insensitive)
func containsSubstring(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// getAllEmployees retrieves all employees from the CSV file
func getAllEmployees() ([]Employee, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		return nil, err
	}

	var employees []Employee
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 9 {
			log.Printf("Skipping record %d: insufficient fields", i)
			continue
		}

		id, _ := strconv.Atoi(record[0])
		salary, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			log.Printf("Error parsing salary for employee with ID %d: %v", id, err)
			return nil, fmt.Errorf("error parsing salary for employee with ID %d: %v", id, err)
		}

		employee := Employee{
			ID:        id,
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			Password:  record[4],
			PhoneNo:   record[5],
			Role:      record[6],
			Salary:    salary,
			Birthdate: record[8],
		}
		employees = append(employees, employee)
	}

	log.Printf("Retrieved %d employees from the CSV file", len(employees))
	return employees, nil
}

// createEmployee creates a new employee and stores it in the CSV file
func createEmployee(employee Employee) error {
	// Check if the file exists
	_, err := os.Stat(csvFilename)    //Checking the file information for a file
	fileExists := !os.IsNotExist(err) //checks whether an error occurred during the attempt to get file information
	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error while getting all employees: %v", err)
		return err
	}

	// Open the file with the appropriate flags
	file, err := os.OpenFile(csvFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening file for createEmployee: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// If the file is newly created, write the header
	if !fileExists {
		err := writer.Write([]string{"ID", "FirstName", "LastName", "Email", "Password", "PhoneNo", "Role", "Salary", "Birthdate"})
		if err != nil {
			log.Printf("Error writing header for createEmployee: %v", err)
			return err
		}
	}

	var b bool = true

	for _, emp := range employees {
		if emp.ID == employee.ID {
			log.Printf("ID %d is already present in database.", employee.ID)
			b = false
		}
		// break
	}

	if b == true {
		err = writer.Write([]string{
			strconv.Itoa(employee.ID),
			employee.FirstName,
			employee.LastName,
			employee.Email,
			employee.Password,
			employee.PhoneNo,
			employee.Role,
			strconv.FormatFloat(employee.Salary, 'f', -1, 64),
			employee.Birthdate,
		})
		if err != nil {
			log.Printf("Error writing data for createEmployee: %v", err)
			return err
		}

		log.Printf("Employee with ID %d created successfully", employee.ID)
		return nil
	}else{
		return nil
	}

}

// updateEmployee updates specific fields of an employee with a given ID in the CSV file
func updateEmployee(employeeID int, updatedFields map[string]interface{}) error {
	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error getting all employees for updateEmployee: %v", err)
		return err
	}

	var foundIndex = -1
	for i, employee := range employees {
		if employee.ID == employeeID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("employee with ID %d not found", employeeID)
	}

	// Update specific fields
	for key, value := range updatedFields {
		switch key {
		case "firstName":
			if firstName, ok := value.(string); ok && firstName != "" {
				employees[foundIndex].FirstName = firstName
			}
		case "lastName":
			if lastName, ok := value.(string); ok && lastName != "" {
				employees[foundIndex].LastName = lastName
			}
		case "email":
			if email, ok := value.(string); ok && email != "" {
				employees[foundIndex].Email = email
			}
		case "password":
			if password, ok := value.(string); ok && password != "" {
				employees[foundIndex].Password = password
			}
		case "phoneNo":
			if phoneNo, ok := value.(string); ok && phoneNo != "" {
				employees[foundIndex].PhoneNo = phoneNo
			}
		case "role":
			if role, ok := value.(string); ok && role != "" {
				employees[foundIndex].Role = role
			}
		case "salary":
			if salary, ok := value.(float64); ok {
				employees[foundIndex].Salary = salary
			}
		case "birthdate":
			if birthdate, ok := value.(string); ok && birthdate != "" {
				employees[foundIndex].Birthdate = birthdate
			}
		}
	}

	file, err := os.Create(csvFilename)
	if err != nil {
		log.Printf("Error creating file for updateEmployee: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	err = writer.Write([]string{"ID", "FirstName", "LastName", "Email", "Password", "PhoneNo", "Role", "Salary", "Birthdate"})
	if err != nil {
		log.Printf("Error writing header for updateEmployee: %v", err)
		return err
	}

	// Write the updated employee details
	for _, employee := range employees {
		err := writer.Write([]string{
			strconv.Itoa(employee.ID),
			employee.FirstName,
			employee.LastName,
			employee.Email,
			employee.Password,
			employee.PhoneNo,
			employee.Role,
			strconv.FormatFloat(employee.Salary, 'f', -1, 64),
			employee.Birthdate,
		})
		if err != nil {
			log.Printf("Error writing data for updateEmployee: %v", err)
			return err
		}
	}

	log.Printf("Employee with ID %d updated successfully", employeeID)
	return nil
}

// deleteEmployee deletes an employee with a given ID from the CSV file
func deleteEmployee(employeeID int) error {
	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error getting all employees for deleteEmployee: %v", err)
		return err
	}

	var foundIndex = -1
	for i, employee := range employees {
		if employee.ID == employeeID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("employee with ID %d not found", employeeID)
	}

	// Remove the employee from the slice
	employees = append(employees[:foundIndex], employees[foundIndex+1:]...)

	file, err := os.Create(csvFilename)
	if err != nil {
		log.Printf("Error creating file for deleteEmployee: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	err = writer.Write([]string{"ID", "FirstName", "LastName", "Email", "Password", "PhoneNo", "Role", "Salary", "Birthdate"})
	if err != nil {
		log.Printf("Error writing header for deleteEmployee: %v", err)
		return err
	}

	// Write the remaining employees
	for _, employee := range employees {
		err := writer.Write([]string{
			strconv.Itoa(employee.ID),
			employee.FirstName,
			employee.LastName,
			employee.Email,
			employee.Password,
			employee.PhoneNo,
			employee.Role,
			strconv.FormatFloat(employee.Salary, 'f', -1, 64),
			employee.Birthdate,
		})
		if err != nil {
			log.Printf("Error writing data for deleteEmployee: %v", err)
			return err
		}
	}

	log.Printf("Employee with ID %d deleted successfully", employeeID)
	return nil
}

// handleGetAllEmployees handles the GET request to retrieve all employees
func handleGetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error getting all employees in handleGetAllEmployees: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, employees)
}

// handleCreateEmployee handles the POST request to create a new employee
func handleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		log.Printf("Error decoding JSON in handleCreateEmployee: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || employee.Role == "" {
		log.Println("Missing required fields in handleCreateEmployee")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Perform the actual creation of the employee
	err = createEmployee(employee)
	if err != nil {
		log.Printf("Error creating employee in handleCreateEmployee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log.Println("Employee created successfully")
	// w.WriteHeader(http.StatusCreated)
}

// handleViewEmployeeByID handles the GET request to retrieve details of a specific employee by ID
func handleViewEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid employee ID in handleViewEmployeeByID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error getting all employees in handleViewEmployeeByID: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var foundEmployee *Employee
	for _, employee := range employees {
		if employee.ID == employeeID {
			foundEmployee = &employee
			break
		}
	}

	if foundEmployee == nil {
		log.Println("Employee not found in handleViewEmployeeByID")
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, foundEmployee)
}

// handleSearchEmployees handles the GET request to search employees by firstname, lastname, email, or role
func handleSearchEmployees(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	employees, err := getAllEmployees()
	if err != nil {
		log.Printf("Error getting all employees in handleSearchEmployees: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var matchedEmployees []Employee
	for _, employee := range employees {
		if containsSubstring(employee.FirstName, queryParams.Get("firstName")) &&
			containsSubstring(employee.LastName, queryParams.Get("lastName")) &&
			containsSubstring(employee.Email, queryParams.Get("email")) &&
			containsSubstring(employee.Role, queryParams.Get("role")) {
			matchedEmployees = append(matchedEmployees, employee)
		}
	}

	writeJSONResponse(w, matchedEmployees)
}

// handleUpdateEmployee handles the PUT request to update details of an employee with a given ID
func handleUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid employee ID in handleUpdateEmployee: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var updatedFields map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		log.Printf("Error decoding JSON in handleUpdateEmployee: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received update fields: %v", updatedFields)

	// Perform the actual update of the employee
	err = updateEmployee(employeeID, updatedFields)
	if err != nil {
		log.Printf("Error updating employee in handleUpdateEmployee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleDeleteEmployee handles the DELETE request to delete an employee with a given ID
func handleDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid employee ID in handleDeleteEmployee: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	err = deleteEmployee(employeeID)
	if err != nil {
		log.Printf("Error deleting employee in handleDeleteEmployee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/employees", handleGetAllEmployees).Methods("GET")
	r.HandleFunc("/employees", handleCreateEmployee).Methods("POST")
	r.HandleFunc("/employees/{id}", handleViewEmployeeByID).Methods("GET")
	r.HandleFunc("/employees/search_by_key/search", handleSearchEmployees).Methods("GET")
	r.HandleFunc("/employees/{id}", handleUpdateEmployee).Methods("PATCH")
	r.HandleFunc("/employees/{id}", handleDeleteEmployee).Methods("DELETE")

	log.Println("Server started at port 8080!!!")
	http.ListenAndServe(":8080", r)
}
