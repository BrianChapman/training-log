package main

import (
	"appengine"
	"appengine/memcache"
	"github.com/emicklei/go-restful"
	// "github.com/emicklei/go-restful/swagger"
	"net/http"
)

// This example is functionally the same as ../restful-user-service.go
// but it`s supposed to run on Goole App Engine (GAE)
//
// contributed by ivanhawkes

type User struct {
	Id, Name string
}

type UserService struct {
	// normally one would use DAO (data access object)
	// but in this example we simple use memcache.
}

type Activity struct {
	Id       string
	Title    string
	Distance int32 //in meters, we don't worry about fractions of a meter.
}

type ActivityService struct {
}

func (u UserService) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.PATCH("").To(u.updateUser).
		// docs
		Doc("update a user").
		Reads(User{})) // from the request

	ws.Route(ws.PUT("/{user-id}").To(u.createUser).
		// docs
		Doc("create a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))
}

func (a ActivityService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/activities").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/").To(a.listActivities).
		Doc("lists activities").
		Writes([]Activity{}))

	ws.Route(ws.GET("/{activity-id}").To(a.findActivity).
		// docs
		Doc("get an activity").
		Param(ws.PathParameter("activity-id", "identifier of the activity").DataType("string")).
		Writes(Activity{})) // on the response

	ws.Route(ws.PATCH("").To(a.updateActivity).
		// docs
		Doc("update an activity").
		Reads(Activity{})) // from the request

	ws.Route(ws.PUT("/{activity-id}").To(a.createActivity).
		// docs
		Doc("create an activity").
		Param(ws.PathParameter("activity-id", "identifier of the activity").DataType("string")).
		Reads(Activity{})) // from the request

	ws.Route(ws.DELETE("/{activity-id}").To(a.removeActivity).
		// docs
		Doc("delete an activity").
		Param(ws.PathParameter("activity-id", "identifier of the activity").DataType("string")))

	restful.Add(ws)
}

// GET http://localhost:8080/users/1
//
func (u UserService) findUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("user-id")
	usr := new(User)
	_, err := memcache.Gob.Get(c, id, &usr)
	if err != nil || len(usr.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PATCH http://localhost:8080/users
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserService) updateUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		item := &memcache.Item{
			Key:    usr.Id,
			Object: &usr,
		}
		err = memcache.Gob.Set(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa</Name></User>
//
func (u *UserService) createUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	usr := User{Id: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		item := &memcache.Item{
			Key:    usr.Id,
			Object: &usr,
		}
		err = memcache.Gob.Add(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserService) removeUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("user-id")
	err := memcache.Delete(c, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// GET http://localhost:8080/activities/1
//
func (a ActivityService) findActivity(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("activity-id")
	activity := new(Activity)
	_, err := memcache.Gob.Get(c, id, &activity)
	if err != nil || len(activity.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "Activity could not be found.")
	} else {
		response.WriteEntity(activity)
	}
}

// GET /activities
//
func (a ActivityService) listActivities(request *restful.Request, response *restful.Response) {
	response.WriteErrorString(http.StatusNotFound, "Listing Activities is not Implemented.")

}

// PATCH http://localhost:8080/activities
// <Activity><Id>1</Id><Title>Tuesday Sprints</Title><Distance>20000<Distance/></Activity>
//
func (a *ActivityService) updateActivity(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	activity := new(Activity)
	err := request.ReadEntity(&activity)
	if err == nil {
		item := &memcache.Item{
			Key:    activity.Id,
			Object: &activity,
		}
		err = memcache.Gob.Set(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteEntity(activity)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/activities/1
// <User><Id>1</Id><Name>Melissa</Name></User>
//
func (a *ActivityService) createActivity(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	activity := Activity{Id: request.PathParameter("activity-id")}
	err := request.ReadEntity(&activity)
	if err == nil {
		item := &memcache.Item{
			Key:    activity.Id,
			Object: &activity,
		}
		err = memcache.Gob.Add(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(activity)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/activities/1
//
func (a *ActivityService) removeActivity(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("activity-id")
	err := memcache.Delete(c, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func getGaeURL() string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080"
	} else {
		/**
		 * Include your URL on App Engine here.
		 * I found no way to get AppID without appengine.Context and this always
		 * based on a http.Request.
		 */
		return "http://training-log.appspot.com"
	}
}

func init() {
	u := UserService{}
	u.Register()
	a := ActivityService{}
	a.Register()

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	// config := swagger.Config{
	// 	WebServices:    restful.RegisteredWebServices(), // you control what services are visible
	// 	WebServicesUrl: getGaeURL(),
	// 	ApiPath:        "/apidocs.json",

	// 	// Optionally, specifiy where the UI is located
	// 	SwaggerPath: "/apidocs/",
	// 	// GAE support static content which is configured in your app.yaml.
	// 	// This example expect the swagger-ui in static/swagger so you should place it there :)
	// 	SwaggerFilePath: "static/swagger"}
	// swagger.InstallSwaggerService(config)
}
