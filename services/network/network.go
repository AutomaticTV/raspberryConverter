package network

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os/exec"
	"raspberryConverter/models"
	"regexp"
	"strconv"
	"strings"
)

type networkTemplateData struct {
	IP      string
	Mask    string
	Gateway string
	DNSs    string
}

var staticIPDelimiter = "#STATIC_IP_ENABLED_BY_RASPBERRY_CONVERTER"
var staticTemplate = template.Must(template.New("StaticIP").Parse(
	staticIPDelimiter + `
interface eth0
static ip_address={{.IP}}/{{.Mask}}
static routers={{.Gateway}}
static domain_name_servers={{.DNSs}}
fallback static_eth0
`))

func SetConfig(config models.NetworkConfig) error {
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}
	if isDHCP(configFile) && config.Mode == "DHCP" {
		return nil
	}
	withoutStatic := removeStaticConfig(configFile)
	if config.Mode == "DHCP" {
		return setConfigFile(withoutStatic)
	}
	var buff bytes.Buffer
	if err = staticTemplate.Execute(&buff, networkTemplateData{
		config.IP,
		strconv.Itoa(maskToIP(config.Netmask)),
		config.Gateway,
		config.DNS1 + " " + config.DNS2,
	}); err != nil {
		return err
	}
	err = setConfigFile(withoutStatic + buff.String())
	if err != nil {
		return err
	}
	return resetConfig()
}

func GetConfig() (models.NetworkConfig, error) {
	ipRouteCommand, err := exec.Command("ip", "route").CombinedOutput()
	if err != nil {
		return networkError(err)
	}

	resolvCommand, err := exec.Command("cat", "/etc/resolv.conf").CombinedOutput()
	if err != nil {
		return networkError(err)
	}
	resolvSplited := regexp.MustCompile("[\n ]").Split(string(resolvCommand), -1)

	ifconfigCommand, err := exec.Command("ifconfig", "eth0").CombinedOutput()
	if err != nil {
		return networkError(err)
	}
	ifconfigSplited := regexp.MustCompile("[\n ]").Split(string(ifconfigCommand), -1)

	Mode := "Static"
	configFile, err := getConfigFile()
	if err != nil {
		return networkError(err)
	}
	if isDHCP(configFile) {
		Mode = "DHCP"
	}
	IP := ""
	Gateway := strings.Split(string(ipRouteCommand), " ")[2]
	Netmask := ""
	DNS1 := ""
	DNS2 := ""

	for i, element := range ifconfigSplited {
		if isIP(element) && i > 0 && ifconfigSplited[i-1] == "inet" {
			IP = element
		}
		if isIP(element) && i > 0 && ifconfigSplited[i-1] == "netmask" {
			Netmask = element
		}
	}

	for i, element := range resolvSplited {
		if isIP(element) && i > 0 && resolvSplited[i-1] == "nameserver" {
			if DNS1 == "" {
				DNS1 = element
			} else if DNS2 == "" {
				DNS2 = element
			}
		}
	}
	if DNS2 == "" {
		DNS2 = DNS1
	}

	network := models.NetworkConfig{
		Mode:    Mode,
		IP:      IP,
		Netmask: string(Netmask),
		Gateway: Gateway,
		DNS1:    DNS1,
		DNS2:    DNS2,
	}
	return network, nil
}

func networkError(err error) (models.NetworkConfig, error) {
	fmt.Println(err)
	return models.NetworkConfig{}, errors.New("Error geting network configuration")
}

func isIP(str string) bool {
	isCorrect, err := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$", str)
	if err != nil {
		return false
	}
	return isCorrect
}

func getConfigFile() (string, error) {
	tmp, err := exec.Command("cat", "/etc/dhcpcd.conf").CombinedOutput()
	return string(tmp), err
}

func setConfigFile(content string) error {
	return ioutil.WriteFile("/etc/dhcpcd.conf", []byte(content), 664)
}

func maskToIP(ip string) int {
	ret := 0
	chunks := strings.Split(ip, ".")
	for _, element := range chunks {
		if element == "0" {
			return ret
		} else if element == "255" {
			ret += 8
		} else if element == "254" {
			return ret + 7
		} else if element == "252" {
			return ret + 6
		} else if element == "248" {
			return ret + 5
		} else if element == "240" {
			return ret + 4
		} else if element == "224" {
			return ret + 3
		} else if element == "192" {
			return ret + 2
		} else if element == "128" {
			return ret + 1
		} else {
			return 24
		}
	}
	return ret
}

func isDHCP(config string) bool {
	return !strings.Contains(config, staticIPDelimiter)
}

func removeStaticConfig(config string) string {
	return strings.Split(config, staticIPDelimiter)[0]
}

func resetConfig() error {
	return exec.Command("/bin/sh", "-c", "sudo systemctl daemon-reload && sudo systemctl restart dhcpcd").Run()
}
