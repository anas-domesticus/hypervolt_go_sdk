package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScheduleSession_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		ss      *ScheduleSession
		args    args
		wantErr bool
	}{
		{
			name: "Valid Data",
			ss:   &ScheduleSession{},
			args: args{
				data: []byte(`{
					"days": ["Monday", "Wednesday"],
					"end_time": "18:00",
					"mode": "Boost",
					"session_type": "recurring",
					"start_time": "09:00"
				}`),
			},
			wantErr: false,
		},
		{
			name: "Invalid Time Format",
			ss:   &ScheduleSession{},
			args: args{
				data: []byte(`{
					"days": ["Tuesday", "Thursday"],
					"end_time": "invalid time",
					"mode": "Boost",
					"session_type": "recurring",
					"start_time": "09:00"
				}`),
			},
			wantErr: true,
		},
		{
			name: "Empty JSON",
			ss:   &ScheduleSession{},
			args: args{
				data: []byte(`""`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ss.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ScheduleSession.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScheduleSession_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  ScheduleSession
		want    string
		wantErr bool
	}{
		{
			name: "All Fields Populated",
			fields: ScheduleSession{
				Days:        []DayOfWeek{MONDAY, TUESDAY},
				EndTime:     time.Date(1, 1, 1, 14, 0, 0, 0, time.UTC),
				Mode:        ECO,
				SessionType: RECURRING,
				StartTime:   time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			want: `{"days":["MONDAY","TUESDAY"],"end_time":"14:00","mode":"eco","session_type":"recurring","start_time":"12:00"}`,
		},
		{
			name:   "Empty ScheduleSession",
			fields: ScheduleSession{},
			want:   `{"days":[],"end_time":"00:00","mode":"","session_type":"","start_time":"00:00"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := ScheduleSession{
				Days:        tt.fields.Days,
				EndTime:     tt.fields.EndTime,
				Mode:        tt.fields.Mode,
				SessionType: tt.fields.SessionType,
				StartTime:   tt.fields.StartTime,
			}

			got, err := ss.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.JSONEq(t, tt.want, string(got))
		})
	}
}
