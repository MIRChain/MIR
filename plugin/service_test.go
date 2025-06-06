package plugin

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	testifyassert "github.com/stretchr/testify/assert"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

func typicalPluginManager(t *testing.T) *PluginManager[nist.PrivateKey,nist.PublicKey] {
	testObject, err := NewPluginManager[nist.PrivateKey,nist.PublicKey]("arbitraryName", &Settings{
		Providers: map[PluginInterfaceName]PluginDefinition{
			HelloWorldPluginInterfaceName: {
				Name:    "arbitrary-helloWorld",
				Version: "1.0.0",
				Config:  "arbitrary config",
			},
		},
	}, false, false, "")

	testifyassert.NoError(t, err)
	return testObject
}

func TestPluginManager_ProvidersPopulation(t *testing.T) {
	arbitraryPluginInterfaceName := PluginInterfaceName("arbitrary")
	defer func() {
		delete(pluginProviders, arbitraryPluginInterfaceName)
	}()
	pluginProviders[arbitraryPluginInterfaceName] = pluginProvider[nist.PrivateKey,nist.PublicKey]{
		pluginSet: plugin.PluginSet{},
	}

	testObject, err := NewPluginManager[nist.PrivateKey,nist.PublicKey]("arbitraryName", &Settings{
		Providers: map[PluginInterfaceName]PluginDefinition{
			HelloWorldPluginInterfaceName: {
				Name:    "arbitrary-helloWorld",
				Version: "1.0.0",
				Config:  "arbitrary config",
			},
			arbitraryPluginInterfaceName: {
				Name:    "foo-bar",
				Version: "2.0.0",
				Config:  "arbitrary config",
			},
		},
	}, false, false, "")

	testifyassert.NoError(t, err)
	testifyassert.Equal(t, fmt.Sprintf("arbitrary-helloWorld-1.0.0-%s-%s", runtime.GOOS, runtime.GOARCH), testObject.initializedPlugins[HelloWorldPluginInterfaceName].(*basePlugin[nist.PrivateKey,nist.PublicKey]).pluginDefinition.FullName())
	testifyassert.Equal(t, fmt.Sprintf("foo-bar-2.0.0-%s-%s", runtime.GOOS, runtime.GOARCH), testObject.initializedPlugins[arbitraryPluginInterfaceName].(*basePlugin[nist.PrivateKey,nist.PublicKey]).pluginDefinition.FullName())
}

func TestPluginManager_GetPluginTemplate_whenTypical(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)

	p := new(HelloWorldPluginTemplate[nist.PrivateKey,nist.PublicKey])
	err := testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, p)

	assert.NoError(err)
	assert.NotNil(p)
}

func TestPluginManager_GetPlugin_whenReadFromCache(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)
	p := new(HelloWorldPluginTemplate[nist.PrivateKey,nist.PublicKey])
	err := testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, p)
	assert.NoError(err)
	assert.NotNil(p)

	actual, ok := testObject.getPlugin(HelloWorldPluginInterfaceName)

	assert.True(ok)
	assert.Equal(p, actual)
}

func TestPluginManager_GetPlugin_whenReadFromInitializedPluginsCache(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)

	actual, ok := testObject.getPlugin(HelloWorldPluginInterfaceName)

	assert.True(ok)
	assert.IsType(new(basePlugin[nist.PrivateKey,nist.PublicKey]), actual)
}

func TestPluginManager_GetPluginTemplate_whenReadFromCache(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)
	p := new(HelloWorldPluginTemplate[nist.PrivateKey,nist.PublicKey])
	err := testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, p)
	assert.NoError(err)
	assert.NotNil(p)

	actual := new(HelloWorldPluginTemplate[nist.PrivateKey,nist.PublicKey])
	err = testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, actual)

	assert.NoError(err)
	assert.Equal(p, actual)
}

func TestPluginManager_GetPluginTemplate_whenPluginTemplateNotExtendBasePlugin(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)

	invalid := new(invalidPluginTemplate)
	err := testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, invalid)

	t.Log(err)
	assert.Error(err)
}

func TestPluginManager_GetPluginTemplate_whenPluginTemplateNotExtendPointerBasePlugin(t *testing.T) {
	assert := testifyassert.New(t)
	testObject := typicalPluginManager(t)

	invalid := new(invalidPluginTemplateNoPointer[nist.PrivateKey,nist.PublicKey])
	err := testObject.GetPluginTemplate(HelloWorldPluginInterfaceName, invalid)

	t.Log(err)
	assert.Error(err)
}

type invalidPluginTemplateNoPointer [T crypto.PrivateKey, P crypto.PublicKey]struct {
	basePlugin[T,P]
}

type invalidPluginTemplate struct {
	someField int
}

func (i invalidPluginTemplate) Start() error {
	panic("implement me")
}

func (i invalidPluginTemplate) Stop() error {
	panic("implement me")
}

func (i invalidPluginTemplate) Info() (PluginInterfaceName, interface{}) {
	panic("implement me")
}
