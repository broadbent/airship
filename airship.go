package main



import (
    "./ripple"
    "net/http"
)
	

func main() {
    // Build the REST application

    app := ripple.NewApplication()

    // Create a controller and register it. Any number of controllers
    // can be registered that way.

    userController := ripple.NewUserController()
    app.RegisterController("users", userController)

    // Setup the routes. The special patterns `_controller` will automatically match
    // an existing controller, as defined above. Likewise, `_action` will match any 
    // existing action.

    app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
    app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
    app.AddRoute(ripple.Route{ Pattern: ":_controller" })

    // Start the server

    http.ListenAndServe(":8080", app)
}

