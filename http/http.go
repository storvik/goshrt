package http

// logResWriterError logs errors when writing to http.ResponseWriter.
// Convenient function to avoid extensive error checking.
func (s *Server) logResWriterError(_ int, err error) {
	if err != nil {
		s.ErrorLog.Printf("Write failed: %v", err)
	}
}
