package mbgl_test

import (
	"testing"
	"github.com/go-spatial/go-mbgl/mbgl"
)

func TestThreadPool(t *testing.T) {
	tpool := mbgl.NewThreadPool(4)
	tpool.Destruct()
}