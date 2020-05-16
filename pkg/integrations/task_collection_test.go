// Copyright 2018-2020 Vikunja and contriubtors. All rights reserved.
//
// This file is part of Vikunja.
//
// Vikunja is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Vikunja is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Vikunja.  If not, see <https://www.gnu.org/licenses/>.

package integrations

import (
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/web/handler"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestTaskCollection(t *testing.T) {
	testHandler := webHandlerTest{
		user: &testuser1,
		strFunc: func() handler.CObject {
			return &models.TaskCollection{}
		},
		t: t,
	}
	t.Run("ReadAll on list", func(t *testing.T) {

		urlParams := map[string]string{"list": "1"}

		t.Run("Normal", func(t *testing.T) {
			rec, err := testHandler.testReadAllWithUser(nil, urlParams)
			assert.NoError(t, err)
			// Not using assert.Equal to avoid having the tests break every time we add new fixtures
			assert.Contains(t, rec.Body.String(), `task #1`)
			assert.Contains(t, rec.Body.String(), `task #2`)
			assert.Contains(t, rec.Body.String(), `task #3`)
			assert.Contains(t, rec.Body.String(), `task #4`)
			assert.Contains(t, rec.Body.String(), `task #5`)
			assert.Contains(t, rec.Body.String(), `task #6`)
			assert.Contains(t, rec.Body.String(), `task #7`)
			assert.Contains(t, rec.Body.String(), `task #8`)
			assert.Contains(t, rec.Body.String(), `task #9`)
			assert.Contains(t, rec.Body.String(), `task #10`)
			assert.Contains(t, rec.Body.String(), `task #11`)
			assert.Contains(t, rec.Body.String(), `task #12`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
			assert.NotContains(t, rec.Body.String(), `task #15`) // Shared via team readonly
			assert.NotContains(t, rec.Body.String(), `task #16`) // Shared via team write
			assert.NotContains(t, rec.Body.String(), `task #17`) // Shared via team admin
			assert.NotContains(t, rec.Body.String(), `task #18`) // Shared via user readonly
			assert.NotContains(t, rec.Body.String(), `task #19`) // Shared via user write
			assert.NotContains(t, rec.Body.String(), `task #20`) // Shared via user admin
			assert.NotContains(t, rec.Body.String(), `task #21`) // Shared via namespace team readonly
			assert.NotContains(t, rec.Body.String(), `task #22`) // Shared via namespace team write
			assert.NotContains(t, rec.Body.String(), `task #23`) // Shared via namespace team admin
			assert.NotContains(t, rec.Body.String(), `task #24`) // Shared via namespace user readonly
			assert.NotContains(t, rec.Body.String(), `task #25`) // Shared via namespace user write
			assert.NotContains(t, rec.Body.String(), `task #26`) // Shared via namespace user admin
			assert.Contains(t, rec.Body.String(), `task #27`)
			assert.Contains(t, rec.Body.String(), `task #28`)
			assert.NotContains(t, rec.Body.String(), `task #32`)
		})
		t.Run("Search", func(t *testing.T) {
			rec, err := testHandler.testReadAllWithUser(url.Values{"s": []string{"task #6"}}, urlParams)
			assert.NoError(t, err)
			assert.NotContains(t, rec.Body.String(), `task #1`)
			assert.NotContains(t, rec.Body.String(), `task #2`)
			assert.NotContains(t, rec.Body.String(), `task #3`)
			assert.NotContains(t, rec.Body.String(), `task #4`)
			assert.NotContains(t, rec.Body.String(), `task #5`)
			assert.Contains(t, rec.Body.String(), `task #6`)
			assert.NotContains(t, rec.Body.String(), `task #7`)
			assert.NotContains(t, rec.Body.String(), `task #8`)
			assert.NotContains(t, rec.Body.String(), `task #9`)
			assert.NotContains(t, rec.Body.String(), `task #10`)
			assert.NotContains(t, rec.Body.String(), `task #11`)
			assert.NotContains(t, rec.Body.String(), `task #12`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
		})
		t.Run("Search case insensitive", func(t *testing.T) {
			rec, err := testHandler.testReadAllWithUser(url.Values{"s": []string{"tASk #6"}}, urlParams)
			assert.NoError(t, err)
			assert.NotContains(t, rec.Body.String(), `task #1`)
			assert.NotContains(t, rec.Body.String(), `task #2`)
			assert.NotContains(t, rec.Body.String(), `task #3`)
			assert.NotContains(t, rec.Body.String(), `task #4`)
			assert.NotContains(t, rec.Body.String(), `task #5`)
			assert.Contains(t, rec.Body.String(), `task #6`)
			assert.NotContains(t, rec.Body.String(), `task #7`)
			assert.NotContains(t, rec.Body.String(), `task #8`)
			assert.NotContains(t, rec.Body.String(), `task #9`)
			assert.NotContains(t, rec.Body.String(), `task #10`)
			assert.NotContains(t, rec.Body.String(), `task #11`)
			assert.NotContains(t, rec.Body.String(), `task #12`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
		})
		t.Run("Sort Order", func(t *testing.T) {
			// TODO: Add more cases
			// should equal priority asc
			t.Run("by priority", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":33,"title":"task #33 with percent done","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0.5,"identifier":"test1-17","index":17,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":1,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-4","index":4,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by priority desc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}, "order_by": []string{"desc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `[{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1`)
			})
			t.Run("by priority asc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}, "order_by": []string{"asc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":33,"title":"task #33 with percent done","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0.5,"identifier":"test1-17","index":17,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":1,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-4","index":4,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			// should equal duedate asc
			t.Run("by due_date", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by duedate desc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}, "order_by": []string{"desc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `[{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":6,"title":"task #6 lower due date`)
			})
			// Due date without unix suffix
			t.Run("by duedate asc without _unix suffix", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date"}, "order_by": []string{"asc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by due_date without _unix suffix", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by duedate desc without _unix suffix", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date"}, "order_by": []string{"desc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `[{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":6,"title":"task #6 lower due date`)
			})
			t.Run("by duedate asc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}, "order_by": []string{"asc"}}, urlParams)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("invalid sort parameter", func(t *testing.T) {
				_, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"loremipsum"}}, urlParams)
				assert.Error(t, err)
				assertHandlerErrorCode(t, err, models.ErrCodeInvalidTaskField)
			})
			t.Run("invalid sort order", func(t *testing.T) {
				_, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"id"}, "order_by": []string{"loremipsum"}}, urlParams)
				assert.Error(t, err)
				assertHandlerErrorCode(t, err, models.ErrCodeInvalidSortOrder)
			})
			t.Run("invalid parameter", func(t *testing.T) {
				// Invalid parameter should not sort at all
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort": []string{"loremipsum"}}, urlParams)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `[{"id":3,"title":"task #3 high prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":100,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":4,"title":"task #4 low prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":1`)
				assert.NotContains(t, rec.Body.String(), `{"id":4,"title":"task #4 low prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":1,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":3,"title":"task #3 high prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":100,"start_date":0,"end_date":0,"assignees":null,"labels":null,"created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}}]`)
				assert.NotContains(t, rec.Body.String(), `[{"id":5,"title":"task #5 higher due date","description":"","done":false,"due_date":1543636724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":6,"title":"task #6 lower due date"`)
				assert.NotContains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"due_date":1543616724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"due_date":1543636724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}}]`)
			})
		})
		t.Run("Date range", func(t *testing.T) {
			t.Run("start and end date", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"start_date", "end_date", "due_date"},
						"filter_value":      []string{"1544500000", "1544700001", "1543500000"},
						"filter_comparator": []string{"greater", "less", "greater"},
					},
					urlParams,
				)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `task #1`)
				assert.NotContains(t, rec.Body.String(), `task #2`)
				assert.NotContains(t, rec.Body.String(), `task #3`)
				assert.NotContains(t, rec.Body.String(), `task #4`)
				assert.Contains(t, rec.Body.String(), `task #5`)
				assert.Contains(t, rec.Body.String(), `task #6`)
				assert.Contains(t, rec.Body.String(), `task #7`)
				assert.Contains(t, rec.Body.String(), `task #8`)
				assert.Contains(t, rec.Body.String(), `task #9`)
				assert.NotContains(t, rec.Body.String(), `task #10`)
				assert.NotContains(t, rec.Body.String(), `task #11`)
				assert.NotContains(t, rec.Body.String(), `task #12`)
				assert.NotContains(t, rec.Body.String(), `task #13`)
				assert.NotContains(t, rec.Body.String(), `task #14`)
			})
			t.Run("start date only", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"start_date"},
						"filter_value":      []string{"1540000000"},
						"filter_comparator": []string{"greater"},
					},
					urlParams,
				)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `task #1`)
				assert.NotContains(t, rec.Body.String(), `task #2`)
				assert.NotContains(t, rec.Body.String(), `task #3`)
				assert.NotContains(t, rec.Body.String(), `task #4`)
				assert.NotContains(t, rec.Body.String(), `task #5`)
				assert.NotContains(t, rec.Body.String(), `task #6`)
				assert.Contains(t, rec.Body.String(), `task #7`)
				assert.NotContains(t, rec.Body.String(), `task #8`)
				assert.Contains(t, rec.Body.String(), `task #9`)
				assert.NotContains(t, rec.Body.String(), `task #10`)
				assert.NotContains(t, rec.Body.String(), `task #11`)
				assert.NotContains(t, rec.Body.String(), `task #12`)
				assert.NotContains(t, rec.Body.String(), `task #13`)
				assert.NotContains(t, rec.Body.String(), `task #14`)
			})
			t.Run("end date only", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"end_date"},
						"filter_value":      []string{"1544700001"},
						"filter_comparator": []string{"greater"},
					},
					urlParams,
				)
				assert.NoError(t, err)
				// If no start date but an end date is specified, this should be null
				// since we don't have any tasks in the fixtures with an end date >
				// the current date.
				assert.Equal(t, "[]\n", rec.Body.String())
			})
		})
	})

	t.Run("ReadAll for all tasks", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			rec, err := testHandler.testReadAllWithUser(nil, nil)
			assert.NoError(t, err)
			// Not using assert.Equal to avoid having the tests break every time we add new fixtures
			assert.Contains(t, rec.Body.String(), `task #1`)
			assert.Contains(t, rec.Body.String(), `task #2`)
			assert.Contains(t, rec.Body.String(), `task #3`)
			assert.Contains(t, rec.Body.String(), `task #4`)
			assert.Contains(t, rec.Body.String(), `task #5`)
			assert.Contains(t, rec.Body.String(), `task #6`)
			assert.Contains(t, rec.Body.String(), `task #7`)
			assert.Contains(t, rec.Body.String(), `task #8`)
			assert.Contains(t, rec.Body.String(), `task #9`)
			assert.Contains(t, rec.Body.String(), `task #10`)
			assert.Contains(t, rec.Body.String(), `task #11`)
			assert.Contains(t, rec.Body.String(), `task #12`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
			assert.Contains(t, rec.Body.String(), `task #15`) // Shared via team readonly
			assert.Contains(t, rec.Body.String(), `task #16`) // Shared via team write
			assert.Contains(t, rec.Body.String(), `task #17`) // Shared via team admin
			assert.Contains(t, rec.Body.String(), `task #18`) // Shared via user readonly
			assert.Contains(t, rec.Body.String(), `task #19`) // Shared via user write
			assert.Contains(t, rec.Body.String(), `task #20`) // Shared via user admin
			assert.Contains(t, rec.Body.String(), `task #21`) // Shared via namespace team readonly
			assert.Contains(t, rec.Body.String(), `task #22`) // Shared via namespace team write
			assert.Contains(t, rec.Body.String(), `task #23`) // Shared via namespace team admin
			assert.Contains(t, rec.Body.String(), `task #24`) // Shared via namespace user readonly
			assert.Contains(t, rec.Body.String(), `task #25`) // Shared via namespace user write
			assert.Contains(t, rec.Body.String(), `task #26`) // Shared via namespace user admin
			// TODO: Add some cases where the user has access to the list, somhow shared
		})
		t.Run("Search", func(t *testing.T) {
			rec, err := testHandler.testReadAllWithUser(url.Values{"s": []string{"task #6"}}, nil)
			assert.NoError(t, err)
			assert.NotContains(t, rec.Body.String(), `task #1`)
			assert.NotContains(t, rec.Body.String(), `task #2`)
			assert.NotContains(t, rec.Body.String(), `task #3`)
			assert.NotContains(t, rec.Body.String(), `task #4`)
			assert.NotContains(t, rec.Body.String(), `task #5`)
			assert.Contains(t, rec.Body.String(), `task #6`)
			assert.NotContains(t, rec.Body.String(), `task #7`)
			assert.NotContains(t, rec.Body.String(), `task #8`)
			assert.NotContains(t, rec.Body.String(), `task #9`)
			assert.NotContains(t, rec.Body.String(), `task #10`)
			assert.NotContains(t, rec.Body.String(), `task #11`)
			assert.NotContains(t, rec.Body.String(), `task #12`)
			assert.NotContains(t, rec.Body.String(), `task #13`)
			assert.NotContains(t, rec.Body.String(), `task #14`)
		})
		t.Run("Sort Order", func(t *testing.T) {
			// should equal priority asc
			t.Run("by priority", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":33,"title":"task #33 with percent done","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0.5,"identifier":"test1-17","index":17,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":1,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-4","index":4,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by priority desc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}, "order_by": []string{"desc"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `[{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1`)
			})
			t.Run("by priority asc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"priority"}, "order_by": []string{"asc"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":33,"title":"task #33 with percent done","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0.5,"identifier":"test1-17","index":17,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":1,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":4,"title":"task #4 low prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":1,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-4","index":4,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":3,"title":"task #3 high prio","description":"","done":false,"done_at":null,"due_date":null,"reminder_dates":null,"list_id":1,"repeat_after":0,"priority":100,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-3","index":3,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			// should equal duedate asc
			t.Run("by due_date", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("by duedate desc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}, "order_by": []string{"desc"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `[{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":6,"title":"task #6 lower due date`)
			})
			t.Run("by duedate asc", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort_by": []string{"due_date_unix"}, "order_by": []string{"asc"}}, nil)
				assert.NoError(t, err)
				assert.Contains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"done_at":null,"due_date":"2018-11-30T22:25:24Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-6","index":6,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":3,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"done_at":null,"due_date":"2018-12-01T03:58:44Z","reminder_dates":null,"list_id":1,"repeat_after":0,"priority":0,"start_date":null,"end_date":null,"assignees":null,"labels":null,"hex_color":"","percent_done":0,"identifier":"test1-5","index":5,"related_tasks":{},"attachments":null,"created":"2018-12-01T01:12:04Z","updated":"2018-12-01T01:12:04Z","bucket_id":2,"position":0,"created_by":{"id":1,"username":"user1","created":null,"updated":null}}]`)
			})
			t.Run("invalid parameter", func(t *testing.T) {
				// Invalid parameter should not sort at all
				rec, err := testHandler.testReadAllWithUser(url.Values{"sort": []string{"loremipsum"}}, nil)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `[{"id":3,"title":"task #3 high prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":100,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":4,"title":"task #4 low prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":1`)
				assert.NotContains(t, rec.Body.String(), `{"id":4,"title":"task #4 low prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":1,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":3,"title":"task #3 high prio","description":"","done":false,"due_date":0,"reminder_dates":null,"repeat_after":0,"priority":100,"start_date":0,"end_date":0,"assignees":null,"labels":null,"created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}}]`)
				assert.NotContains(t, rec.Body.String(), `[{"id":5,"title":"task #5 higher due date","description":"","done":false,"due_date":1543636724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":6,"title":"task #6 lower due date"`)
				assert.NotContains(t, rec.Body.String(), `{"id":6,"title":"task #6 lower due date","description":"","done":false,"due_date":1543616724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"hex_color":"","created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}},{"id":5,"title":"task #5 higher due date","description":"","done":false,"due_date":1543636724,"reminder_dates":null,"repeat_after":0,"priority":0,"start_date":0,"end_date":0,"assignees":null,"labels":null,"created":1543626724,"updated":1543626724,"created_by":{"id":0,"username":"","email":"","created":0,"updated":0}}]`)
			})
		})
		t.Run("Date range", func(t *testing.T) {
			t.Run("start and end date", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"start_date", "end_date", "due_date"},
						"filter_value":      []string{"1544500000", "1544700001", "1543500000"},
						"filter_comparator": []string{"greater", "less", "greater"},
					},
					nil,
				)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `task #1`)
				assert.NotContains(t, rec.Body.String(), `task #2`)
				assert.NotContains(t, rec.Body.String(), `task #3`)
				assert.NotContains(t, rec.Body.String(), `task #4`)
				assert.Contains(t, rec.Body.String(), `task #5`)
				assert.Contains(t, rec.Body.String(), `task #6`)
				assert.Contains(t, rec.Body.String(), `task #7`)
				assert.Contains(t, rec.Body.String(), `task #8`)
				assert.Contains(t, rec.Body.String(), `task #9`)
				assert.NotContains(t, rec.Body.String(), `task #10`)
				assert.NotContains(t, rec.Body.String(), `task #11`)
				assert.NotContains(t, rec.Body.String(), `task #12`)
				assert.NotContains(t, rec.Body.String(), `task #13`)
				assert.NotContains(t, rec.Body.String(), `task #14`)
			})
			t.Run("start date only", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"start_date"},
						"filter_value":      []string{"1540000000"},
						"filter_comparator": []string{"greater"},
					},
					nil,
				)
				assert.NoError(t, err)
				assert.NotContains(t, rec.Body.String(), `task #1`)
				assert.NotContains(t, rec.Body.String(), `task #2`)
				assert.NotContains(t, rec.Body.String(), `task #3`)
				assert.NotContains(t, rec.Body.String(), `task #4`)
				assert.NotContains(t, rec.Body.String(), `task #5`)
				assert.NotContains(t, rec.Body.String(), `task #6`)
				assert.Contains(t, rec.Body.String(), `task #7`)
				assert.NotContains(t, rec.Body.String(), `task #8`)
				assert.Contains(t, rec.Body.String(), `task #9`)
				assert.NotContains(t, rec.Body.String(), `task #10`)
				assert.NotContains(t, rec.Body.String(), `task #11`)
				assert.NotContains(t, rec.Body.String(), `task #12`)
				assert.NotContains(t, rec.Body.String(), `task #13`)
				assert.NotContains(t, rec.Body.String(), `task #14`)
			})
			t.Run("end date only", func(t *testing.T) {
				rec, err := testHandler.testReadAllWithUser(
					url.Values{
						"filter_by":         []string{"end_date"},
						"filter_value":      []string{"1544700001"},
						"filter_comparator": []string{"greater"},
					},
					nil,
				)
				assert.NoError(t, err)
				// If no start date but an end date is specified, this should be null
				// since we don't have any tasks in the fixtures with an end date >
				// the current date.
				assert.Equal(t, "[]\n", rec.Body.String())
			})
		})
	})

}
