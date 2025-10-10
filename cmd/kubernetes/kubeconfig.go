package kubernetes

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	configlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	configv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

type authType string

const (
	AuthTypeClientCerts authType = "certs"
	AuthTypeExec        authType = "exec"
)

type Kubeconfig struct {
	Raw         []byte
	Path        string
	ContextName string
}

func NewKubeconfigFromRaw(fs afero.Fs, in io.Reader, out io.Writer, raw []byte, projectName *string, projectid, clusterid string) (*Kubeconfig, error) {
	path := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if userPath := viper.GetString("kubeconfig"); userPath != "" {
		path = userPath
	}
	if path == "" {
		path = clientcmd.RecommendedHomeFile
	}

	if strings.Contains(path, ":") {
		return nil, fmt.Errorf("it is currently not supported to merge when multiple kubeconfigs are provided")
	}

	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if file exists: %w", err)
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

	currentConfig := &api.Config{
		Clusters:  map[string]*api.Cluster{},
		Contexts:  map[string]*api.Context{},
		AuthInfos: map[string]*api.AuthInfo{},
	}

	if exists {
		if viper.GetBool("merge") {
			var err error
			currentConfig, err = clientcmd.LoadFromFile(path)
			if err != nil {
				return nil, fmt.Errorf("error loading kubeconfig: %w", err)
			}
		} else {
			err = genericcli.PromptCustom(&genericcli.PromptConfig{
				Message:     fmt.Sprintf(color.YellowString("There is already a file at %s. In combination with --merge=false this file will be overwritten. Are you sure you want to continue?"), path),
				ShowAnswers: true,
				In:          in,
				Out:         out,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	currentConfig.Contexts[contextName] = &api.Context{
		Cluster:  contextName,
		AuthInfo: contextName,
	}

	currentConfig.Clusters[contextName] = &api.Cluster{
		Server:                   cluster.Cluster.Server,
		CertificateAuthorityData: cluster.Cluster.CertificateAuthorityData,
	}

	if currentConfig.CurrentContext == "" {
		currentConfig.CurrentContext = contextName
	}

	auth := AuthTypeExec
	if viper.GetString("auth-type") != "" {
		auth = authType(viper.GetString("auth-type"))
	}

	switch auth {
	case AuthTypeExec:
		metalcli, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("unable to get executable path: %w", err)
		}

		currentConfig.AuthInfos[contextName] = &api.AuthInfo{
			Exec: &api.ExecConfig{
				Command:         metalcli,
				Args:            []string{"cluster", "exec-config", "-p", projectid, clusterid},
				APIVersion:      "client.authentication.k8s.io/v1",
				InteractiveMode: api.IfAvailableExecInteractiveMode,
			},
		}

		ec, err := NewUserExecCache(fs)
		if err != nil {
			return nil, err
		}
		// remove cached credentials so a new one will be created
		_ = ec.Clean(clusterid)
	case AuthTypeClientCerts:
		currentConfig.AuthInfos[contextName] = &api.AuthInfo{
			ClientCertificateData: authInfo.ClientCertificateData,
			ClientKeyData:         authInfo.ClientKeyData,
		}
	default:
		return nil, fmt.Errorf("unsupported auth type for kubeconfig: %s", auth)
	}

	merged, err := runtime.Encode(configlatest.Codec, currentConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to encode kubeconfig: %w", err)
	}

	return &Kubeconfig{
		Raw:         merged,
		ContextName: contextName,
		Path:        path,
	}, nil
}
