package bigippg

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
)

type Config struct {
	Address        string
	Port           string
	Username       string
	Password       string
	LoginReference string `json:"loginProviderName"`
	ConfigOptions  *bigip.ConfigOptions
}

func (c *Config) Client() (*bigip.BigIP, error) {

	if c.Address != "" && c.Username != "" && c.Password != "" {
		log.Println("[INFO] Initializing BigIP connection")
		var client *bigip.BigIP
		var err error
		if c.LoginReference != "" {
			client, err = bigip.NewTokenSession(c.Address, c.Port, c.Username, c.Password, c.LoginReference, c.ConfigOptions)
			if err != nil {
				log.Printf("[ERROR] Error creating New Token Session %s ", err)
				return nil, err
			}

		} else {
			client = bigip.NewSession(c.Address, c.Port, c.Username, c.Password, c.ConfigOptions)
		}
		err = c.validateConnection(client)
		if err == nil {
			return client, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("BigIP provider requires address, username and password")
}

func (c *Config) validateConnection(client *bigip.BigIP) error {
	t, err := client.SelfIPs()
	if err != nil {
		log.Printf("[ERROR] Connection to BigIP device could not have been validated: %v ", err)
		return err
	}

	if t == nil {
		log.Printf("[WARN] Could not validate connection to BigIP")
		return nil
	}
	return nil
}
