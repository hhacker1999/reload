package reload

import (
	"os"
)

func (r *Reload) initialise() {
	os.Mkdir(r.root+"/.reload", 0775)
}
