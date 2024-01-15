package helpers

import (
	"context"
	"fmt"
	"io"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"

	metalssh "github.com/metal-stack/metal-lib/pkg/ssh"
	metalvpn "github.com/metal-stack/metal-lib/pkg/vpn"
)

func SSHViaVPN(out io.Writer, machineID string, creds *adminv1.ClusterServiceCredentialsResponse) error {
	if creds.SshKeypair == nil || len(creds.SshKeypair.Privatekey) == 0 {
		return fmt.Errorf("no ssh key found")
	}
	if creds.Vpn == nil || creds.Vpn.Authkey == "" {
		return fmt.Errorf("no vpn connection possible")
	}

	fmt.Fprintf(out, "accessing firewall through vpn at %s ", creds.Vpn.Address)

	ctx := context.Background()
	v, err := metalvpn.Connect(ctx, machineID, creds.Vpn.Address, creds.Vpn.Authkey)
	if err != nil {
		return err
	}
	defer v.Close()

	s, err := metalssh.NewClientWithConnection("metal", v.TargetIP, creds.SshKeypair.Privatekey, v.Conn)
	if err != nil {
		return err
	}

	return s.Connect(nil)
}
