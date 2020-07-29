package sldap

import (
	"fmt"
	ldap "github.com/go-ldap/ldap"
	util "github.com/hangchan/shield/pkg/util"
	"sort"
)

type LdapConn struct {
	LdapURL 		string
	BaseDN			string
	BindUser		string
	BindPassword	string
	LdapUser		string
}

func ldapConn(lc LdapConn) *ldap.Conn {
	l, err := ldap.DialURL(lc.LdapURL)
	util.LogError(err)
	err = l.Bind(lc.BindUser, lc.BindPassword)
	util.LogError(err)

	return l
}

func Search(lc LdapConn, memberUid string) []string {
	l := ldapConn(lc)
	defer l.Close()

	var resultArr []string
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%v))", memberUid)

	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, // The filter to apply
		[]string{"dn", "cn"},                    // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	util.LogError(err)

	for _, entry := range sr.Entries {
		//results = append(results, fmt.Sprintf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn")))
		resultArr = append(resultArr, fmt.Sprintf("%v", entry.GetAttributeValue("cn")))
	}

	sort.Strings(resultArr)

	return resultArr

}

func RemoveFromGroup(lc LdapConn, memberUid string, group string) {
	l := ldapConn(lc)
	defer l.Close()

	groupOU := "Groups"

	mr := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=%s,%s", group, groupOU, lc.BaseDN), []ldap.Control{})
	mr.Delete("memberUid", []string{fmt.Sprintf("%s", memberUid)})

	err := l.Modify(mr)
	util.LogError(err)

}
