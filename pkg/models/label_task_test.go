package models

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/user"
	"gopkg.in/d4l3k/messagediff.v1"

	"code.vikunja.io/web"
)

func TestLabelTask_ReadAll(t *testing.T) {
	type fields struct {
		ID       int64
		TaskID   int64
		LabelID  int64
		Created  time.Time
		CRUDable web.CRUDable
		Rights   web.Rights
	}
	type args struct {
		search string
		a      web.Auth
		page   int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLabels interface{}
		wantErr    bool
		errType    func(error) bool
	}{
		{
			name: "normal",
			fields: fields{
				TaskID: 1,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantLabels: []*labelWithTaskID{
				{
					TaskID: 1,
					Label: Label{
						ID:          4,
						Title:       "Label #4 - visible via other task",
						Created:     testCreatedTime,
						Updated:     testUpdatedTime,
						CreatedByID: 2,
						CreatedBy: &user.User{
							ID:                    2,
							Username:              "user2",
							Password:              "$2a$14$dcadBoMBL9jQoOcZK8Fju.cy0Ptx2oZECkKLnaa8ekRoTFe1w7To.",
							Issuer:                "local",
							EmailRemindersEnabled: true,
							Created:               testCreatedTime,
							Updated:               testUpdatedTime,
						},
					},
				},
			},
		},
		{
			name: "no right to see the task",
			fields: fields{
				TaskID: 14,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantErr: true,
			errType: IsErrNoRightToSeeTask,
		},
		{
			name: "nonexistant task",
			fields: fields{
				TaskID: 9999,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantErr: true,
			errType: IsErrTaskDoesNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.LoadAndAssertFixtures(t)

			l := &LabelTask{
				ID:       tt.fields.ID,
				TaskID:   tt.fields.TaskID,
				LabelID:  tt.fields.LabelID,
				Created:  tt.fields.Created,
				CRUDable: tt.fields.CRUDable,
				Rights:   tt.fields.Rights,
			}
			gotLabels, _, _, err := l.ReadAll(tt.args.a, tt.args.search, tt.args.page, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("LabelTask.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("LabelTask.ReadAll() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
			if diff, equal := messagediff.PrettyDiff(gotLabels, tt.wantLabels); !equal {
				t.Errorf("LabelTask.ReadAll() = %v, want %v, diff: %v", l, tt.wantLabels, diff)
			}
		})
	}
}

func TestLabelTask_Create(t *testing.T) {
	type fields struct {
		ID       int64
		TaskID   int64
		LabelID  int64
		Created  time.Time
		CRUDable web.CRUDable
		Rights   web.Rights
	}
	type args struct {
		a web.Auth
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		errType       func(error) bool
		wantForbidden bool
	}{
		{
			name: "normal",
			fields: fields{
				TaskID:  1,
				LabelID: 1,
			},
			args: args{
				a: &user.User{ID: 1},
			},
		},
		{
			name: "already existing",
			fields: fields{
				TaskID:  1,
				LabelID: 4,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantErr: true,
			errType: IsErrLabelIsAlreadyOnTask,
		},
		{
			name: "nonexisting label",
			fields: fields{
				TaskID:  1,
				LabelID: 9999,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantForbidden: true,
		},
		{
			name: "nonexisting task",
			fields: fields{
				TaskID:  9999,
				LabelID: 1,
			},
			args: args{
				a: &user.User{ID: 1},
			},
			wantForbidden: true,
			wantErr:       true,
			errType:       IsErrTaskDoesNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.LoadAndAssertFixtures(t)

			l := &LabelTask{
				ID:       tt.fields.ID,
				TaskID:   tt.fields.TaskID,
				LabelID:  tt.fields.LabelID,
				Created:  tt.fields.Created,
				CRUDable: tt.fields.CRUDable,
				Rights:   tt.fields.Rights,
			}
			allowed, _ := l.CanCreate(tt.args.a)
			if !allowed && !tt.wantForbidden {
				t.Errorf("LabelTask.CanCreate() forbidden, want %v", tt.wantForbidden)
			}
			err := l.Create(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("LabelTask.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("LabelTask.Create() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
			if !tt.wantErr {
				db.AssertExists(t, "label_task", map[string]interface{}{
					"id":       l.ID,
					"task_id":  l.TaskID,
					"label_id": l.LabelID,
				}, false)
			}
		})
	}
}

func TestLabelTask_Delete(t *testing.T) {
	type fields struct {
		ID       int64
		TaskID   int64
		LabelID  int64
		Created  time.Time
		CRUDable web.CRUDable
		Rights   web.Rights
	}
	tests := []struct {
		name          string
		fields        fields
		wantErr       bool
		errType       func(error) bool
		auth          web.Auth
		wantForbidden bool
	}{
		{
			name: "normal",
			fields: fields{
				TaskID:  1,
				LabelID: 4,
			},
			auth: &user.User{ID: 1},
		},
		{
			name: "delete nonexistant",
			fields: fields{
				TaskID:  1,
				LabelID: 1,
			},
			auth:          &user.User{ID: 1},
			wantForbidden: true,
		},
		{
			name: "nonexisting label",
			fields: fields{
				TaskID:  1,
				LabelID: 9999,
			},
			auth:          &user.User{ID: 1},
			wantForbidden: true,
		},
		{
			name: "nonexisting task",
			fields: fields{
				TaskID:  9999,
				LabelID: 1,
			},
			auth:          &user.User{ID: 1},
			wantForbidden: true,
		},
		{
			name: "existing, but forbidden task",
			fields: fields{
				TaskID:  14,
				LabelID: 1,
			},
			auth:          &user.User{ID: 1},
			wantForbidden: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.LoadAndAssertFixtures(t)

			l := &LabelTask{
				ID:       tt.fields.ID,
				TaskID:   tt.fields.TaskID,
				LabelID:  tt.fields.LabelID,
				Created:  tt.fields.Created,
				CRUDable: tt.fields.CRUDable,
				Rights:   tt.fields.Rights,
			}
			allowed, _ := l.CanDelete(tt.auth)
			if !allowed && !tt.wantForbidden {
				t.Errorf("LabelTask.CanDelete() forbidden, want %v", tt.wantForbidden)
			}
			err := l.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("LabelTask.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("LabelTask.Delete() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
			if !tt.wantForbidden {
				db.AssertMissing(t, "label_task", map[string]interface{}{
					"label_id": l.LabelID,
					"task_id":  l.TaskID,
				})
			}
		})
	}
}
