package main

import "os"

func main() {
	if os.Getenv("YIFI_EXAMPLE") == "public" {
		PublicReadOnlyCycle()
		return
	}
	if os.Getenv("YIFI_EXAMPLE") == "regular" {
		FullCycle()
		return
	}
	if os.Getenv("YIFI_EXAMPLE") == "approve" {
		Approve()
		return
	}
	PartnerFullCycle()
}
