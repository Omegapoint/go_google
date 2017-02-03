package stuff

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "golang.org/x/net/context"   //for datastore
	_ "golang.org/x/oauth2/google" //for Datastore

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	_ "github.com/go-sql-driver/mysql" //Needed for SQL connection
	_ "google.golang.org/appengine/remote_api"
)

func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/sup", what)

}
func hello(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r) //For Google logs
	log.Debugf(ctx, "Start")
	io.WriteString(w, "hello interweb")
	io.WriteString(w, "twice")
	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := os.Getenv("CLOUDSQL_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/go_test", user, password, connectionName))
	if err != nil {
		log.Debugf(ctx, "Open")
		log.Errorf(ctx, "%v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Debugf(ctx, "Ping")
		log.Errorf(ctx, "%v", err)
	}
	var name string
	var quantity int
	rows, err := db.Query("SELECT entry, quantity FROM go_data")
	if err != nil {
		log.Debugf(ctx, "Query")
		log.Errorf(ctx, "%v", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &quantity)
		if err != nil {
			log.Debugf(ctx, "Scan")
			log.Errorf(ctx, "%v", err)
		}
		result := fmt.Sprintf("%s, %d", name, quantity)
		io.WriteString(w, result)
	}
}

func what(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "sup")
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Printf("Environment variable: %s not set", k)
	}
	return v
}
