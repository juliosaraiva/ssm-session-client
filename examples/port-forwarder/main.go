package main

import (
	"flag"
	"os"
	"log"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/juliosaraiva/ssm-session-client/ssmclient"
)

// Start a SSM port forwarding session.
// Usage: port-forwarder [profile_name] target_spec
//   The profile_name argument is the name of profile in the local AWS configuration to use for credentials.
//   if unset, it will consult the AWS_PROFILE environment variable, and if that is unset, will use credentials
//   set via environment variables, or from the default profile.
//
//   The target_spec parameter is required, and is in the form of ec2_instance_id:port_number (ex: i-deadbeef:80)

func main() {
	var profile, instanceName string 
	var srcPort, dstPort int

	flag.StringVar(&profile, "profile", "default", "AWS Profile")
	flag.StringVar(&instanceName, "instance-name", "JenkinsServerPrivateNetworkTestFromSnapshot", "AWS Instance Name")
	flag.IntVar(&srcPort, "src-port", 8080, "Instance Port to Connect to")
	flag.IntVar(&dstPort, "dst-port", 50000, "Local port")

	flag.Parse()

	if v, ok := os.LookupEnv("AWS_PROFILE"); ok {
		profile = v
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatal(err)
	}

	tgt, err := ssmclient.ResolveTarget(instanceName, cfg)
	if err != nil {
		log.Fatal(err)
	}


	in := ssmclient.PortForwardingInput{
		Target:     tgt,
		RemotePort: srcPort,
		LocalPort:  dstPort, // just use random port for demo purposes (this is the default, if not set > 0)
	}
	log.Fatal(ssmclient.PortForwardingSession(cfg, &in))
}