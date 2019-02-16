package pyruletest

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include "greeter.h"
import "C"
import (
	"github.com/project-flogo/rules/common/model"
	"fmt"
	"context"
)


