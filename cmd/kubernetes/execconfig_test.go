package kubernetes

import (
	"testing"
	"time"

	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	c "k8s.io/client-go/pkg/apis/clientauthentication/v1"
)

func TestLoadCachedCredentials(t *testing.T) {
	type args struct {
		clusterid string
	}
	tests := []struct {
		name      string
		content   string
		args      args
		hasResult bool
		wantErr   bool
	}{
		{
			name: "no existing cache file",
			args: args{
				clusterid: "cid",
			},
			wantErr: true,
		},
		{
			name:    "existing cache file",
			content: "{}",
			args: args{
				clusterid: "cid",
			},
			wantErr: true,
		},
		{
			name:    "data with expired date in cache file",
			content: `{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1","spec":{"interactive":false},"status":{"expirationTimestamp":"2020-03-14T18:52:24Z"}}`,
			args: args{
				clusterid: "cid",
			},
		},
		{
			name:    "data with non expired date in cache file",
			content: `{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1","spec":{"interactive":false},"status":{"expirationTimestamp":"2050-03-14T18:52:24Z"}}`,
			args: args{
				clusterid: "cid",
			},
			hasResult: true,
		},
	}
	fs := afero.NewMemMapFs()
	ec := NewExecCache(fs, "/tmp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = afero.WriteFile(fs, ec.cacheFilePath(tt.args.clusterid), []byte(tt.content), 0644)
			got, err := ec.LoadCachedCredentials(tt.args.clusterid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCachedCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.hasResult {
				t.Errorf("LoadCachedCredentials wants result = %v", got)
			}
		})
	}
}

func TestExecCache_ExecConfig(t *testing.T) {
	type args struct {
		clusterid string
		kubeRaw   string
		exp       time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *c.ExecCredential
		wantErr bool
	}{
		{
			name: "illegal kubeconfig",
			args: args{
				clusterid: "clusterid",
				kubeRaw:   "no valid kubeconfig",
				exp:       10 * time.Minute,
			},
			wantErr: true,
		},
		{
			name: "multiple authinfos",
			args: args{
				clusterid: "clusterid",
				kubeRaw:   `{"users":[{"name":"user1", user: {"client-certificate-data":"Y2VydDE=", "client-key-data":"a2V5MQ=="}},{"name":"user2",user: {"client-certificate-data":"Y2VydDI=", "client-key-data":"a2V5Mg=="}}]}`,
				exp:       10 * time.Minute,
			},
			wantErr: true,
		},
		{
			name: "one authinfo",
			args: args{
				clusterid: "clusterid",
				kubeRaw:   `{"users":[{"name":"user1", user: {"client-certificate-data":"Y2VydDE=", "client-key-data":"a2V5MQ=="}}]}`,
				exp:       10 * time.Minute,
			},
			wantErr: false,
			want: &c.ExecCredential{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ExecCredential",
					APIVersion: "client.authentication.k8s.io/v1",
				},
				Spec: c.ExecCredentialSpec{
					Interactive: false,
				},
				Status: &c.ExecCredentialStatus{
					ClientCertificateData: "cert1",
					ClientKeyData:         "key1",
				},
			},
		},
	}

	fs := afero.NewMemMapFs()
	ec := NewExecCache(fs, "/tmp")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ec.ExecConfig(tt.args.clusterid, tt.args.kubeRaw, tt.args.exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecCache.ExecConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.Status.ClientCertificateData != tt.want.Status.ClientCertificateData {
					t.Errorf("ClientCertificateData = %v, want %v", got.Status.ClientCertificateData, tt.want.Status.ClientCertificateData)
				}
				if got.Status.ClientKeyData != tt.want.Status.ClientKeyData {
					t.Errorf("ClientKeyData = %v, want %v", got.Status.ClientKeyData, tt.want.Status.ClientKeyData)
				}
				// check if cache file was written
				fi, err := fs.Stat(ec.cacheFilePath(tt.args.clusterid))
				if err != nil {
					t.Errorf("cache file not found")
				} else {
					if fi.Size() == 0 {
						t.Errorf("cache file is empty")
					}
				}
			}
		})
	}
}
