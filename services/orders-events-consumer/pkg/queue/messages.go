package queue

import (
	"time"
)

type OrderCreatedSQSMessage struct {
	Value         OrderCreatedSQSMessageValue
	receiptHandle *string
}

type OrderCreatedSQSMessageValue struct {
	AwsRegion string `json:"AwsRegion,omitempty"`
	Dynamodb  struct {
		ApproximateCreationDateTime time.Time `json:"ApproximateCreationDateTime,omitempty"`
		Keys                        struct {
			ID struct {
				Value string `json:"Value,omitempty"`
			} `json:"ID,omitempty"`
		} `json:"Keys,omitempty"`
		NewImage struct {
			Date struct {
				Value time.Time `json:"Value,omitempty"`
			} `json:"Date,omitempty"`
			ID struct {
				Value string `json:"Value,omitempty"`
			} `json:"ID,omitempty"`
			OrderedProducts struct {
				Value []struct {
					Value struct {
						ID struct {
							Value string `json:"Value,omitempty"`
						} `json:"ID,omitempty"`
						Name struct {
							Value string `json:"Value,omitempty"`
						} `json:"Name,omitempty"`
						PhotoURL struct {
							Value string `json:"Value,omitempty"`
						} `json:"PhotoURL,omitempty"`
						UnitsOrdered struct {
							Value string `json:"Value,omitempty"`
						} `json:"UnitsOrdered,omitempty"`
					} `json:"Value,omitempty"`
				} `json:"Value,omitempty"`
			} `json:"OrderedProducts,omitempty"`
		} `json:"NewImage,omitempty"`
		OldImage       interface{} `json:"OldImage,omitempty"`
		SequenceNumber string      `json:"SequenceNumber,omitempty"`
		SizeBytes      int         `json:"SizeBytes,omitempty"`
		StreamViewType string      `json:"StreamViewType,omitempty"`
	} `json:"Dynamodb,omitempty"`
	EventID      string      `json:"EventID,omitempty"`
	EventName    string      `json:"EventName,omitempty"`
	EventSource  string      `json:"EventSource,omitempty"`
	EventVersion string      `json:"EventVersion,omitempty"`
	UserIdentity interface{} `json:"UserIdentity,omitempty"`
}
