package stdlibtemplate

import (
	"errors"

	"github.com/luigimorel/gogen/internal/fsutils"
)

func CreateRouterSetup() error {
	job, err := fsutils.NewJobFromJSON([]byte(routerPkgTemplate))
	switch {
	case err != nil:
		return err
	case job == nil:
		return errors.New("job is nil")
	}
	return job.Execute()
}
