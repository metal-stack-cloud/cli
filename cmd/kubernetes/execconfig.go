package kubernetes

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	c "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	configlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	configv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

type ExecCache struct {
	cachedir string
	fs       afero.Fs
}

func NewExecCache(fs afero.Fs, cachedir string) *ExecCache {
	return &ExecCache{
		cachedir: cachedir,
		fs:       fs,
	}
}

func NewUserExecCache(fs afero.Fs) (*ExecCache, error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	return NewExecCache(fs, cachedir), nil
}

func (ec *ExecCache) Clean(clusterid string) error {
	return ec.fs.Remove(ec.cacheFilePath(clusterid))
}

func (ec *ExecCache) cacheFilePath(clusterid string) string {
	return path.Join(ec.cachedir, fmt.Sprintf("metal_%s.json", clusterid))
}
func (ec *ExecCache) LoadCachedCredentials(clusterid string) (*c.ExecCredential, error) {
	credFile, err := ec.fs.Open(ec.cacheFilePath(clusterid))
	if err != nil {
		return nil, err
	}
	cachedCredentials, err := io.ReadAll(credFile)
	if err != nil {
		return nil, err
	}
	var execCredential c.ExecCredential
	err = json.Unmarshal(cachedCredentials, &execCredential)
	if err != nil {
		return nil, err
	}
	if execCredential.Status == nil || execCredential.Status.ExpirationTimestamp == nil {
		return nil, fmt.Errorf("cached credentials are invalid")
	}
	if execCredential.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		return nil, nil
	}
	return &execCredential, nil
}

func (ec *ExecCache) saveCachedCredentials(clusterid string, execCredential *c.ExecCredential) error {
	cachedCredentials, err := json.Marshal(execCredential)
	if err != nil {
		return fmt.Errorf("unable to marshal cached credentials: %w", err)
	}
	f, err := ec.fs.OpenFile(ec.cacheFilePath(clusterid), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("unable to write cached credentials: %w", err)
	}
	_, err = f.Write(cachedCredentials)
	return err
}

func (ec *ExecCache) ExecConfig(clusterid string, kubeRaw string, exp time.Duration) (*c.ExecCredential, error) {
	kubeconfig := &configv1.Config{}
	err := runtime.DecodeInto(configlatest.Codec, []byte(kubeRaw), kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode kubeconfig: %w", err)
	}
	if len(kubeconfig.AuthInfos) != 1 {
		return nil, fmt.Errorf("expected 1 auth info, got %d", len(kubeconfig.AuthInfos))
	}
	ai := kubeconfig.AuthInfos[0]
	expiration := metav1.NewTime(time.Now().Add(exp))
	ed := c.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1", // since k8s 1.24, if earlier versions are used, the API version is client.authentication.k8s.io/v1beta1
			Kind:       "ExecCredential",
		},
		Status: &c.ExecCredentialStatus{
			ClientCertificateData: string(ai.AuthInfo.ClientCertificateData),
			ClientKeyData:         string(ai.AuthInfo.ClientKeyData),
			ExpirationTimestamp:   &expiration,
		},
	}
	// ignoring error, so a failed save doesn't break the flow
	_ = ec.saveCachedCredentials(clusterid, &ed)
	return &ed, nil
}
