// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package routes

import (
	"net/http"

	"code.vikunja.io/api/pkg/models"

	"code.vikunja.io/web"
	"github.com/asaskevich/govalidator"
)

// CustomValidator is a dummy struct to use govalidator with echo
type CustomValidator struct{}

func init() {
	govalidator.TagMap["time"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsTime(str, "15:04")
	})
}

// Validate validates stuff
func (cv *CustomValidator) Validate(i interface{}) error {
	if _, err := govalidator.ValidateStruct(i); err != nil {

		var errs []string
		for field, e := range govalidator.ErrorsByField(err) {
			errs = append(errs, field+": "+e)
		}

		return models.ValidationHTTPError{
			HTTPError: web.HTTPError{
				HTTPCode: http.StatusPreconditionFailed,
				Code:     models.ErrCodeInvalidData,
				Message:  "Invalid Data",
			},
			InvalidFields: errs,
		}
	}
	return nil
}
