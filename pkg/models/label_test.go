//  Vikunja is a todo-list application to facilitate your life.
//  Copyright 2018 Vikunja and contributors. All rights reserved.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"reflect"
	"runtime"
	"testing"

	"code.vikunja.io/web"
)

func TestLabel_ReadAll(t *testing.T) {
	type fields struct {
		ID          int64
		Title       string
		Description string
		HexColor    string
		CreatedByID int64
		CreatedBy   *User
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	type args struct {
		search string
		a      web.Auth
		page   int
	}
	user1 := &User{
		ID:       1,
		Username: "user1",
		Password: "1234",
		Email:    "user1@example.com",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLs  interface{}
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				a: &User{ID: 1},
			},
			wantLs: []*labelWithTaskID{
				{
					Label: Label{
						ID:          1,
						Title:       "Label #1",
						CreatedByID: 1,
						CreatedBy:   user1,
					},
				},
				{
					Label: Label{
						ID:          2,
						Title:       "Label #2",
						CreatedByID: 1,
						CreatedBy:   user1,
					},
				},
				{
					TaskID: 1,
					Label: Label{
						ID:          4,
						Title:       "Label #4 - visible via other task",
						CreatedByID: 2,
						CreatedBy: &User{
							ID:       2,
							Username: "user2",
							Password: "1234",
							Email:    "user2@example.com",
						},
					},
				},
			},
		},
		{
			name:    "invalid user",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Label{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				HexColor:    tt.fields.HexColor,
				CreatedByID: tt.fields.CreatedByID,
				CreatedBy:   tt.fields.CreatedBy,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			gotLs, err := l.ReadAll(tt.args.search, tt.args.a, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("Label.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLs, tt.wantLs) {
				t.Errorf("Label.ReadAll() = %v, want %v", gotLs, tt.wantLs)
			}
		})
	}
}

