package application

type Application struct {
	Services *Services
}

type Services struct {
	Networks NetworkService
}

func New(services *Services) *Application {
	app := &Application{services}
	return app
}
