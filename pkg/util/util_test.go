package util

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"github.com/stretchr/testify/assert"
)

const flytectlReleaseURL = "/repos/flyteorg/flytectl/releases/latest"
const baseURL = "https://api.github.com"
const wrongBaseURL = "htts://api.github.com"
const testVersion = "v0.1.20"

func TestGetRequest(t *testing.T) {
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
	})
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := GetRequest(wrongBaseURL, flytectlReleaseURL)
		assert.NotNil(t, err)
	})
	t.Run("Get request with 400", func(t *testing.T) {
		_, err := GetRequest("https://github.com", "/flyteorg/flyte/releases/download/latest/flyte_eks_manifest.yaml")
		assert.NotNil(t, err)
	})
}

func TestParseGithubTag(t *testing.T) {
	t.Run("Parse Github tag with success", func(t *testing.T) {
		response, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		data, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		tag, err := ParseGithubTag(data)
		assert.Nil(t, err)
		assert.Contains(t, tag, "v")
	})
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := ParseGithubTag([]byte("string"))
		assert.NotNil(t, err)
	})
}

func TestWriteIntoFile(t *testing.T) {
	t.Run("Successfully write into a file", func(t *testing.T) {
		response, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		data, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		err = WriteIntoFile(data, "version.yaml")
		assert.Nil(t, err)
	})
	t.Run("Error in writing file", func(t *testing.T) {
		response, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		data, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		err = WriteIntoFile(data, "/githubtest/version.yaml")
		assert.NotNil(t, err)
	})
}

func TestSetupFlyteDir(t *testing.T) {
	assert.Nil(t, SetupFlyteDir())
}

func TestIsVersionGreaterThan(t *testing.T) {
	t.Run("Compare flytectl version when upgrade available", func(t *testing.T) {
		_, err := IsVersionGreaterThan("v1.1.21", testVersion)
		assert.Nil(t, err)
	})
	t.Run("Compare flytectl version greater then", func(t *testing.T) {
		ok, err := IsVersionGreaterThan("v1.1.21", testVersion)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	})
	t.Run("Compare flytectl version smaller then", func(t *testing.T) {
		ok, err := IsVersionGreaterThan("v0.1.19", testVersion)
		assert.Nil(t, err)
		assert.Equal(t, false, ok)
	})
	t.Run("Compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan(testVersion, testVersion)
		assert.Nil(t, err)
	})
	t.Run("Error in compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan("vvvvvvvv", testVersion)
		assert.NotNil(t, err)
	})
	t.Run("Error in compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan(testVersion, "vvvvvvvv")
		assert.NotNil(t, err)
	})
}

func TestProgressBarForFlyteStatus(t *testing.T) {
	t.Run("Progress bar success", func(t *testing.T) {
		count := make(chan int64)
		go func() {
			time.Sleep(1 * time.Second)
			count <- 1
		}()
		ProgressBarForFlyteStatus(1, count, "")
	})
}

func TestGetLatestVersion(t *testing.T) {
	t.Run("Get latest release with wrong url", func(t *testing.T) {
		tag, err := GetLatestVersion("h://api.github.com/repos/flyteorg/flytectreleases/latest")
		assert.NotNil(t, err)
		assert.Equal(t, len(tag), 0)
	})
	t.Run("Get latest release", func(t *testing.T) {
		tag, err := GetLatestVersion(FlytectlReleasePath)
		assert.NotNil(t, err)
		assert.Equal(t, len(tag), 0)
	})
}

func TestDetectNewVersion(t *testing.T) {
	stdlibversion.Version = "v0.2.10"
	assert.Nil(t, DetectNewVersion(context.Background()))
	stdlibversion.Version = "v100.0.0"
	assert.Nil(t, DetectNewVersion(context.Background()))
	stdlibversion.Version = "v100.0.0"
	assert.Nil(t, DetectNewVersion(context.Background()))
	stdlibversion.Version = "v0"
	assert.NotNil(t, DetectNewVersion(context.Background()))
}
