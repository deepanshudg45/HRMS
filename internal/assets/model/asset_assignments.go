package model

type ReturnRequest struct {
	ReturnedOn        string `json:"returnedOn"`
	ConditionAtReturn string `json:"conditionAtReturn"`
	ReturnReason      string `json:"returnReason"`
}

type ReturnDTO struct {
	AssetID           string `json:"assetId"`
	AssignmentID      string `json:"assignmentId"`
	ConditionAtReturn string `json:"conditionAtReturn"`
	ReturnReason      string `json:"returnReason"`
	ReturnedOn        string `json:"returnedOn"`
	AssetStatus       string `json:"assetStatus"`
}

type AssignmentDTO struct {
	ID                    string `json:"id"`
	EmployeeID            string `json:"employeeId"`
	AcknowledgementStatus string `json:"acknowledgementStatus"`
	AcknowledgedAt        string `json:"acknowledgedAt"`
}

type AssignmentHistoryDTO struct {
	ID         string `json:"id"`
	AssetID    string `json:"assetId"`
	EmployeeID string `json:"employeeId"`
	IsActive   bool   `json:"isActive"`
	AssignedOn string `json:"assignedOn"`
}
