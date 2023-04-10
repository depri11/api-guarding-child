package gc

func (g *GC) RegisterRouter() {
	publicApiRouter := g.Router.PathPrefix("/api/v1/public").Subrouter()
	protectedApiRouter := g.Router.PathPrefix("/api/v1").Subrouter()

	g.RegisterStatussRouter(publicApiRouter, protectedApiRouter)
}
