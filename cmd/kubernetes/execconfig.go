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

func cacheFilePath(cachedir, clusterid string) string {
	return path.Join(cachedir, fmt.Sprintf("metal_%s.json", clusterid))
}
func LoadCachedCredentials(fs afero.Fs, clusterid string) (*c.ExecCredential, error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	credFile, err := fs.Open(cacheFilePath(cachedir, clusterid))
	if err != nil {
		// if the file cannot be opened, we assume no error and return nil
		// so the caller will (re)create it
		return nil, nil //nolint:nilerr
	}
	cachedCredentials, err := io.ReadAll(credFile)
	if err != nil {
		return nil, nil //nolint:nilerr
	}
	var execCredential c.ExecCredential
	err = json.Unmarshal(cachedCredentials, &execCredential)
	if err != nil {
		return nil, nil //nolint:nilerr
	}
	if execCredential.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		return nil, nil
	}
	return &execCredential, nil
}

func saveCachedCredentials(fs afero.Fs, clusterid string, execCredential *c.ExecCredential) error {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cachedCredentials, err := json.Marshal(execCredential)
	if err != nil {
		return fmt.Errorf("unable to marshal cached credentials: %w", err)
	}
	f, err := fs.OpenFile(cacheFilePath(cachedir, clusterid), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("unable to write cached credentials: %w", err)
	}
	_, err = f.Write(cachedCredentials)
	return err
}

func ExecConfig(fs afero.Fs, clusterid string, kubeRaw string, exp time.Duration) (*c.ExecCredential, error) {
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
	_ = saveCachedCredentials(fs, clusterid, &ed)
	return &ed, nil
}
