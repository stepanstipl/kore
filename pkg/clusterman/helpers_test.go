package clusterman

import (
	"testing"

	"github.com/appvia/kore/pkg/clusterapp"
	"github.com/stretchr/testify/assert"
	cc "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestLoadAllManifests(t *testing.T) {
	options, err := clusterapp.GetClientOptions()
	assert.NoError(t, err)
	client := cc.NewFakeClientWithScheme(options.Scheme)
	err = LoadAllManifests(client)
	assert.NoError(t, err)
}
