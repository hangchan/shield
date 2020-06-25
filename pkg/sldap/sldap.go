package sldap

import (
	"fmt"
	ldap "github.com/go-ldap/ldap"
	util "github.com/hangchan/shield/pkg/util"
)

type LdapConn struct {
	LdapURL 		string
	BaseDN			string
	BindUser		string
	BindPassword	string
	LdapUser		string
}

func Test(lc LdapConn, memberUid string) []string {
	var results []string
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%v))", memberUid)
	//filter := "(&(objectClass=posixGroup)(memberUid=hchan))"

	l, err := ldap.DialURL(lc.LdapURL)
	util.LogError(err)
	err = l.Bind(lc.BindUser, lc.BindPassword)
	util.LogError(err)
	defer l.Close()

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
		results = append(results, fmt.Sprintf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn")))
	}

	return results

}