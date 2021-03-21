package remotePinning

import (
	"context"
	"fmt"

	"github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
)

func dnsaddrFormatter(domains []string) ([]string, error) {
	addrList := []string{}
	for _, domain := range domains {
		a, err := multiaddr.NewMultiaddr(fmt.Sprintf("/dnsaddr/%s", domain))

		if err != nil {
			return []string{}, err
		}

		res, err := madns.Resolve(context.TODO(), a)

		if err != nil {
			return []string{}, err
		}

		addrs := []string{}

		for _, addr := range res {
			addrs = append(addrs, addr.String())
		}

		addrList = append(addrList, addrs...)

	}
	fmt.Println(addrList)
	return addrList, nil
}
