package email

type EmailJobStatus string

const (
	EmailJobStatusPending          EmailJobStatus = "pending"
	EmailJobStatusProcessing       EmailJobStatus = "processing"
	EmailJobStatusSuccess          EmailJobStatus = "success"
	EmailJobStatusPartiallySuccess EmailJobStatus = "partially_success"
	EmailJobStatusFailed           EmailJobStatus = "failed"
	EmailJobStatusCanceled         EmailJobStatus = "canceled"
)

func (s EmailJobStatus) CanCancel() bool {
	return s == EmailJobStatusPending || s == EmailJobStatusProcessing
}

func (s EmailJobStatus) CanStart() bool {
	return s == EmailJobStatusCanceled
}

func (s EmailJobStatus) ToPointer() *EmailJobStatus {
	return &s
}

type EmailLogStatus string

const (
	EmailLogStatusPending  EmailLogStatus = "pending"
	EmailLogStatusSuccess  EmailLogStatus = "success"
	EmailLogStatusFailed   EmailLogStatus = "failed"
	EmailLogStatusCanceled EmailLogStatus = "canceled"
)

func (s EmailLogStatus) ToPointer() *EmailLogStatus {
	return &s
}
