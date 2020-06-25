package main

// Todo:
// Allow for hot reloading config (maybe watch file?)
// In haws, make of HA's message id system
// I should be able to figure out domain without the user telling me. I /could/ strip off the first part of entity_id, but that sounds risky.
//     I'm pretty sure it's accessible in get_states, so I'll try that.
//
// Better Icon
// Allow for choice of standard or light/dark following system ui theme

import (
	"log"

	"github.com/adriancampos/ha-taskbar/config"
	"github.com/adriancampos/ha-taskbar/haws"
	"github.com/adriancampos/ha-taskbar/icon"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/browser"
)

var c *haws.Connection

type entityItem struct {
	menuItem *systray.MenuItem
	entityID string
	onClick  func()
}

var userEntityItems []entityItem

func main() {
	systray.Run(onReady, nil)
	log.Println("Finished main(); exiting")
}

func onReady() {

	var conf *config.Config

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("HomeAssistant")

	mURL := systray.AddMenuItem("Open Home Assistant", "Open web HA")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mURL.ClickedCh:
				browser.OpenURL(conf.Server.WebAddress)
			}
		}
	}()

	systray.AddSeparator()

	// Attempt to load config
	if _conf, err := config.LoadConfig("configuration.yaml"); err == nil {
		conf = _conf
		loadFromConfig(conf)
	}

	// Listen for menu close events
	go func() {
		for {
			<-systray.MenuClosedCh
			log.Println("systray menu hidden")
			// Close websocket connection when menu closes
			haws.Close(c)
		}
	}()

	// Listen for menu open events
	go func() {
		for {
			<-systray.MenuOpenedCh
			log.Println("systray menu shown")

			// Load config
			if _conf, err := config.LoadConfig("configuration.yaml"); err != nil {
				dlgs.Warning("Invalid Configuration", "configuration.yaml is either not found or invalid:\n"+err.Error())
				log.Println(err)
				continue
			} else {
				if _conf != nil && !cmp.Equal(_conf, conf) {
					log.Println("Config changed!")
					conf = _conf
					loadFromConfig(conf)
				}
			}

			// Close any preexisting connection
			haws.Close(c)

			// Open connection to HA
			c = haws.OpenConnection(conf.Server.Address, conf.Server.AccessToken, func() { haws.RefreshAllEntities(c) }, func() {
				log.Println("Read error!")
			})

			// Subscribe to entity changes
			for _, mItem := range userEntityItems {
				mItem := mItem
				log.Println("subscribing to ", mItem.entityID)
				haws.SubscribeToggleableEntity(c, mItem.entityID, func(te haws.ToggleableEntity) {
					if te.State {
						mItem.menuItem.Check()
					} else {
						mItem.menuItem.Uncheck()
					}

					mItem.menuItem.SetTitle(te.FriendlyName)
				})

			}

		}
	}()
}

func loadFromConfig(conf *config.Config) {
	log.Println("Loading config: ", conf)
	// Remove any existing menu items
	for _, mItem := range userEntityItems {
		// TODO: I'd like to remove instead of hide the menu items
		mItem.menuItem.Disable()
		mItem.menuItem.Hide()
	}

	// Load entities
	haEntities := make([]haws.ToggleableEntity, len(conf.ToggleableEntities))
	for i, te := range conf.ToggleableEntities {
		haEntities[i] = haws.ToggleableEntity{Domain: te.Domain, EntityID: te.EntityID}
	}
	log.Println("Loaded ToggleableEntities:", conf.ToggleableEntities)

	// Define Menu Items
	userEntityItems = make([]entityItem, len(haEntities))
	for i, v := range haEntities {
		v := v
		i := i
		userEntityItems[i] = entityItem{
			menuItem: systray.AddMenuItem(v.EntityID, "Toggle "+v.EntityID),
			entityID: v.EntityID,
			onClick: func() {
				haws.ToggleDevice(c, v)

				// To avoid UI flashes when the menu is next opened, assume that the toggle succeeded.
				// If the assumption is incorrect, it'll be corrected through the haws subscription.
				if userEntityItems[i].menuItem.Checked() {
					userEntityItems[i].menuItem.Uncheck()
				} else {
					userEntityItems[i].menuItem.Check()
				}
			},
		}
	}

	// Set up click listeners
	for _, mItem := range userEntityItems {
		mItem := mItem
		go func(mItem entityItem) {
			for {
				<-mItem.menuItem.ClickedCh
				mItem.onClick()
			}
		}(mItem)
	}
}
