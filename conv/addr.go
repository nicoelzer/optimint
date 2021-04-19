package conv

import (
	"strings"

	"github.com/lazyledger/optimint/config"
	"github.com/multiformats/go-multiaddr"
)

func TranslateAddresses(conf *config.NodeConfig) error {
	if conf.P2P.ListenAddress != "" {
		addr, err := GetMultiAddr(conf.P2P.ListenAddress)
		if err != nil {
			return err
		}
		conf.P2P.ListenAddress = addr.String()
	}

	seeds := strings.Split(conf.P2P.Seeds, ",")
	for i, seed := range seeds {
		if seed != "" {
			addr, err := GetMultiAddr(seed)
			if err != nil {
				return err
			}
			seeds[i] = addr.String()
		}
	}
	conf.P2P.Seeds = strings.Join(seeds, ",")

	return nil
}

func GetMultiAddr(addr string) (multiaddr.Multiaddr, error) {
	var err error
	var p2pId multiaddr.Multiaddr
	if at := strings.IndexRune(addr, '@'); at != -1 {
		p2pId, err = multiaddr.NewMultiaddr("/p2p/" + addr[:at])
		if err != nil {
			return nil, err
		}
		addr = addr[at+1:]
	}
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return nil, ErrInvalidAddress
	}
	maddr, err := multiaddr.NewMultiaddr("/ip4/" + parts[0] + "/tcp/" + parts[1])
	if err != nil {
		return nil, err
	}
	if p2pId != nil {
		maddr = maddr.Encapsulate(p2pId)
	}
	return maddr, nil
}