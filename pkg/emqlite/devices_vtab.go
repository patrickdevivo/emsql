package emqlite

import (
	"encoding/json"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/packethost/packngo"
)

type DevicesModule struct {
	client *packngo.Client
}

func NewDevicesModule(client *packngo.Client) *DevicesModule {
	return &DevicesModule{client}
}

type devicesTable struct {
	module *DevicesModule
}

func (m *DevicesModule) EponymousOnlyModule() {}

func (m *DevicesModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab(fmt.Sprintf(`
		CREATE TABLE %q (
			project_id HIDDEN,
			id TEXT,
			href TEXT,
			hostname TEXT,
			description TEXT,
			state TEXT,
			created DATETIME,
			updated DATETIME,
			locked BOOL,
			billing_cycle TEXT,
			storage TEXT,
			tags TEXT,
			network TEXT,
			volumes TEXT,
			os TEXT,
			plan TEXT,
			facility TEXT,
			project TEXT,
			provision_events TEXT,
			provision_per REAL,
			user_data TEXT,
			user TEXT,
			root_password TEXT,
			ipxe_script_url TEXT,
			always_pxe BOOL,
			hardware_reservation TEXT,
			spot_instance BOOL,
			spot_price_max REAL,
			termination_time DATETIME,
			network_ports TEXT,
			custom_data TEXT,
			ssh_keys TEXT,
			short_id TEXT,
			switch_uuid TEXT
		)`, args[0]))
	if err != nil {
		return nil, err
	}

	return &devicesTable{m}, nil
}

func (m *DevicesModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *DevicesModule) DestroyModule() {}

func (v *devicesTable) Open() (sqlite3.VTabCursor, error) {

	return &devicesCursor{v, "", nil, 0}, nil
}

func (v *devicesTable) Disconnect() error {
	return nil
}
func (v *devicesTable) Destroy() error { return nil }

type devicesCursor struct {
	table     *devicesTable
	projectID string
	devices   []packngo.Device
	index     int
}

func (vc *devicesCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	device := vc.devices[vc.index]
	switch col {
	case 0:
		c.ResultText(vc.projectID)
	case 1:
		c.ResultText(device.ID)
	case 2:
		c.ResultText(device.Href)
	case 3:
		c.ResultText(device.Hostname)
	case 4:
		if device.Description == nil {
			c.ResultNull()
		} else {
			c.ResultText(*device.Description)
		}
	case 5:
		c.ResultText(device.State)
	case 6:
		c.ResultText(device.Created)
	case 7:
		c.ResultText(device.Updated)
	case 8:
		c.ResultBool(device.Locked)
	case 9:
		c.ResultText(device.BillingCycle)
	case 10:
		s, err := json.Marshal(device.Storage)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 11:
		s, err := json.Marshal(device.Tags)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 12:
		s, err := json.Marshal(device.Network)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 13:
		s, err := json.Marshal(device.Volumes)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 14:
		s, err := json.Marshal(device.OS)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 15:
		s, err := json.Marshal(device.Plan)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 16:
		s, err := json.Marshal(device.Facility)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 17:
		s, err := json.Marshal(device.Project)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 18:
		s, err := json.Marshal(device.ProvisionEvents)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 19:
		c.ResultDouble(float64(device.ProvisionPer))
	case 20:
		c.ResultText(device.UserData)
	case 21:
		c.ResultText(device.User)
	case 22:
		c.ResultText(device.RootPassword)
	case 23:
		c.ResultText(device.IPXEScriptURL)
	case 24:
		c.ResultBool(device.AlwaysPXE)
	case 25:
		s, err := json.Marshal(device.HardwareReservation)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 26:
		c.ResultBool(device.SpotInstance)
	case 27:
		c.ResultDouble(device.SpotPriceMax)
	case 28:
		if device.TerminationTime == nil {
			c.ResultNull()
		} else {
			c.ResultText(device.TerminationTime.String())
		}
	case 29:
		s, err := json.Marshal(device.NetworkPorts)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 30:
		s, err := json.Marshal(device.CustomData)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 31:
		s, err := json.Marshal(device.SSHKeys)
		if err != nil {
			return err
		}
		c.ResultText(string(s))
	case 32:
		c.ResultText(device.ShortID)
	case 33:
		c.ResultText(device.SwitchUUID)
	}
	return nil
}

func (v *devicesTable) BestIndex(constraints []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	used := make([]bool, len(constraints))
	cost := 1000.0
	projectIDCstUsed := false
	for c, cst := range constraints {
		if !cst.Usable || cst.Op != sqlite3.OpEQ {
			continue
		}
		switch cst.Column {
		case 0: // project_id
			used[c] = true
			projectIDCstUsed = true
		}
	}

	if projectIDCstUsed {
		// if the project ID constraint is used, cost is 0 to force sqlite to always use this index
		cost = 0
	}

	return &sqlite3.IndexResult{
		IdxNum:        0,
		IdxStr:        "default",
		Used:          used,
		EstimatedCost: cost,
	}, nil
}

func (vc *devicesCursor) Filter(idxNum int, idxStr string, vals []interface{}) error {
	projectID := vals[0].(string)
	vc.projectID = projectID

	devices, _, err := vc.table.module.client.Devices.List(vc.projectID, &packngo.ListOptions{})
	if err != nil {
		return err
	}

	vc.devices = devices
	vc.index = 0

	return nil
}

func (vc *devicesCursor) Next() error {
	vc.index++
	return nil
}

func (vc *devicesCursor) EOF() bool {
	return vc.index >= len(vc.devices)
}

func (vc *devicesCursor) Rowid() (int64, error) {
	return int64(vc.index), nil
}

func (vc *devicesCursor) Close() error {
	return nil
}
