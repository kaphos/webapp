package webapp

// APIServer contains the data of an OpenAPI-spec server.
type APIServer struct {
	URL         string
	Description string
}

// GenDocs writes an OpenAPI documentation in JSON at the provided filename.
// "servers" is used just to decorate the file (as part of the OpenAPI spec,
// rather than being functional).
func (s *Server) GenDocs(servers []APIServer, filename string) error {
	for _, server := range servers {
		s.apiDocs.AddServer(server.URL, server.Description)
	}

	return s.apiDocs.Write(filename)
}
