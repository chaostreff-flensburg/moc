package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func HandleSQLError(err error) error {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			return BadRequestError("duplicate slug").WithInternalError(mysqlError)
		}
	} else if sqliteError, ok := err.(sqlite3.Error); ok {
		if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
			return BadRequestError("duplicate slug").WithInternalError(mysqlError)
		}
	}

	return InternalServerError("internal server error").WithInternalError(err)
}

func SendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
