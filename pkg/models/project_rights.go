// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2021 Vikunja and contributors. All rights reserved.
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

package models

import (
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/web"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// CanWrite return whether the user can write on that project or not
func (l *Project) CanWrite(s *xorm.Session, a web.Auth) (bool, error) {

	// The favorite project can't be edited
	if l.ID == FavoritesPseudoProject.ID {
		return false, nil
	}

	// Get the project and check the right
	originalProject, err := GetProjectSimpleByID(s, l.ID)
	if err != nil {
		return false, err
	}

	// We put the result of the is archived check in a separate variable to be able to return it later without
	// needing to recheck it again
	errIsArchived := originalProject.CheckIsArchived(s)

	var canWrite bool

	// Check if we're dealing with a share auth
	shareAuth, ok := a.(*LinkSharing)
	if ok {
		return originalProject.ID == shareAuth.ProjectID &&
			(shareAuth.Right == RightWrite || shareAuth.Right == RightAdmin), errIsArchived
	}

	// Check if the user is either owner or can write to the project
	if originalProject.isOwner(&user.User{ID: a.GetID()}) {
		canWrite = true
	}

	if canWrite {
		return canWrite, errIsArchived
	}

	canWrite, _, err = originalProject.checkRight(s, a, RightWrite, RightAdmin)
	if err != nil {
		return false, err
	}
	return canWrite, errIsArchived
}

// CanRead checks if a user has read access to a project
func (l *Project) CanRead(s *xorm.Session, a web.Auth) (bool, int, error) {

	// The favorite project needs a special treatment
	if l.ID == FavoritesPseudoProject.ID {
		owner, err := user.GetFromAuth(a)
		if err != nil {
			return false, 0, err
		}

		*l = FavoritesPseudoProject
		l.Owner = owner
		return true, int(RightRead), nil
	}

	// Saved Filter Projects need a special case
	if getSavedFilterIDFromProjectID(l.ID) > 0 {
		sf := &SavedFilter{ID: getSavedFilterIDFromProjectID(l.ID)}
		return sf.CanRead(s, a)
	}

	// Check if the user is either owner or can read
	var err error
	originalProject, err := GetProjectSimpleByID(s, l.ID)
	if err != nil {
		return false, 0, err
	}

	*l = *originalProject

	// Check if we're dealing with a share auth
	shareAuth, ok := a.(*LinkSharing)
	if ok {
		return l.ID == shareAuth.ProjectID &&
			(shareAuth.Right == RightRead || shareAuth.Right == RightWrite || shareAuth.Right == RightAdmin), int(shareAuth.Right), nil
	}

	if l.isOwner(&user.User{ID: a.GetID()}) {
		return true, int(RightAdmin), nil
	}
	return l.checkRight(s, a, RightRead, RightWrite, RightAdmin)
}

// CanUpdate checks if the user can update a project
func (l *Project) CanUpdate(s *xorm.Session, a web.Auth) (canUpdate bool, err error) {
	// The favorite project can't be edited
	if l.ID == FavoritesPseudoProject.ID {
		return false, nil
	}

	// Get the project
	ol, err := GetProjectSimpleByID(s, l.ID)
	if err != nil {
		return false, err
	}

	// Check if we're moving the project into a different namespace.
	// If that is the case, we need to verify permissions to do so.
	if l.NamespaceID != 0 && l.NamespaceID != ol.NamespaceID {
		newNamespace := &Namespace{ID: l.NamespaceID}
		can, err := newNamespace.CanWrite(s, a)
		if err != nil {
			return false, err
		}
		if !can {
			return false, ErrGenericForbidden{}
		}
	}

	fid := getSavedFilterIDFromProjectID(l.ID)
	if fid > 0 {
		sf, err := getSavedFilterSimpleByID(s, fid)
		if err != nil {
			return false, err
		}

		return sf.CanUpdate(s, a)
	}

	canUpdate, err = l.CanWrite(s, a)
	// If the project is archived and the user tries to un-archive it, let the request through
	if IsErrProjectIsArchived(err) && !l.IsArchived {
		err = nil
	}
	return canUpdate, err
}

// CanDelete checks if the user can delete a project
func (l *Project) CanDelete(s *xorm.Session, a web.Auth) (bool, error) {
	return l.IsAdmin(s, a)
}

// CanCreate checks if the user can create a project
func (l *Project) CanCreate(s *xorm.Session, a web.Auth) (bool, error) {
	// A user can create a project if they have write access to the namespace
	n := &Namespace{ID: l.NamespaceID}
	return n.CanWrite(s, a)
}

// IsAdmin returns whether the user has admin rights on the project or not
func (l *Project) IsAdmin(s *xorm.Session, a web.Auth) (bool, error) {
	// The favorite project can't be edited
	if l.ID == FavoritesPseudoProject.ID {
		return false, nil
	}

	originalProject, err := GetProjectSimpleByID(s, l.ID)
	if err != nil {
		return false, err
	}

	// Check if we're dealing with a share auth
	shareAuth, ok := a.(*LinkSharing)
	if ok {
		return originalProject.ID == shareAuth.ProjectID && shareAuth.Right == RightAdmin, nil
	}

	// Check all the things
	// Check if the user is either owner or can write to the project
	// Owners are always admins
	if originalProject.isOwner(&user.User{ID: a.GetID()}) {
		return true, nil
	}
	is, _, err := originalProject.checkRight(s, a, RightAdmin)
	return is, err
}

// Little helper function to check if a user is project owner
func (l *Project) isOwner(u *user.User) bool {
	return l.OwnerID == u.ID
}

// Checks n different rights for any given user
func (l *Project) checkRight(s *xorm.Session, a web.Auth, rights ...Right) (bool, int, error) {

	/*
			The following loop creates a sql condition like this one:

		    (ul.user_id = 1 AND ul.right = 1) OR (un.user_id = 1 AND un.right = 1) OR
			(tm.user_id = 1 AND tn.right = 1) OR (tm2.user_id = 1 AND tl.right = 1) OR

			for each passed right. That way, we can check with a single sql query (instead if 8)
			if the user has the right to see the project or not.
	*/

	var conds []builder.Cond
	for _, r := range rights {
		// User conditions
		// If the project was shared directly with the user and the user has the right
		conds = append(conds, builder.And(
			builder.Eq{"ul.user_id": a.GetID()},
			builder.Eq{"ul.right": r},
		))
		// If the namespace this project belongs to was shared directly with the user and the user has the right
		conds = append(conds, builder.And(
			builder.Eq{"un.user_id": a.GetID()},
			builder.Eq{"un.right": r},
		))

		// Team rights
		// If the project was shared directly with the team and the team has the right
		conds = append(conds, builder.And(
			builder.Eq{"tm2.user_id": a.GetID()},
			builder.Eq{"tl.right": r},
		))
		// If the namespace this project belongs to was shared directly with the team and the team has the right
		conds = append(conds, builder.And(
			builder.Eq{"tm.user_id": a.GetID()},
			builder.Eq{"tn.right": r},
		))
	}

	// If the user is the owner of a namespace, it has any right, all the time
	conds = append(conds, builder.Eq{"n.owner_id": a.GetID()})

	type allProjectRights struct {
		UserNamespace *NamespaceUser `xorm:"extends"`
		UserProject   *ProjectUser   `xorm:"extends"`

		TeamNamespace *TeamNamespace `xorm:"extends"`
		TeamProject   *TeamProject   `xorm:"extends"`

		NamespaceOwnerID int64 `xorm:"namespaces_owner_id"`
	}

	r := &allProjectRights{}
	var maxRight = 0
	exists, err := s.
		Select("l.*, un.right, ul.right, tn.right, tl.right, n.owner_id as namespaces_owner_id").
		Table("projects").
		Alias("l").
		// User stuff
		Join("LEFT", []string{"users_namespaces", "un"}, "un.namespace_id = l.namespace_id").
		Join("LEFT", []string{"users_projects", "ul"}, "ul.project_id = l.id").
		Join("LEFT", []string{"namespaces", "n"}, "n.id = l.namespace_id").
		// Team stuff
		Join("LEFT", []string{"team_namespaces", "tn"}, " l.namespace_id = tn.namespace_id").
		Join("LEFT", []string{"team_members", "tm"}, "tm.team_id = tn.team_id").
		Join("LEFT", []string{"team_projects", "tl"}, "l.id = tl.project_id").
		Join("LEFT", []string{"team_members", "tm2"}, "tm2.team_id = tl.team_id").
		// The actual condition
		Where(builder.And(
			builder.Or(
				conds...,
			),
			builder.Eq{"l.id": l.ID},
		)).
		Get(r)

	// Figure out the max right and return it
	if int(r.UserNamespace.Right) > maxRight {
		maxRight = int(r.UserNamespace.Right)
	}
	if int(r.UserProject.Right) > maxRight {
		maxRight = int(r.UserProject.Right)
	}
	if int(r.TeamNamespace.Right) > maxRight {
		maxRight = int(r.TeamNamespace.Right)
	}
	if int(r.TeamProject.Right) > maxRight {
		maxRight = int(r.TeamProject.Right)
	}
	if r.NamespaceOwnerID == a.GetID() {
		maxRight = int(RightAdmin)
	}

	return exists, maxRight, err
}