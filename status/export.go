package status

import (
	"github.com/caicloud/rudder-client/status/internal"
)

// export internal types

// Umpire can employs many assistant to handle many kinds of objects.
type Umpire = internal.Umpire

// Assistant handles a kind of object. It will generates the resourceStatus for the object.
type Assistant = internal.Assistant
