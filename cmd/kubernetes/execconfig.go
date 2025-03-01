package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	c "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	configlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	configv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

func cacheFilePath(cachedir, clusterid string) string {
	return path.Join(cachedir, fmt.Sprintf("metal_%s.json", clusterid))
}
func LoadCachedCredentials(clusterid string) (*c.ExecCredential, error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	cachedCredentials, err := os.ReadFile(cacheFilePath(cachedir, clusterid))
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

func saveCachedCredentials(clusterid string, execCredential *c.ExecCredential) error {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cachedCredentials, err := json.Marshal(execCredential)
	if err != nil {
		return fmt.Errorf("unable to marshal cached credentials: %w", err)
	}
	err = os.WriteFile(cacheFilePath(cachedir, clusterid), cachedCredentials, 0600)
	if err != nil {
		return fmt.Errorf("unable to write cached credentials: %w", err)
	}
	return nil
}

func ExecConfig(clusterid string, kubeRaw string, exp time.Duration) (*c.ExecCredential, error) {
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
			APIVersion: "client.authentication.k8s.io/v1", // since k8s 1.22, if earlier versions are used, the API version is client.authentication.k8s.io/v1beta1
			Kind:       "ExecCredential",
		},
		Status: &c.ExecCredentialStatus{
			ClientCertificateData: string(ai.AuthInfo.ClientCertificateData),
			ClientKeyData:         string(ai.AuthInfo.ClientKeyData),
			ExpirationTimestamp:   &expiration,
		},
	}
	// ignoring error, so a failed save doesn't break the flow
	_ = saveCachedCredentials(clusterid, &ed)
	return &ed, nil
}