func TestLabel_ReadOne(t *testing.T) {
	type fields struct {
		ID          int64
		Title       string
		Description string
		HexColor    string
		CreatedByID int64
		CreatedBy   *User
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	user1 := &User{
		ID:       1,
		Username: "user1",
		Password: "1234",
		Email:    "user1@example.com",
	}
	tests := []struct {
		name          string
		fields        fields
		want          *Label
		wantErr       bool
		errType       func(error) bool
		auth          web.Auth
		wantForbidden bool
	}{
		{
			name: "Get label #1",
			fields: fields{
				ID: 1,
			},
			want: &Label{
				ID:          1,
				Title:       "Label #1",
				CreatedByID: 1,
				CreatedBy:   user1,
			},
			auth: &User{ID: 1},
		},
		{
			name: "Get nonexistant label",
			fields: fields{
				ID: 9999,
			},
			wantErr:       true,
			errType:       IsErrLabelDoesNotExist,
			wantForbidden: true,
			auth:          &User{ID: 1},
		},
		{
			name: "no rights",
			fields: fields{
				ID: 3,
			},
			wantForbidden: true,
			auth:          &User{ID: 1},
		},
		{
			name: "Get label #4 - other user",
			fields: fields{
				ID: 4,
			},
			want: &Label{
				ID:          4,
				Title:       "Label #4 - visible via other task",
				CreatedByID: 2,
				CreatedBy: &User{
					ID:       2,
					Username: "user2",
					Password: "1234",
					Email:    "user2@example.com",
				},
			},
			auth: &User{ID: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Label{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				HexColor:    tt.fields.HexColor,
				CreatedByID: tt.fields.CreatedByID,
				CreatedBy:   tt.fields.CreatedBy,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}

			if !l.CanRead(tt.auth) && !tt.wantForbidden {
				t.Errorf("Label.CanRead() forbidden, want %v", tt.wantForbidden)
			}
			err := l.ReadOne()
			if (err != nil) != tt.wantErr {
				t.Errorf("Label.ReadOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && tt.wantErr && !tt.errType(err) {
				t.Errorf("Label.ReadOne() Wrong error type! Error = %v, want = %v", err, runtime.FuncForPC(reflect.ValueOf(tt.errType).Pointer()).Name())
			}
			if !reflect.DeepEqual(l, tt.want) && !tt.wantErr && !tt.wantForbidden {
				t.Errorf("Label.ReadOne() = %v, want %v", l, tt.want)
			}
		})
	}
}

func TestLabel_Create(t *testing.T) {
	type fields struct {
		ID          int64
		Title       string
		Description string
		HexColor    string
		CreatedByID int64
		CreatedBy   *User
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	type args struct {
		a web.Auth
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantForbidden bool
	}{
		{
			name: "normal",
			fields: fields{
				Title:       "Test #1",
				Description: "Lorem Ipsum",
				HexColor:    "ffccff",
			},
			args: args{
				a: &User{ID: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Label{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				HexColor:    tt.fields.HexColor,
				CreatedByID: tt.fields.CreatedByID,
				CreatedBy:   tt.fields.CreatedBy,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			if !l.CanCreate(tt.args.a) && !tt.wantForbidden {
				t.Errorf("Label.CanCreate() forbidden, want %v", tt.wantForbidden)
			}
			if err := l.Create(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Label.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLabel_Update(t *testing.T) {
	type fields struct {
		ID          int64
		Title       string
		Description string
		HexColor    string
		CreatedByID int64
		CreatedBy   *User
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	tests := []struct {
		name          string
		fields        fields
		wantErr       bool
		auth          web.Auth
		wantForbidden bool
	}{
		{
			name: "normal",
			fields: fields{
				ID:    1,
				Title: "new and better",
			},
			auth: &User{ID: 1},
		},
		{
			name: "nonexisting",
			fields: fields{
				ID:    99999,
				Title: "new and better",
			},
			auth:          &User{ID: 1},
			wantForbidden: true,
			wantErr:       true,
		},
		{
			name: "no rights",
			fields: fields{
				ID:    3,
				Title: "new and better",
			},
			auth:          &User{ID: 1},
			wantForbidden: true,
		},
		{
			name: "no rights other creator but access",
			fields: fields{
				ID:    4,
				Title: "new and better",
			},
			auth:          &User{ID: 1},
			wantForbidden: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Label{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				HexColor:    tt.fields.HexColor,
				CreatedByID: tt.fields.CreatedByID,
				CreatedBy:   tt.fields.CreatedBy,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			if !l.CanUpdate(tt.auth) && !tt.wantForbidden {
				t.Errorf("Label.CanUpdate() forbidden, want %v", tt.wantForbidden)
			}
			if err := l.Update(); (err != nil) != tt.wantErr {
				t.Errorf("Label.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLabel_Delete(t *testing.T) {
	type fields struct {
		ID          int64
		Title       string
		Description string
		HexColor    string
		CreatedByID int64
		CreatedBy   *User
		Created     int64
		Updated     int64
		CRUDable    web.CRUDable
		Rights      web.Rights
	}
	tests := []struct {
		name          string
		fields        fields
		wantErr       bool
		auth          web.Auth
		wantForbidden bool
	}{

		{
			name: "normal",
			fields: fields{
				ID: 1,
			},
			auth: &User{ID: 1},
		},
		{
			name: "nonexisting",
			fields: fields{
				ID: 99999,
			},
			auth:          &User{ID: 1},
			wantForbidden: true, // When the label does not exist, it is forbidden. We should fix this, but for everything.
		},
		{
			name: "no rights",
			fields: fields{
				ID: 3,
			},
			auth:          &User{ID: 1},
			wantForbidden: true,
		},
		{
			name: "no rights but visible",
			fields: fields{
				ID: 4,
			},
			auth:          &User{ID: 1},
			wantForbidden: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Label{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				HexColor:    tt.fields.HexColor,
				CreatedByID: tt.fields.CreatedByID,
				CreatedBy:   tt.fields.CreatedBy,
				Created:     tt.fields.Created,
				Updated:     tt.fields.Updated,
				CRUDable:    tt.fields.CRUDable,
				Rights:      tt.fields.Rights,
			}
			if !l.CanDelete(tt.auth) && !tt.wantForbidden {
				t.Errorf("Label.CanDelete() forbidden, want %v", tt.wantForbidden)
			}
			if err := l.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("Label.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}