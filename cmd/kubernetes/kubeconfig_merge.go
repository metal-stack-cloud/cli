package kubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	configlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	configv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

type MergedKubeconfig struct {
	Raw         []byte
	Path        string
	ContextName string
}

func MergeKubeconfig(fs afero.Fs, raw []byte, kubeconfigPath, projectName *string) (*MergedKubeconfig, error) {
	path := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if kubeconfigPath != nil {
		path = *kubeconfigPath
	}
	if path == "" {
		path = clientcmd.RecommendedHomeFile
	}

	if strings.Contains(path, ":") {
		return nil, fmt.Errorf("it is currently not supported to merge when multiple kubeconfigs are provided")
	}

	if _, err := fs.Stat(path); os.IsNotExist(err) {
		err := afero.WriteFile(fs, path, nil, 0600)
		if err != nil {
			return nil, fmt.Errorf("error to write to: %w", err)
		}
	}

	currentConfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %w", err)
	}

	kubeconfig := &configv1.Config{}
	err = runtime.DecodeInto(configlatest.Codec, raw, kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode kubeconfig: %w", err)
	}

	var (
		authInfo    *configv1.AuthInfo
		cluster     *configv1.NamedCluster
		clusterName string
	)

	for _, cl := range kubeconfig.Clusters {
		cl := cl

		prefix, _, found := strings.Cut(cl.Name, "-external")
		if !found {
			continue
		}
		_, clusterName, found = strings.Cut(prefix, "--")
		if !found {
			continue
		}

		cluster = &cl
	}

	for _, a := range kubeconfig.AuthInfos {
		a := a

		if !strings.HasSuffix(a.Name, clusterName+"-external") {
			continue
		}

		authInfo = &a.AuthInfo
	}

	if authInfo == nil || cluster == nil || clusterName == "" {
		return nil, fmt.Errorf("internal error: kubeconfig does not contain all required information, please update client or raise ticket on metalstack.cloud")
	}

	contextName := fmt.Sprintf("%s@metalstack.cloud", clusterName)
	if projectName != nil {
		contextName = fmt.Sprintf("%s-%s@metalstack.cloud", clusterName, *projectName)
	}

	currentConfig.Contexts[contextName] = &api.Context{
		Cluster:  contextName,
		AuthInfo: contextName,
	}
	currentConfig.Clusters[contextName] = &api.Cluster{
		Server:                   cluster.Cluster.Server,
		CertificateAuthorityData: cluster.Cluster.CertificateAuthorityData,
	}
	currentConfig.AuthInfos[contextName] = &api.AuthInfo{
		ClientCertificateData: authInfo.ClientCertificateData,
		ClientKeyData:         authInfo.ClientKeyData,
	}

	if currentConfig.CurrentContext == "" {
		currentConfig.CurrentContext = contextName
	}

	merged, err := runtime.Encode(configlatest.Codec, currentConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to encode kubeconfig: %w", err)
	}

	return &MergedKubeconfig{
		Raw:         merged,
		ContextName: contextName,
		Path:        path,
	}, nil
}
