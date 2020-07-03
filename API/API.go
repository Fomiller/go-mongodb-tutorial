package API

import (
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/config"
)

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	config.TPL.ExecuteTemplate(res, "index.html", nil)
}
