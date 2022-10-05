package filestorage

import (
	"api/internal/constants"
	customerrors "api/internal/customerrors"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"sort"

	"api/data"
)

type EmailSubscriptionRepository struct {
	filepath string
	logger   *logrus.Logger
}

func NewEmailSubscriptionRepository(filepath string, logger *logrus.Logger) *EmailSubscriptionRepository {
	return &EmailSubscriptionRepository{
		filepath: filepath,
		logger:   logger,
	}
}

func (r *EmailSubscriptionRepository) Add(email string) error {
	emails, err := r.GetAll()
	if err != nil {
		return err
	}

	emails, err = r.addToSorted(emails, email)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	records := data.SubscribedEmails{
		Emails: emails,
	}

	updatedData, err := json.Marshal(records)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.filepath, updatedData, constants.WriteFilePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *EmailSubscriptionRepository) GetAll() ([]string, error) {
	records := data.SubscribedEmails{}
	file, err := os.ReadFile(r.filepath)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	err = json.Unmarshal(file, &records)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	return records.Emails, nil
}

func (r *EmailSubscriptionRepository) addToSorted(sourceSlice []string, itemToAdd string) ([]string, error) {
	index := sort.SearchStrings(sourceSlice, itemToAdd)
	if index != len(sourceSlice) {
		if sourceSlice[index] == itemToAdd {
			return nil, customerrors.ErrEmailDuplicate
		}
	}

	sourceSlice = append(sourceSlice, "")
	copy(sourceSlice[index+1:], sourceSlice[index:])
	sourceSlice[index] = itemToAdd

	return sourceSlice, nil
}
