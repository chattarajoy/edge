package httpproxy

import (
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (writer *statusWriter) WriteHeader(status int) {
	writer.status = status
	writer.ResponseWriter.WriteHeader(status)
}

func (writer *statusWriter) Write(b []byte) (int, error) {
	if writer.status == 0 {
		writer.status = 200
	}

	length, err := writer.ResponseWriter.Write(b)
	writer.length += length
	return length, err

}

func(server *httpServer) dbUpdateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
		// update database
		next.ServeHTTP(writer, request)
	})
}

func(server *httpServer) recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
		defer func(){
			if err := recover(); err != nil {
				_ = server.Logger.Log("Error Occurred in handler: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte{})
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func (server *httpServer) logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		statWriter := statusWriter{ResponseWriter: writer}

		defer func(begin time.Time, request *http.Request) {
			_ = server.Logger.Log("host", request.Host, "path", request.URL.Path, "remote",
				request.RemoteAddr, "method", request.Method, "status", statWriter.status,
				"content-length", statWriter.length, "took", time.Since(begin))
		}(time.Now(), request)

		next.ServeHTTP(&statWriter, request)
	})
}