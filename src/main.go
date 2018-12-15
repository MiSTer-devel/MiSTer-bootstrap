package main

import "C"

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Repo is the default repo structure.
type Repo struct {
	File string `json:"file"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	Update()
}

var metadata = map[string]string{
	"Altair8800_Mister":           "Computer",
	"Amstrad_MiSTer":              "Computer",
	"Apogee_MiSTer":               "Computer",
	"Apple-II_MiSTer":             "Computer",
	"Aquarius_MISTer":             "Computer",
	"Arcade-Alibaba_MiSTer":       "Arcade",
	"Arcade-Amidar_MiSTer":        "Arcade",
	"Arcade-Azurian_MiSTer":       "Arcade",
	"Arcade-Bagman_MiSTer":        "Arcade",
	"Arcade-BlackHole_MiSTer":     "Arcade",
	"Arcade-BombJack_MiSTer":      "Arcade",
	"Arcade-BurgerTime_MiSTer":    "Arcade",
	"Arcade-BurningRubber_MiSTer": "Arcade",
	"Arcade-Catacomb_MiSTer":      "Arcade",
	"Arcade-ComputerSpace_MiSTer": "Arcade",
	"Arcade-CosmicAvenger_MiSTer": "Arcade",
	"Arcade-CrazyClimber_MiSTer":  "Arcade",
	"Arcade-CrazyKong_MiSTer":     "Arcade",
	"Arcade-CrushRoller_MiSTer":   "Arcade",
	"Arcade-Defender_MiSTer":      "Arcade",
	"Arcade-DonkeyKong_MiSTer":    "Arcade",
	"Arcade-Dorodon_MiSTer":       "Arcade",
	"Arcade-DreamShopper_MiSTer":  "Arcade",
	"Arcade-Eeekk_MiSTer":         "Arcade",
	"Arcade-Eyes_MiSTer":          "Arcade",
	"Arcade-Frogger_MiSTer":       "Arcade",
	"Arcade-Galaga_MiSTer":        "Arcade",
	"Arcade-Galaxian_MiSTer":      "Arcade",
	"Arcade-Gorkans_MiSTer":       "Arcade",
	"Arcade-LadyBug_MiSTer":       "Arcade",
	"Arcade-LizardWizard_MiSTer":  "Arcade",
	"Arcade-MoonCresta_MiSTer":    "Arcade",
	"Arcade-MoonPatrol_MiSTer":    "Arcade",
	"Arcade-MrDoNightmare_MiSTer": "Arcade",
	"Arcade-MrTNT_MiSTer":         "Arcade",
	"Arcade-MsPacman_MiSTer":      "Arcade",
	"Arcade-Omega_MiSTer":         "Arcade",
	"Arcade-Orbitron_MiSTer":      "Arcade",
	"Arcade-PacmanClub_MiSTer":    "Arcade",
	"Arcade-PacmanPlus_MiSTer":    "Arcade",
	"Arcade-Pacman_MiSTer":        "Arcade",
	"Arcade-PacmanicMiner_MiSTer": "Arcade",
	"Arcade-Pengo_MiSTer":         "Arcade",
	"Arcade-Phoenix_MiSTer":       "Arcade",
	"Arcade-Pisces_MiSTer":        "Arcade",
	"Arcade-Ponpoko_MiSTer":       "Arcade",
	"Arcade-Pooyan_MiSTer":        "Arcade",
	"Arcade-Scramble_MiSTer":      "Arcade",
	"Arcade-SnapJack_MiSTer":      "Arcade",
	"Arcade-SuperGlob_MiSTer":     "Arcade",
	"Arcade-TheEnd_MiSTer":        "Arcade",
	"Arcade-TimePilot_MiSTer":     "Arcade",
	"Arcade-VanVanCar_MiSTer":     "Arcade",
	"Arcade-WarOfTheBugs_MiSTer":  "Arcade",
	"Arcade-Woodpecker_MiSTer":    "Arcade",
	"Arcade-Xevious_MiSTer":       "Arcade",
	"Archie_MiSTer":               "Computer",
	"Atari2600_MiSTer":            "Console",
	"Atari800_MiSTer":             "Computer",
	"BBCMicro_MiSTer":             "Computer",
	"BK0011M_MiSTer":              "Computer",
	"C16_MiSTer":                  "Computer",
	"C64_MiSTer":                  "Computer",
	"ColecoVision_MiSTer":         "Console",
	"Gameboy_MiSTer":              "Console",
	"Genesis_MiSTer":              "Console",
	"Jupiter_MiSTer":              "Computer",
	"Main_MiSTer":				   "Main",
	"MSX_MiSTer":                  "Computer",
	"MacPlus_MiSTer":              "Computer",
	"MemTest_MiSTer":              "Utility",
	"Menu_MiSTer":                 "Main",
	"Minimig-AGA_MiSTer":          "Computer",
	"MultiComp_MiSTer":            "Computer",
	"NES_MiSTer":                  "Console",
	"PET2001_MiSTer":              "Computer",
	"QL_MiSTer":                   "Computer",
	"SAM-Coupe_MiSTer":            "Computer",
	"SMS_MiSTer":                  "Console",
	"SNES_MiSTer":                 "Console",
	"SharpMZ_MiSTer":              "Computer",
	"Specialist_MiSTer":           "Computer",
	"TI-99_4A_MiSTer":             "Computer",
	"TSConf_MiSTer":               "Computer",
	"TurboGrafx16_MiSTer":         "Console",
	"VIC20_MiSTer":                "Computer",
	"Vector-06C_MiSTer":           "Computer",
	"Vectrex_MiSTer":              "Console",
	"X68000_MiSTer":               "Computer",
	"ZX-Spectrum_MISTer":          "Computer",
	"ZX81_MiSTer":                 "Computer",
	"ao486_MiSTer":                "Computer",
}

//export Update
func Update() {
	repo := flag.String("r", "https://raw.githubusercontent.com/OpenVGS/MiSTer-repository/master/repo.json", "Repo URL")
	output := flag.String("o", ".", "Output Directory")

	flag.Parse()

	absOutput, _ := filepath.Abs(*output)

	err := os.MkdirAll(absOutput, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating output directory:", err)
		return
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", *repo, nil)
	if err != nil {
		log.Fatal("Building repo request:", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Repo request:", err)
		return
	}

	defer Close(resp.Body)

	var repos []Repo

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Println(err)
	}

	regex, err := regexp.Compile("\\d{8}")
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		fmt.Println(repo.Name)
	}

	for _, repo := range repos {
		coreType, found := metadata[repo.Name]
		if found == false {
			log.Println("Unrecognized repo", repo.Name)

			if strings.Contains(repo.File, "Arcade") {
				coreType = "Arcade"
			} else {
				coreType = "New"
			}
		}

		outputDir := ""

		if coreType != "Main" {
			outputDir = fmt.Sprintf("%s/_%s", absOutput, coreType)

			err := os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				log.Fatal("Error creating output directory:", err)
				return
			}
		} else {
			outputDir = absOutput
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", outputDir, repo.File)); os.IsNotExist(err) {
			log.Printf("Downloading updated %s core...\n", repo.Name)

			downloadLocation := ""

			if coreType == "Main" {
				if repo.Name == "Menu_MiSTer" {
					downloadLocation = fmt.Sprintf("%s/%s", outputDir, "menu.rbf")
				} else if repo.Name == "Main_MiSTer" {
					downloadLocation = fmt.Sprintf("%s/%s", outputDir, "MiSTer")
				}
			} else {
				downloadLocation = fmt.Sprintf("%s/%s", outputDir, repo.File)
			}

			err := DownloadCore(downloadLocation, repo.URL)

			if err != nil {
				panic(err)
			}

			globString := regex.ReplaceAllString(repo.File, "*")
			files, _ := filepath.Glob(fmt.Sprintf("%s/%s", outputDir, globString))
			currentFile := fmt.Sprintf("%s/%s", outputDir, repo.File)
			fmt.Println("Current file", currentFile)
			for _, file := range files {
				if file != currentFile {
					err := os.Remove(file)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("Deleted:", file)
				}
			}
		} else {
			log.Printf("Core %s already latest version, skipping\n", repo.Name)
		}
	}
}

// DownloadCore downloads core to local filesystem.
func DownloadCore(path string, url string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	out, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer Close(out)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer Close(resp.Body)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Close is a generic io Closer with error handling
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal("IO close error:", err)
	}
}
