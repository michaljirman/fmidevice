package fmidevice

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaljirman/fmidevice/applews/icloud/fmip"
	"github.com/michaljirman/fmidevice/applews/icloud/setup"
	"github.com/michaljirman/fmidevice/applews/idmsa"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LocateCmd represents locate command
var LocateCmd = &cobra.Command{
	Use:   "locate",
	Short: "Locate iDevices for an Apple account",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Locate(args); err != nil {
			log.Fatalln("Locate Failed:", err)
		}
	},
}

const (
	clientBuildNumber     = "1809Project50"
	clientMasteringNumber = "1809B29"
)

var clientID string
var accountName string
var password string
var silent bool
var authToken string
var prsID int

// PrintDeviceInfo prints device information into the console
func PrintDeviceInfo(fmipService *fmip.FmipService) {
	for i, device := range fmipService.FmipClientResponse.Content {
		fmt.Println("----------")
		fmt.Printf("[%d] Device: \033[92m%s, %s ... Location: %f, %f\033[0m\n",
			i+1, device.DeviceDisplayName, device.Name,
			device.Location.Latitude, device.Location.Longitude)
	}
}

// Locate handles comunication with location services and printing the results to the console
func Locate(args []string) error {
	if silent {
		if accountName == "" {
			accountName = viper.GetString("accountName")
			if accountName == "" {
				return fmt.Errorf("Apple account name must be provided as a flag or defined in configuration file for silent mode")
			}
		}
		if password == "" {
			authToken = viper.GetString("AuthToken")
			prsID = viper.GetInt("prsID")
			if authToken == "" || !silent {
				return fmt.Errorf("Apple password must be provided as a flag or an AuthToken must be defined in configuration file for silent mode")
			}
		}
	} else {
		if accountName == "" {
			return fmt.Errorf("Apple account name must be provided as a flag for default mode")
		}
		if password == "" {
			return fmt.Errorf("Apple password must be provided as a flag for default mode")
		}
	}

	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	locateErrChan := make(chan error)
	go func() {
		if silent {
			//Location aquired via mobile client
			fmipService, err := fmip.NewFmipMobileService("https://fmipmobile.icloud.com")
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to create a new fmip service: %s", err)
			}
			if authToken == "" || prsID <= 0 {
				fmipService.InitMobileClient(accountName, password)
				prsID = fmipService.FmipClientResponse.ServerContext.PrsID
				authToken = fmipService.FmipClientResponse.ServerContext.AuthToken
				fmt.Printf("Successfully authenticated:\n---\nAccountName: %s\nPrsId: %d\nAuthToken: %s\n",
					accountName, prsID, authToken)
				PrintDeviceInfo(fmipService)
			}

			for {
				fmipService.RefreshMobileClient(accountName, prsID, authToken)
				PrintDeviceInfo(fmipService)
				time.Sleep(5 * time.Second)
			}
		} else {
			//Generates a new clientID (uuid)
			clientIDUUID, err := uuid.NewV4()
			if err != nil {
				locateErrChan <- err
			}
			clientID = clientIDUUID.String()
			//Location aquired via web client
			setupService, err := setup.NewSetupService("https://setup.icloud.com", clientBuildNumber, clientID, clientMasteringNumber)
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to create a new setup service: %s", err)
			}
			setupService.Validate()
			accountLoginUIBaseURL, accountLoginUIResource, err := setupService.GetAccountLoginUIHostAndResource()
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to retrieve accountLoginUI URL")
			}
			_, accountLoginResource, err := setupService.GetAccountLoginHostAndResource()
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to retrieve accountLogin URL")
			}

			idmsaService, err := idmsa.NewIdmsaService(accountLoginUIBaseURL)
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to create a new idmsa service: %s", err)
			}
			idmsaService.AccountLoginUI(accountLoginUIResource, accountName, password)

			setupService.AccountLogin(accountLoginResource, idmsaService.XappleSessionToken, idmsaService.XappleIDAccountCountry)

			dsid := setupService.AccountLoginResponse.DsInfo.Dsid
			findmeURL := setupService.AccountLoginResponse.Webservices.Findme.URL
			xAppleCookiesHeader := setupService.GetXappleCookiesHeader()
			fmipService, err := fmip.NewFmipWebService(findmeURL)
			if err != nil {
				locateErrChan <- fmt.Errorf("Unable to create a new fmip service: %s", err)
			}
			fmipService.InitWebClient(xAppleCookiesHeader, clientBuildNumber, clientID, clientMasteringNumber, dsid)
			PrintDeviceInfo(fmipService)
			for {
				// fmt.Println("\033[92mRefreshing client ...\033[0m")
				time.Sleep(5 * time.Second)
				PrintDeviceInfo(fmipService)
			}
		}
	}()

loop:
	for {
		select {
		case <-sigs:
			fmt.Println(" Got shutdown event, exiting gracefully ...")
			// Break out of the outer for statement and end the program
			break loop
		case err := <-locateErrChan:
			return err
		}
	}

	return nil
}

func init() {
	LocateCmd.PersistentFlags().StringVar(&accountName, "account", "", "An Apple account name")
	LocateCmd.PersistentFlags().StringVar(&password, "password", "", "An Apple password")
	LocateCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Silent mode allows to access location through a mobile client")
}
