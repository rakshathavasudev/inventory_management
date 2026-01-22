package services

import (
    "errors"
    "printflow/models"
)

func Transition(order *models.Order, newStatus string) error {
    validTransitions := map[string][]string{
        models.StatusCreated:         {models.StatusMockupGenerated},
        models.StatusMockupGenerated: {models.StatusApproved},
        models.StatusApproved:        {models.StatusReady},
    }

    allowed := validTransitions[order.Status]
    for _, s := range allowed {
        if s == newStatus {
            order.Status = newStatus
            return nil
        }
    }

    return errors.New("invalid state transition")
}
