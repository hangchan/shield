package main

import (
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
	BindPassword:	        "<BindPassword>",
	LdapUser:		"",
}

var mc = smysql.MysqlConn{
	DbDriver:		"mysql",
	DbUser:			"<DbUser>",
	DbPass:			"<DbPass>",
	DbName:			"mysql",
	DbAddress:		"<MysqlServer>",
}

type MysqlConn struct {
	DbDriver 	string
	DbUser		string
	DbPass		string
	DbName		string
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

		r.Route("/getUserGroup", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Nothing here"))
			})

			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")
					ldapGroups := sldap.Search(lc, username)
					w.Write([]byte("These are the ldap groups the user belongs to:\n\n"))
					w.Write([]byte(username + ": " + ldapGroups + "\n\n"))

					mysqlHosts := smysql.Search(mc, username)
					w.Write([]byte("These are the mysql host entries for the user:\n\n"))
					w.Write([]byte(username + ": " + mysqlHosts + "\n"))
				})
			})
		})

		r.Route("/deleteUser", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Nothing here"))
			})

			r.Route("/{username}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					username := chi.URLParam(r, "username")

					mysqlHosts := smysql.Search(mc, username)
					if len(mysqlHosts) == 0 {
						w.Write([]byte("No entries for " + username + "\n"))
					} else {
						mysqlHost := strings.Split(mysqlHosts, ",")
						for _, host := range mysqlHost {
							smysql.DropUser(mc, username, host)
							w.Write([]byte("Deleting host from mysql for user:\n"))
							w.Write([]byte(username + "@" + host + "\n\n"))
						}
						mysqlUser := smysql.Search(mc, username)
						w.Write([]byte("These are the mysql host entries for the user:\n\n"))
						w.Write([]byte(username + ": " + mysqlUser + "\n\n"))
					}
				})
			})
		})
	})

	http.ListenAndServe(":3333", r)
}
