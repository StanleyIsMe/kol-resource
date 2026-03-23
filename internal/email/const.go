package email

type JobStatus string

const (
	JobStatusPending          JobStatus = "pending"
	JobStatusProcessing       JobStatus = "processing"
	JobStatusSuccess          JobStatus = "success"
	JobStatusPartiallySuccess JobStatus = "partially_success"
	JobStatusFailed           JobStatus = "failed"
	JobStatusCanceled         JobStatus = "canceled"
)

func (s JobStatus) CanCancel() bool {
	return s == JobStatusPending || s == JobStatusProcessing
}

func (s JobStatus) CanStart() bool {
	return s == JobStatusCanceled
}

func (s JobStatus) ToPointer() *JobStatus {
	return &s
}

type LogStatus string

const (
	LogStatusPending  LogStatus = "pending"
	LogStatusSuccess  LogStatus = "success"
	LogStatusFailed   LogStatus = "failed"
	LogStatusCanceled LogStatus = "canceled"
)

func (s LogStatus) ToPointer() *LogStatus {
	return &s
}
