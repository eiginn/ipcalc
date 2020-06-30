package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xgfone/netaddr"
	"log"
	"os"
	"text/tabwriter"
)

func handleInput(c *cli.Context) error {
	var argzero string
	var net netaddr.IPNetwork

	argzero = c.Args().Get(0)
	net, err := netaddr.NewIPNetwork(argzero)
	if err != nil {
		return err
	}
	err = outputTable(net)
	if err != nil {
		return err
	}
	return nil
}

func bits(addr netaddr.IPAddress, mask int, host bool) string {
	var padding int
	switch {
	case mask < 8:
		padding = 0
	case mask < 16:
		padding = 1
	case mask < 24:
		padding = 2
	default:
		padding = 3
	}
	if host {
		return addr.Bits()[mask+padding:]
	}
	return addr.Bits()[:mask+padding]
}

func outputTable(net netaddr.IPNetwork) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Address\t:\t%s\t%s %s\n", net.Address(), bits(net.Address(), net.Mask(), false), bits(net.Address(), net.Mask(), true))
	fmt.Fprintf(w, "Netmask\t:\t%s = %d\t%s %s\n", net.NetworkMask(), net.Mask(), bits(net.NetworkMask(), net.Mask(), false), bits(net.NetworkMask(), net.Mask(), true))
	fmt.Fprintf(w, "Wildcard\t:\t%s\t%s %s\n", net.HostMask(), bits(net.HostMask(), net.Mask(), false), bits(net.HostMask(), net.Mask(), true))
	// we need to have same number of fields to get padding right
	fmt.Fprintf(w, "=>\t \t \t \n")
	fmt.Fprintf(w, "Network\t:\t%s\t%s %s\n", net.CIDR(), bits(net.CIDR().Address(), net.Mask(), false), bits(net.CIDR().Address(), net.Mask(), true))

	fmt.Fprintf(w, "HostMin\t:\t%s\t%s %s\n", net.First(), bits(net.First(), net.Mask(), false), bits(net.First(), net.Mask(), true))
	fmt.Fprintf(w, "HostMax\t:\t%s\t%s %s\n", net.Last(), bits(net.Last(), net.Mask(), false), bits(net.Last(), net.Mask(), true))

	fmt.Fprintf(w, "Broadcast\t:\t%s\t%s %s\n", net.Broadcast(), bits(net.Broadcast(), net.Mask(), false), bits(net.Broadcast(), net.Mask(), true))
	if net.Version() == 4 {
		// ridiculous for ipv6, should do /64 networks or something
		fmt.Fprintf(w, "Hosts/Net\t:\t%d\t \n", int(net.Size()))
	}
	return nil
}

func main() {
	app := &cli.App{
		Action: handleInput,
	}

	// IPv4 options:
	//  -s MASK     Split the IPv4 network into subnets of MASK size
	// IPv6 options:
	//  -e          IPv4 compatible IPv6 information
	//  -r          IPv6 reverse DNS output
	//  -S MASK     Split the IPv6 network into subnets of MASK size

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
