package ipaddress

import "fmt"

type IP [4]uint8
type IPNetmask struct {
	Ip      IP
	Netmask uint8
}

func (in *IPNetmask) GetWildcardMask() IP {
	var wildcard IP
	var bits uint32
  suffix := 32 - in.Netmask

	for suffix > 0 {
		bits <<= 1
		bits |= 1
		suffix = suffix - 1
	}

	wildcard[0] = uint8((bits & 0xFF000000) >> (8 * 3))
	wildcard[1] = uint8((bits & 0x00FF0000) >> (8 * 2))
	wildcard[2] = uint8((bits & 0x0000FF00) >> (8 * 1))
	wildcard[3] = uint8((bits & 0x000000FF) >> (8 * 0))

  return wildcard
}

func (in *IPNetmask) GetSubnetMask() IP {
  wildcard := in.GetWildcardMask()

  for i, v := range wildcard {
    wildcard[i] = ^v;
  }

  return wildcard
}

func (in *IPNetmask) GetNetwork() IP {
  netmask := in.GetSubnetMask()
  var network IP

  for i := 0; i < 4; i++ {
    network[i] = netmask[i] & in.Ip[i]
  }

  return network
}

// XOR in Go
// https://stackoverflow.com/a/23025720
func xor(a, b uint8) uint8 {
  return (a | b) & ^(a & b)
}

func (in *IPNetmask) GetBroadcast() IP {
  network := in.GetNetwork()
  wildcard := in.GetWildcardMask()
  var broadcast IP
  
  for i := 0; i < 4; i++ {
    broadcast[i] = xor(network[i], wildcard[i]) 
  }

  return broadcast
}

func (ip *IP) Print() {
  fmt.Printf("%d.%d.%d.%d\n", ip[0], ip[1], ip[2], ip[3])
}

func (in *IPNetmask) Print() {
  fmt.Printf("%d.%d.%d.%d/%d\n", in.Ip[0], in.Ip[1], in.Ip[2], in.Ip[3], in.Netmask)
}
