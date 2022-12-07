package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

const configURL = "https://raw.githubusercontent.com/olafhartong/sysmon-modular/master/sysmonconfig.xml"
const sysmonURL = "https://download.sysinternals.com/files/Sysmon.zip"

func main() {
	cmd := exec.Command("sysmon", "-version")
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("Sysmon is not installed. Downloading sysmon...")

		// download sysmon config
		resp, err := http.Get(configURL)
		if err != nil {
			fmt.Println("Error downloading sysmon config: ", err)
			return
		}
		defer resp.Body.Close()

		config, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading sysmon config: ", err)
			return
		}

		// write sysmon config to file
		err = ioutil.WriteFile("sysmonconfig.xml", config, 0644)
		if err != nil {
			fmt.Println("Error writing sysmon config to file: ", err)
			return
		}

		fmt.Println("Sysmon config downloaded and saved to file.")

		// download sysmon
		resp, err = http.Get(sysmonURL)
		if err != nil {
			fmt.Println("Error downloading sysmon: ", err)
			return
		}
		defer resp.Body.Close()

		sysmonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading sysmon data: ", err)
			return
		}

		// write sysmon data to file
		err = ioutil.WriteFile("sysmon.zip", sysmonData, 0644)
		if err != nil {
			fmt.Println("Error writing sysmon data to file: ", err)
			return
		}

		fmt.Println("Sysmon downloaded and saved to file.")

		// Open the zip file
		zipReader, zipErr := zip.OpenReader("sysmon.zip")
		if zipErr != nil {
			fmt.Println("Error opening zip reader: ", zipErr)
			return
		}
		defer zipReader.Close()

		// Extract the files in the zip archive to the current working directory
		for _, file := range zipReader.File {
			fileReader, fileErr := file.Open()
			if fileErr != nil {
				fmt.Println("Error opening file from zip: ", fileErr)
				return
			}
			defer fileReader.Close()

			fileBytes, readErr := ioutil.ReadAll(fileReader)
			if readErr != nil {
				fmt.Println("Error reading file from zip: ", readErr)
				return
			}

			writeErr := ioutil.WriteFile(file.Name, fileBytes, 0644)
			if writeErr != nil {
				fmt.Println("Error writing file: ", writeErr)
				return
			}
		}
		fmt.Println("Successfully extracted files from sysmon.zip")

		// install sysmon
		install := exec.Command(".\\Sysmon.exe", "-accepteula", "-i", ".\\sysmonconfig.xml")
		installOut, installErr := install.Output()

		if installErr != nil {
			fmt.Println("Error installing sysmon: ", installErr)
			return
		}

		fmt.Println(string(installOut))
	} else {
		fmt.Println("Sysmon is already installed.")
	}
}
