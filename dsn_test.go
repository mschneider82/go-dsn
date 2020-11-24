// Package dsn contains the utilities used for dsn message (DSN) generation.
//
// It implements RFC 3464 and RFC 3462.
package dsn

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-message/textproto"
	"github.com/emersion/go-smtp"
)

func TestGenerateDSN(t *testing.T) {
	type args struct {
		utf8         bool
		envelope     Envelope
		mtaInfo      ReportingMTAInfo
		rcptsInfo    []RecipientInfo
		failedHeader textproto.Header
	}
	tests := []struct {
		name                  string
		args                  args
		want                  textproto.Header
		wantStringInOutWriter string
		wantErr               bool
	}{
		{
			name: "test1",
			args: args{
				utf8: false,
				envelope: Envelope{
					MsgID: "msgid1",
					From:  "from@example.com",
					To:    "to@example.com",
				},
				mtaInfo: ReportingMTAInfo{
					ReportingMTA:    "reportingmta.example.com",
					ReceivedFromMTA: "receivedmta.example.com",
					XSender:         "XSender@example.com",
					XMessageID:      "XMessageID",
					ArrivalDate:     time.Date(2020, 01, 02, 15, 04, 05, 06, time.UTC),
					LastAttemptDate: time.Date(2020, 01, 02, 15, 04, 05, 07, time.UTC),
				},
				rcptsInfo: []RecipientInfo{{
					FinalRecipient: "finalrcpt@example.com",
					RemoteMTA:      "remotemta.example.com",
					Action:         ActionDelivered,
					Status:         smtp.EnhancedCode{2, 0, 0},
					DiagnosticCode: nil,
				}},
				failedHeader: textproto.Header{},
			},
			want:                  textproto.Header{},
			wantStringInOutWriter: "This is the mail delivery system at reportingmta.example.com",
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWriter := &bytes.Buffer{}
			got, err := GenerateDSN(tt.args.utf8, tt.args.envelope, tt.args.mtaInfo, tt.args.rcptsInfo, tt.args.failedHeader, outWriter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Get("Message-Id") != tt.args.envelope.MsgID {
				t.Errorf("Message Id is: %s, want %v", got.Get("Message-Id"), tt.want)
			}
			if !strings.ContainsAny(outWriter.String(), tt.wantStringInOutWriter) {
				t.Errorf("outWriter should contain %s, but got: %s", tt.wantStringInOutWriter, outWriter.String())
			}
		})
	}
}
