package schedule

import (
	"context"
	"errors"
	"flag"
	commonErrors "kolresource/internal/common/errors"
	"kolresource/internal/email"
	"kolresource/internal/email/domain"
	"kolresource/internal/email/domain/entities"
	"kolresource/internal/email/mock/repositorymock"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestNewEmailSchedule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		interval         time.Duration
		expectedInterval time.Duration
	}{
		{
			name:             "default_interval_when_zero",
			interval:         0,
			expectedInterval: defaultEmailScheduleInterval,
		},
		{
			name:             "default_interval_when_negative",
			interval:         -1 * time.Second,
			expectedInterval: defaultEmailScheduleInterval,
		},
		{
			name:             "custom_interval",
			interval:         5 * time.Minute,
			expectedInterval: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repoMock := repositorymock.NewMockRepository(ctrl)
			emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

			s := NewEmailSchedule(repoMock, emailRepoMock, tt.interval)
			if s == nil {
				t.Fatal("NewEmailSchedule returned nil")
			}

			if s.interval != tt.expectedInterval {
				t.Errorf("interval = %v, want %v", s.interval, tt.expectedInterval)
			}
		})
	}
}

func TestSendEmailJob_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"test body","images":[]}`,
		Status:               email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	emailLog := &entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return([]*entities.EmailJob{emailJob}, nil)
	repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(emailSender, nil)
	repoMock.EXPECT().CountSentEmailsLast24Hours(gomock.Any(), senderID).Return(int64(0), nil)

	// defaultEmailCountPerMinute is 2, so executeJob is called twice.
	// First iteration: processes the pending email log successfully.
	firstWithTx := repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"test body","images":[]}`,
		Status:               email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)
	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(emailLog, nil)
	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	emailRepoMock.EXPECT().SendEmail(gomock.Any(), domain.SendEmailParams{
		Subject: "test",
		Body:    "test body",
		ToEmails: []domain.ToEmail{
			{Email: "kol@example.com", Name: "Test Kol"},
		},
		Images:      []domain.SendEmailImage{},
		SenderName:  "Test Sender",
		SenderEmail: "sender@example.com",
		SenderPwd:   "sender-key",
	}).Return(nil)

	successStatus := email.LogStatusSuccess
	repoMock.EXPECT().UpdateEmailLog(gomock.Any(), domain.UpdateEmailLogParam{
		ID:     logID,
		Status: &successStatus,
		Memo:   "",
	}).Return(nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusFailed).Return(int64(0), nil)

	repoMock.EXPECT().UpdateEmailJob(gomock.Any(), domain.UpdateEmailJobParam{
		JobID:                jobID,
		Status:               email.JobStatusSuccess.ToPointer(),
		IncreaseSuccessCount: 1,
	}).Return(nil)

	// Second iteration: no more pending logs, finalizeJob is called.
	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).After(firstWithTx).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         1,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"test body","images":[]}`,
		Status:               email.JobStatusProcessing,
	}, nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(nil, commonErrors.ErrDataNotFound)
	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusSuccess).Return(int64(1), nil)
	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusFailed).Return(int64(0), nil)
	repoMock.EXPECT().UpdateEmailJob(gomock.Any(), domain.UpdateEmailJobParam{
		JobID:  jobID,
		Status: email.JobStatusSuccess.ToPointer(),
	}).Return(nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err != nil {
		t.Errorf("SendEmailJob() error = %v, want nil", err)
	}
}

func TestSendEmailJob_GrabError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return(nil, errors.New("database error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err == nil {
		t.Fatal("SendEmailJob() expected error, got nil")
	}

	expectedErr := "repo.GrabEmailJob error: database error"
	if err.Error() != expectedErr {
		t.Errorf("SendEmailJob() error = %v, want %v", err, expectedErr)
	}
}

func TestSendEmailJob_NoJobs(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return([]*entities.EmailJob{}, nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err != nil {
		t.Errorf("SendEmailJob() error = %v, want nil", err)
	}
}

func TestSendEmailJob_GetSenderError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:       1,
		SenderID: senderID,
		Status:   email.JobStatusPending,
	}

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return([]*entities.EmailJob{emailJob}, nil)
	repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(nil, errors.New("sender not found"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err != nil {
		t.Errorf("SendEmailJob() error = %v, want nil", err)
	}
}

func TestSendEmailJob_RateLimitExceeded(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:       1,
		SenderID: senderID,
		Status:   email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 10,
	}

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return([]*entities.EmailJob{emailJob}, nil)
	repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(emailSender, nil)
	repoMock.EXPECT().CountSentEmailsLast24Hours(gomock.Any(), senderID).Return(int64(10), nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err != nil {
		t.Errorf("SendEmailJob() error = %v, want nil", err)
	}
}

func TestSendEmailJob_CountSentError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:       1,
		SenderID: senderID,
		Status:   email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().GrabEmailJob(gomock.Any()).Return([]*entities.EmailJob{emailJob}, nil)
	repoMock.EXPECT().GetEmailSenderByID(gomock.Any(), senderID).Return(emailSender, nil)
	repoMock.EXPECT().CountSentEmailsLast24Hours(gomock.Any(), senderID).Return(int64(0), errors.New("count error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.SendEmailJob(context.Background())
	if err != nil {
		t.Errorf("SendEmailJob() error = %v, want nil", err)
	}
}

func TestExecuteJob_FullSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	emailRepoMock.EXPECT().SendEmail(gomock.Any(), domain.SendEmailParams{
		Subject: "test",
		Body:    "body",
		ToEmails: []domain.ToEmail{
			{Email: "kol@example.com", Name: "Test Kol"},
		},
		Images:      []domain.SendEmailImage{},
		SenderName:  "Test Sender",
		SenderEmail: "sender@example.com",
		SenderPwd:   "sender-key",
	}).Return(nil)

	successStatus := email.LogStatusSuccess
	repoMock.EXPECT().UpdateEmailLog(gomock.Any(), domain.UpdateEmailLogParam{
		ID:     logID,
		Status: &successStatus,
		Memo:   "",
	}).Return(nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusFailed).Return(int64(0), nil)

	repoMock.EXPECT().UpdateEmailJob(gomock.Any(), domain.UpdateEmailJobParam{
		JobID:                jobID,
		Status:               email.JobStatusSuccess.ToPointer(),
		IncreaseSuccessCount: 1,
	}).Return(nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err != nil {
		t.Errorf("executeJob() error = %v, want nil", err)
	}
}

func TestExecuteJob_SendEmailFails(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	sendErr := errors.New("SMTP connection failed")
	emailRepoMock.EXPECT().SendEmail(gomock.Any(), domain.SendEmailParams{
		Subject: "test",
		Body:    "body",
		ToEmails: []domain.ToEmail{
			{Email: "kol@example.com", Name: "Test Kol"},
		},
		Images:      []domain.SendEmailImage{},
		SenderName:  "Test Sender",
		SenderEmail: "sender@example.com",
		SenderPwd:   "sender-key",
	}).Return(sendErr)

	failedStatus := email.LogStatusFailed
	repoMock.EXPECT().UpdateEmailLog(gomock.Any(), domain.UpdateEmailLogParam{
		ID:     logID,
		Status: &failedStatus,
		Memo:   "failed to send email: SMTP connection failed",
	}).Return(nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusFailed).Return(int64(1), nil)

	repoMock.EXPECT().UpdateEmailJob(gomock.Any(), domain.UpdateEmailJobParam{
		JobID:  jobID,
		Status: email.JobStatusFailed.ToPointer(),
	}).Return(nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err != nil {
		t.Errorf("executeJob() error = %v, want nil", err)
	}
}

func TestExecuteJob_GetEmailJobByIDForUpdateError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)

	emailJob := &entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Status:   email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(nil, errors.New("db lock error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error, got nil")
	}
}

func TestExecuteJob_GrabPendingEmailLogError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)

	emailJob := &entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Status:   email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Payload:  `{"subject":"test","body":"body","images":[]}`,
		Status:   email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)
	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(nil, errors.New("grab log error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error, got nil")
	}
}

func TestExecuteJob_CountPendingError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Payload:  `{"subject":"test","body":"body","images":[]}`,
		Status:   email.JobStatusProcessing,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Payload:  `{"subject":"test","body":"body","images":[]}`,
		Status:   email.JobStatusProcessing,
	}, nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(0), errors.New("count error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error, got nil")
	}
}

func TestExecuteJob_UpdateEmailLogError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	emailRepoMock.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil)

	successStatus := email.LogStatusSuccess
	repoMock.EXPECT().UpdateEmailLog(gomock.Any(), domain.UpdateEmailLogParam{
		ID:     logID,
		Status: &successStatus,
		Memo:   "",
	}).Return(errors.New("update log error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error, got nil")
	}
}

func TestExecuteJob_UpdateEmailJobError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:                   jobID,
		ExpectedReciverCount: 1,
		SuccessCount:         0,
		SenderID:             senderID,
		Payload:              `{"subject":"test","body":"body","images":[]}`,
		Status:               email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	emailRepoMock.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil)

	successStatus := email.LogStatusSuccess
	repoMock.EXPECT().UpdateEmailLog(gomock.Any(), domain.UpdateEmailLogParam{
		ID:     logID,
		Status: &successStatus,
		Memo:   "",
	}).Return(nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusFailed).Return(int64(0), nil)

	repoMock.EXPECT().UpdateEmailJob(gomock.Any(), gomock.Any()).Return(errors.New("update job error"))

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error, got nil")
	}
}

func TestExecuteJob_InvalidPayload(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)
	logID := int64(100)
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	emailJob := &entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Payload:  `{invalid json`,
		Status:   email.JobStatusPending,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Payload:  `{invalid json`,
		Status:   email.JobStatusPending,
	}, nil)

	repoMock.EXPECT().UpdateEmailJobStats(gomock.Any(), jobID, email.JobStatusProcessing).Return(nil)

	repoMock.EXPECT().GrabPendingEmailLogByJobID(gomock.Any(), jobID).Return(&entities.EmailLog{
		ID:      logID,
		JobID:   jobID,
		KolID:   kolID,
		KolName: "Test Kol",
		Email:   "kol@example.com",
		Status:  email.LogStatusPending,
	}, nil)

	repoMock.EXPECT().CountEmailLogsByJobIDAndStatus(gomock.Any(), jobID, email.LogStatusPending).Return(int64(1), nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err == nil {
		t.Fatal("executeJob() expected error for invalid payload, got nil")
	}
}

func TestExecuteJob_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repositorymock.NewMockRepository(ctrl)
	emailRepoMock := repositorymock.NewMockEmailRepository(ctrl)

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	jobID := int64(1)

	emailJob := &entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Status:   email.JobStatusSuccess,
	}

	emailSender := &entities.EmailSender{
		ID:        senderID,
		Email:     "sender@example.com",
		Name:      "Test Sender",
		Key:       "sender-key",
		RateLimit: 100,
	}

	repoMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetEmailJobByIDForUpdate(gomock.Any(), jobID).Return(&entities.EmailJob{
		ID:       jobID,
		SenderID: senderID,
		Status:   email.JobStatusSuccess,
	}, nil)

	s := NewEmailSchedule(repoMock, emailRepoMock, 1*time.Minute)

	err := s.executeJob(context.Background(), emailJob, emailSender)
	if err != nil {
		t.Errorf("executeJob() error = %v, want nil", err)
	}
}
