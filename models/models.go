// Package models is where structs that are used across the project are declared
package models

import "github.com/jinzhu/gorm"

// PlayerConfig is a struct that represents the player configuration
type PlayerConfig struct {
	gorm.Model
	Video         string `validate:""`
	AudioDecoding string `validate:""`
	URL           string `validate:""`
	Transport     string `validate:""`
	Buffer        int    `validate:""`
	Username      string `validate:""`
	Password      string `validate:""`
	Volume        int    `validate:""`
	Autoplay      string `validate:""`
}

// Status is a struct that represents the status of the system
type Status struct {
	CPU         int
	RAM         int
	Temperature float64
	URL         string
	Video       string
	Status      string
}

// NetworkConfig is a struct that represents the network configuration of the system
type NetworkConfig struct {
	Mode    string
	IP      string
	Gateway string
	Netmask string
	DNS1    string
	DNS2    string
}
