package helpers

import (
	"context"
	"fmt"
	"io"
	"os"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"

	metalssh "github.com/metal-stack/metal-lib/pkg/ssh"
	metalvpn "github.com/metal-stack/metal-lib/pkg/vpn"
)

func FirewallSSHViaVPN(out io.Writer, machineID string, vpn *adminv1.VPN, sshPrivateKey []byte) (err error) {
	fmt.Fprintf(out, "accessing firewall through vpn at %s ", vpn.Address)

	ctx := context.Background()
	v, err := metalvpn.Connect(ctx, machineID, vpn.Address, vpn.Authkey)
	if err != nil {
		return err
	}
	defer v.Close()

	s, err := metalssh.NewClientWithConnection("metal", v.TargetIP, sshPrivateKey, v.Conn)
	if err != nil {
		return err
	}

	return s.Connect(nil)
}

// sshClient opens an interactive ssh session to the host on port with user, authenticated by the key.
func sshClient(user, keyfile, host string, port int, idToken *string) error {
	privateKey, err := os.ReadFile(keyfile)
	if err != nil {
		return err
	}

	s, err := metalssh.NewClient(user, host, privateKey, port)
	if err != nil {
		return err
	}

	var env *metalssh.Env
	if idToken != nil {
		env = &metalssh.Env{"LC_METAL_STACK_OIDC_TOKEN": *idToken}
	}

	return s.Connect(env)
}
