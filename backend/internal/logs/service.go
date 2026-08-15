package logs

import (
	"database/sql"
	"encoding/json"
	"log"
)

type Service struct {
	db *sql.DB
}

var defaultService *Service

func InitService(db *sql.DB) {
	defaultService = &Service{db: db}
}

func GetService() *Service {
	return defaultService
}

func (s *Service) LogAudit(userID, action, entity, entityID, details string, ipAddress string) {
	if s.db == nil {
		return
	}
	// handle non-uuid gracefully
	userIDVal := sql.NullString{String: userID, Valid: true}
	if len(userID) != 36 { userIDVal.Valid = false }

	_, err := s.db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity, entity_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userIDVal, action, entity, entityID, details, ipAddress)
	if err != nil {
		log.Printf("Failed to insert audit log: %v", err)
	}
}

func (s *Service) LogPayment(txID, provider, endpoint string, reqPayload, resPayload interface{}, statusCode int, errorMsg string) {
	if s.db == nil {
		return
	}
	reqBytes, _ := json.Marshal(reqPayload)
	resBytes, _ := json.Marshal(resPayload)
	_, err := s.db.Exec(`
		INSERT INTO payment_logs (transaction_id, provider, endpoint, request_payload, response_payload, status_code, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, txID, provider, endpoint, string(reqBytes), string(resBytes), statusCode, errorMsg)
	if err != nil {
		log.Printf("Failed to insert payment log: %v", err)
	}
}

func (s *Service) LogError(serviceName, errorType, msg, stackTrace string, context map[string]interface{}) {
	if s.db == nil {
		return
	}
	ctxBytes, _ := json.Marshal(context)
	_, err := s.db.Exec(`
		INSERT INTO app_errors (service, error_type, message, stack_trace, context)
		VALUES ($1, $2, $3, $4, $5)
	`, serviceName, errorType, msg, stackTrace, string(ctxBytes))
	if err != nil {
		log.Printf("Failed to insert app error log: %v", err)
	}
}
