package helpersaudit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Convert() (string, any, any, error) {
	return "", nil, nil, fmt.Errorf("not implemented for audit traces")
}

func Delete() (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}

func Create() (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}

func Update() (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}

func EventuallyRelativeDateTime(s string) (*timestamppb.Timestamp, error) {
	if s == "" {
		return nil, nil
	}
	duration, err := time.ParseDuration(s)
	if err == nil {
		return timestamppb.New(time.Now().Add(-duration)), nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return timestamppb.Now(), fmt.Errorf("failed to convert time: %w", err)
	}
	return timestamppb.New(t), nil
}

func ToPhase(phase string) *apiv1.AuditPhase {
	p, ok := apiv1.AuditPhase_value[phase]
	if !ok {
		return nil
	}

	return pointer.Pointer(apiv1.AuditPhase(p))
}

func TryPrettifyBody(trace *apiv1.AuditTrace) *apiv1.AuditTrace {
	if trace.Body != nil {
		trimmed := strings.Trim(*trace.Body, `"`)
		body := map[string]any{}
		if err := json.Unmarshal([]byte(trimmed), &body); err == nil {
			if pretty, err := json.MarshalIndent(body, "", "    "); err == nil {
				trace.Body = pointer.Pointer(string(pretty))
			}
		}
	}

	return trace
}
