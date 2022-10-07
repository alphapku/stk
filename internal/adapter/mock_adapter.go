package adapter

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type MockAdapter struct {
	file *os.File
}

func (f *MockAdapter) Close() {
	f.file.Close()
}

func (f *MockAdapter) Read() (interface{}, error) {
	return nil, nil
}

func (f *MockAdapter) Do(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, struct{}{})
}
