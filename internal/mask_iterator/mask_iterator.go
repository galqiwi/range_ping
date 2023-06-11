package mask_iterator

import (
	"net"
)

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IPGenerator(mask string) (chan net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(mask)
	if err != nil {
		return nil, err
	}

	output := make(chan net.IP)

	go func() {
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			outputIP := make(net.IP, 4)
			copy(outputIP, ip)
			output <- outputIP
		}
		close(output)
	}()

	return output, nil
}
