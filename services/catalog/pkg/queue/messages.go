package queue

type OrderCreatedSQSMessage struct {
	Value         OrderCreatedSQSMessageValue
	receiptHandle *string
}

type OrderCreatedSQSBody struct {
	Records []OrderCreatedSQSMessageValue `json:"Records"`
}

type OrderCreatedSQSMessageValue struct {
	AwsRegion string `json:"awsRegion,omitempty"`
	Dynamodb  struct {
		// ApproximateCreationDateTime time.Time `json:"ApproximateCreationDateTime,omitempty"`
		Keys struct {
			ID struct {
				Value string `json:"Value,omitempty"`
			} `json:"ID,omitempty"`
		} `json:"Keys,omitempty"`
		NewImage struct {
			// Date struct {
			// 	Value time.Time `json:"Value,omitempty"`
			// } `json:"Date,omitempty"`
			ID struct {
				Value string `json:"S,omitempty"`
			} `json:"ID,omitempty"`
			OrderedProducts struct {
				Value []struct {
					Value struct {
						ID struct {
							Value string `json:"S,omitempty"`
						} `json:"ID,omitempty"`
						Name struct {
							Value string `json:"S,omitempty"`
						} `json:"Name,omitempty"`
						PhotoURL struct {
							Value string `json:"S,omitempty"`
						} `json:"PhotoURL,omitempty"`
						UnitsOrdered struct {
							Value string `json:"N,omitempty"`
						} `json:"UnitsOrdered,omitempty"`
					} `json:"M,omitempty"`
				} `json:"L,omitempty"`
			} `json:"OrderedProducts,omitempty"`
		} `json:"NewImage,omitempty"`
		OldImage       interface{} `json:"OldImage,omitempty"`
		SequenceNumber string      `json:"SequenceNumber,omitempty"`
		SizeBytes      int         `json:"SizeBytes,omitempty"`
		StreamViewType string      `json:"StreamViewType,omitempty"`
	} `json:"dynamodb,omitempty"`
	EventID        string      `json:"eventID,omitempty"`
	EventName      string      `json:"eventName,omitempty"`
	EventSource    string      `json:"eventSource,omitempty"`
	EventSourceArn string      `json:"eventSourceArn,omitempty"`
	EventVersion   string      `json:"eventVersion,omitempty"`
	UserIdentity   interface{} `json:"userIdentity,omitempty"`
}
