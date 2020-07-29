package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	sldap "github.com/hangchan/shield/pkg/sldap"
	smysql "github.com/hangchan/shield/pkg/smysql"
	"strings"
)

var lc = sldap.LdapConn{
	LdapURL:		"ldap://<LdapServer>",
	BaseDN:			"dc=example,dc=com",
	BindUser:		"uid=Admin,ou=People,dc=example,dc=com",
	BindPassword:	"<BindPassword>",
	LdapUser:		"",
}

var mc = smysql.MysqlConn{
	DbDriver:		"mysql",
	DbUser:			"<DbUser>",
	DbPass:			"<DbPass>",
	DbName:			"mysql",
	DbAddress:		"<MysqlServer>",
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// RESTy routes for "user" resource
	r.Route("/user", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Nothing here"))
		})

		r.Route("/get", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Nothing here"))
			})

			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")
					ldapGroups := sldap.Search(lc, username)
					ldapGroupsStr := strings.Join(ldapGroups, ",")
					if len(ldapGroups) == 0 {
						w.Write([]byte(fmt.Sprintf("No entries for %s\n\n", username)))
					} else {
						w.Write([]byte("These are the ldap groups the user belongs to:\n\n"))
						w.Write([]byte(fmt.Sprintf("%s: %s\n\n", username, ldapGroupsStr )))
					}

					mysqlHosts := smysql.Search(mc, username)
					mysqlHostsStr := strings.Join(mysqlHosts, ",")
					w.Write([]byte("These are the mysql host entries for the user:\n\n"))
					w.Write([]byte(fmt.Sprintf("username: %s\n", mysqlHostsStr)))
				})
			})
		})

		r.Route("/delete", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Nothing here"))
			})

			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")

					ldapGroups := sldap.Search(lc, username)
					if len(ldapGroups) == 0 {
						w.Write([]byte(fmt.Sprintf("No ldap group entries for %s\n\n", username)))
					} else {
						for _, group := range ldapGroups {
							if group != username {
								sldap.RemoveFromGroup(lc, username, group)
								w.Write([]byte(fmt.Sprintf("Removing user %s from group %s\n", username, group)))
							}
						}
						w.Write([]byte("\n"))
					}

					ldapGroups = sldap.Search(lc, username)
					ldapGroupsStr := strings.Join(ldapGroups, ",")
					if len(ldapGroups) == 0 {
						w.Write([]byte(fmt.Sprintf("No ldap group entries for %s\n\n", username)))
					} else {
						w.Write([]byte("These are the ldap groups the user belongs to:\n\n"))
						w.Write([]byte(fmt.Sprintf("%s: %s\n\n", username, ldapGroupsStr )))
					}

					mysqlHosts := smysql.Search(mc, username)
					if len(mysqlHosts) == 0 {
						w.Write([]byte("No mysql host entries for " + username + "\n"))
					} else {
						for _, host := range mysqlHosts {
							smysql.DropUser(mc, username, host)
							w.Write([]byte("Deleting host from mysql for user:\n"))
							w.Write([]byte(username + "@" + host + "\n\n"))
						}

						mysqlHosts = smysql.Search(mc, username)
						mysqlHostsStr := strings.Join(mysqlHosts, ",")
						w.Write([]byte("These are the mysql host entries for the user:\n\n"))
						w.Write([]byte(fmt.Sprintf("username: %s\n", mysqlHostsStr)))
					}
				})
			})
		})
	})

	http.ListenAndServe(":3333", r)
}
