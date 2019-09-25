package httpproxy

import (
	"github.com/chattarajoy/edge/pkg/mysql"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
	requestID string
}

func (writer *statusWriter) WriteHeader(status int) {
	writer.status = status
	writer.ResponseWriter.WriteHeader(status)
}

func (writer *statusWriter) Write(b []byte) (int, error) {
	if writer.status == 0 {
		writer.status = 200

	}

	if writer.requestID == "" {
		writer.requestID = uuid.New().String()
	}

	length, err := writer.ResponseWriter.Write(b)
	writer.length += length
	return length, err

}

func(server *httpServer) dbUpdateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
		defer func(request *http.Request, server *httpServer){
			dbLog := &mysql.DbLog{
				Path: request.URL.Path,
			}
			err := dbLog.Insert(server.Db)
			if err != nil {
				panic(err.Error())
			}
		}(request, server)
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
			// _ = server.Logger.Log("Headers: ", strings.Join(helpers.GetHeaders(request), ", "))
			_ = server.Logger.Log("requestID", statWriter.requestID, "host", request.Host, "path", request.URL.Path, "remote",
				request.RemoteAddr, "method", request.Method, "status", statWriter.status,
				"content-length", statWriter.length, "took", time.Since(begin))
		}(time.Now(), request)

		next.ServeHTTP(&statWriter, request)
	})
}

func (server *httpServer) wrappedHandlers(h http.Handler) http.Handler {
	return server.recoverHandler(server.logHandler(server.dbUpdateHandler(h)))
}
