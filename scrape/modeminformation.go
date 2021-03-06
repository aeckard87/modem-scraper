package scrape

import (
	"encoding/json"

	_ "github.com/influxdata/influxdb1-client" // this is important because of a bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

// ModemInformation holds all information from the
// SB8200 status pages.
type ModemInformation struct {
	ConnectionStatus    ConnectionStatus
	SoftwareInformation SoftwareInformation
	EventLog            []EventLog
}

// ToJSON converts ModemInformation to JSON string.
func (m ModemInformation) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	return jsonString, nil
}

// ToInfluxPoints converts ModemInformation to "points"
// for use with InfluxDB.
// TODO: Really, this should be broken up and implemented in:
// - downstreambondedchannels ([]string)
//   - downstreambondedchannel,channel_id=# channel_id=#,frequency=#,etc
// - upstreambondedchannels ([]string)
// - softwareinformation
// - startupprocedure
//
// or maybe... do the above, and also have this (ModemInformation)
// call each of the above to aggregate them all. Then this can
// simply return []string with all the "lines" and they can
// all be sent to InfluxDB at once.
func (m ModemInformation) ToInfluxPoints() ([]*client.Point, error) {
	var points []*client.Point

	influxPoints, err := m.ConnectionStatus.ToInfluxPoints()
	if err != nil {
		return nil, err
	}
	points = append(points, influxPoints...)

	influxPoints, err = m.SoftwareInformation.ToInfluxPoints()
	if err != nil {
		return nil, err
	}
	points = append(points, influxPoints...)

	influxPoints, err = buildEventLogPoints(m.EventLog)
	if err != nil {
		return nil, err
	}
	points = append(points, influxPoints...)

	return points, nil
}
