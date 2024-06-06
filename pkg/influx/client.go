package influx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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

func aggregateMetrics(data map[string]interface{}) map[string]interface{} {
	averagedMetrics := map[string]bool{
		"process.memory": true,
		"process.cpu":    true,
	}
	out := make(map[string]interface{})
	out["process.memory.num_processes"] = 0
	out["process.cpu.num_processes"] = 0
	out["process.memory"] = uint64(0)
	out["process.cpu"] = float64(0)

	for key, val := range data {
		split := strings.Split(key, ".")
		if len(split) == 3 {
			averagedKey := split[0] + "." + split[2]
			if averagedMetrics[averagedKey] {
				switch averagedKey {
				case "process.memory":
					out["process.memory.num_processes"] = out["process.memory.num_processes"].(int) + 1
					out["process.memory.total"] = out["process.memory"].(uint64) + val.(uint64)
				case "process.cpu":
					out["process.cpu.num_processes"] = out["process.cpu.num_processes"].(int) + 1
					out["process.cpu.total"] = out["process.cpu"].(float64) + val.(float64)
				default:
					yellow := "\033[33m"
					cyan := "\033[36m"
					resetColor := "\033[0m"
					fmt.Println(
						yellow, "Warning: Unexpected system metric found: \n-> ",
						cyan, averagedKey, resetColor,
					)
					out[key] = val
				}
			}
		} else {
			out[key] = val
		}
	}

	// Check if "process.memory.num_processes" exists and is of type int
	if numProcesses, ok := out["process.memory.num_processes"].(int); ok && numProcesses > 0 {
		if cpuUsage, ok := out["process.cpu.total"].(float64); ok {
			out["process.memory.average"] = cpuUsage / float64(numProcesses)
		} else {
			fmt.Println("Warning: 'process.cpu.total' is not of type float64")
		}
	} else {
		fmt.Println("Warning: 'process.memory.num_processes' is missing or not of type int")
	}

	// Check if "process.cpu.num_processes" exists and is of type int
	if numProcesses, ok := out["process.cpu.num_processes"].(int); ok && numProcesses > 0 {
		if cpuUsage, ok := out["process.cpu.total"].(float64); ok {
			out["process.cpu.average"] = cpuUsage / float64(numProcesses)
		} else {
			fmt.Println("Warning: 'process.cpu.total' is not of type float64")
		}
	} else {
		fmt.Println("Warning: 'process.cpu.num_processes' is missing or not of type int")
	}

	return out
}

func LogRuntime(processName string, elapsed time.Duration) {
	client, closeClient := InfluxClientV3()
	defer closeClient(client)
	point := influxdb3.NewPointWithMeasurement("benchmarks").
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

	for timestamp, mapp := range data {
		mapp2 := aggregateMetrics(mapp)
		point := influxdb3.NewPointWithMeasurement("system_metrics")

		// Parse the timestamp string into a time.Time object
		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			fmt.Println("error parsing timestamp:", err)
			continue
		}

		// Set the timestamp as the point's timestamp
		point.SetTimestamp(t)

		for key, val := range mapp2 {
			point.SetField(key, val)
		}

		if err := client.WritePointsWithOptions(context.Background(), &options, point); err != nil {
			fmt.Println("error writing point:", err)
			return err
		}
	}

	green := "\033[32m"
	resetColor := "\033[0m"
	fmt.Println(green + "wrote system metrics to InfluxDB" + resetColor)
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
