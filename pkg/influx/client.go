package influx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const dbName = "test"

func InfluxClientV2() (*influxdb2.Client, error) {
	// Generate a Token from the "Tokens Tab" in the UI
	token := os.Getenv("INFLUXDB_TOKEN")
	hostUrl := os.Getenv("INFLUXDB_URL")

	if token == "" || hostUrl == "" {
		return nil, errors.New("InfluxDB couldn't be located.")
	}

	client := influxdb2.NewClient(hostUrl, token)
	return &client, nil
}

func InfluxClientV3() (*influxdb3.Client, func(*influxdb3.Client)) {
	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")
	fmt.Printf("INFLUXDB_URL: %s\n", url)
	fmt.Printf("INFLUXDB_TOKEN: %s\n", token)

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         url,
		Token:        token,
		Database:     dbName,
		WriteOptions: &influxdb3.WriteOptions{Database: dbName},
	})

	if err != nil {
		panic(err)
	}
	callback := func(client *influxdb3.Client) {
		fmt.Println("closing influx client")
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}

	return client, callback
}

func LogRuntime(processName string, elapsed time.Duration) {
	client, closeClient := InfluxClientV3()
	defer closeClient(client)
	point := influxdb3.NewPointWithMeasurement("system").
			SetTag("process_name", processName).
			SetField("benchmark[ms]", elapsed.Milliseconds()).
			SetField("benchmark[ns]", elapsed.Nanoseconds())

	if err := client.WritePoints(context.Background(), []*influxdb3.Point{point}); err != nil {
		fmt.Println("error: ", err)
	}
}

func LogSystemMetrics(data map[string]map[string]interface{}) error {
	client, closeClient := InfluxClientV3()
	defer closeClient(client)

	options := influxdb3.WriteOptions{
		Database: dbName,
	}

	points := make([]*influxdb3.Point, 0)
	
	for key := range data {
		point := influxdb3.NewPointWithMeasurement("system").
			SetTag(key, data[key]["label"].(string))
		for key, val := range data[key] {
			point.SetField(key, val)
			points = append(points, point)
		}

		if err := client.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			fmt.Println("error: ", err)
			return err
		}
	}

	return nil
}

func TestInfluxV3() {
	client, closeClient := InfluxClientV3()
	defer closeClient(client)

	data := map[string]map[string]interface{}{
		"point1": {
			"location": "London",
			"species":  "bees",
			"count":    23,
		},
		"point2": {
			"location": "Portland",
			"species":  "ants",
			"count":    30,
		},
	}

	options := influxdb3.WriteOptions{
		Database: dbName,
	}

	for key := range data {
		point := influxdb3.NewPointWithMeasurement("census").
			SetTag("location", data[key]["location"].(string)).
			SetField(data[key]["species"].(string), data[key]["count"])

		if err := client.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second) // separate points by 1 second
	}

	query := `SELECT * FROM 'census'`

	ctx := context.Background()
	// dbOpt := influxdb3.QueryOption(influxdb3.WithDatabase(database))
	iterator, err := client.Query(ctx, query)
	// iterator, err := client.QueryWithOptions(context.Background(), &queryOptions, query)

	if err != nil {
		if err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
			panic(errors.New("Did you set your .env vars? - " + err.Error()))
		}
		panic(err)
	}

	for iterator.Next() {
		value := iterator.Value()

		location := value["location"]
		ants := value["ants"]
		if ants == nil {
			ants = 0
		}
		bees := value["bees"]
		if bees == nil {
			bees = 0
		}
		fmt.Printf("in %s are %d ants and %d bees\n", location, ants, bees)
	}
}
