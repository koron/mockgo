package mockrt3

import (
	"context"

	"github.com/google/go-cmp/cmp/cmpopts"
)

// IgnoreContext is an option to ignore context.Context.
var IgnoreContext = cmpopts.IgnoreInterfaces(struct{ context.Context }{})
