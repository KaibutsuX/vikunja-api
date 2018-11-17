package models

// Namespace holds informations about a namespace
type Namespace struct {
	ID          int64  `xorm:"int(11) autoincr not null unique pk" json:"id" param:"namespace"`
	Name        string `xorm:"varchar(250)" json:"name" valid:"required,runelength(5|250)"`
	Description string `xorm:"varchar(1000)" json:"description" valid:"runelength(0|250)"`
	OwnerID     int64  `xorm:"int(11) not null INDEX" json:"-"`

	Owner User `xorm:"-" json:"owner"`

	Created int64 `xorm:"created" json:"created" valid:"range(0|0)"`
	Updated int64 `xorm:"updated" json:"updated" valid:"range(0|0)"`

	CRUDable `xorm:"-" json:"-"`
	Rights   `xorm:"-" json:"-"`
}

// TableName makes beautiful table names
func (Namespace) TableName() string {
	return "namespaces"
}

// GetNamespaceByID returns a namespace object by its ID
func GetNamespaceByID(id int64) (namespace Namespace, err error) {
	if id < 1 {
		return namespace, ErrNamespaceDoesNotExist{ID: id}
	}

	namespace.ID = id
	exists, err := x.Get(&namespace)
	if err != nil {
		return namespace, err
	}

	if !exists {
		return namespace, ErrNamespaceDoesNotExist{ID: id}
	}

	// Get the namespace Owner
	namespace.Owner, err = GetUserByID(namespace.OwnerID)
	if err != nil {
		return namespace, err
	}

	return namespace, err
}

// ReadOne gets one namespace
// @Summary Gets one namespace
// @Description Returns a namespace by its ID.
// @tags namespace
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Namespace ID"
// @Success 200 {object} models.Namespace "The Namespace"
// @Failure 403 {object} models.HTTPError "The user does not have access to that namespace."
// @Failure 500 {object} models.Message "Internal error"
// @Router /namespaces/{id} [get]
func (n *Namespace) ReadOne() (err error) {
	*n, err = GetNamespaceByID(n.ID)
	return
}

// NamespaceWithLists represents a namespace with list meta informations
type NamespaceWithLists struct {
	Namespace `xorm:"extends"`
	Lists     []*List `xorm:"-" json:"lists"`
}

// ReadAll gets all namespaces a user has access to
// @Summary Get all namespaces a user has access to
// @Description Returns all namespaces a user has access to.
// @tags namespace
// @Accept json
// @Produce json
// @Param p query int false "The page number. Used for pagination. If not provided, the first page of results is returned."
// @Param s query string false "Search namespaces by name."
// @Security ApiKeyAuth
// @Success 200 {array} models.NamespaceWithLists "The Namespaces."
// @Failure 500 {object} models.Message "Internal error"
// @Router /namespaces [get]
func (n *Namespace) ReadAll(search string, doer *User, page int) (interface{}, error) {

	all := []*NamespaceWithLists{}

	err := x.Select("namespaces.*").
		Table("namespaces").
		Join("LEFT", "team_namespaces", "namespaces.id = team_namespaces.namespace_id").
		Join("LEFT", "team_members", "team_members.team_id = team_namespaces.team_id").
		Join("LEFT", "users_namespace", "users_namespace.namespace_id = namespaces.id").
		Where("team_members.user_id = ?", doer.ID).
		Or("namespaces.owner_id = ?", doer.ID).
		Or("users_namespace.user_id = ?", doer.ID).
		GroupBy("namespaces.id").
		Limit(getLimitFromPageIndex(page)).
		Where("namespaces.name LIKE ?", "%"+search+"%").
		Find(&all)

	if err != nil {
		return all, err
	}

	// Get all users
	users := []*User{}
	err = x.Select("users.*").
		Table("namespaces").
		Join("LEFT", "team_namespaces", "namespaces.id = team_namespaces.namespace_id").
		Join("LEFT", "team_members", "team_members.team_id = team_namespaces.team_id").
		Join("INNER", "users", "users.id = namespaces.owner_id").
		Where("team_members.user_id = ?", doer.ID).
		Or("namespaces.owner_id = ?", doer.ID).
		GroupBy("users.id").
		Find(&users)

	if err != nil {
		return all, err
	}

	// Make a list of namespace ids
	var namespaceids []int64
	for _, nsp := range all {
		namespaceids = append(namespaceids, nsp.ID)
	}

	// Get all lists
	lists := []*List{}
	err = x.Table(&lists).
		In("namespace_id", namespaceids).
		Find(&lists)
	if err != nil {
		return all, err
	}

	// More details for the lists
	AddListDetails(lists)

	// Put objects in our namespace list
	for i, n := range all {

		// Users
		for _, u := range users {
			if n.OwnerID == u.ID {
				all[i].Owner = *u
				break
			}
		}

		// List infos
		for _, l := range lists {
			if n.ID == l.NamespaceID {
				all[i].Lists = append(all[i].Lists, l)
			}
		}
	}

	return all, nil
}