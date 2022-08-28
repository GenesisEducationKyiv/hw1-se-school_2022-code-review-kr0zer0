package repository

import (
	"encoding/json"
	"os"

	"api/data"
)

type EmailSubscriptionRepository struct {
	filepath string
}

func NewEmailSubscriptionRepository(filepath string) *EmailSubscriptionRepository {
	return &EmailSubscriptionRepository{
		filepath: filepath,
	}
}

const WriteFilePerm = 0o600

func (r *EmailSubscriptionRepository) Add(email string) error {
	emails, err := r.GetAll()
	if err != nil {
		return err
	}

	emails = append(emails, email)

	records := data.Data{
		Emails: emails,
	}

	updatedData, err := json.Marshal(records)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.filepath, updatedData, WriteFilePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *EmailSubscriptionRepository) CheckIfExists(emailToFind string) (bool, error) {
	records, err := r.GetAll()
	if err != nil {
		return false, err
	}

	for _, email := range records {
		if email == emailToFind {
			return true, nil
		}
	}

	return false, nil
}

func (r *EmailSubscriptionRepository) GetAll() ([]string, error) {
	records := data.Data{}
	file, err := os.ReadFile(r.filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &records)
	if err != nil {
		return nil, err
	}
	return records.Emails, nil
}
