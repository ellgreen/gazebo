package modules

import (
	"io/ioutil"
	"net/http"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/g"
)

// HTTP holds the definitions for the http module
var HTTP = &Module{
	Name: "http",
	Values: map[string]g.Object{
		"get": g.NewObject(func(args g.Args) g.Object {
			url := g.ToString(args.Self())

			resp, err := http.Get(url)
			assert.Nil(err)

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.Nil(err)

			response := g.NewObject("response")
			response.Attributes().Set("body", g.NewObject(string(body)))
			response.Attributes().Set("status", g.NewObject(resp.StatusCode))

			return response
		}),
	},
}
