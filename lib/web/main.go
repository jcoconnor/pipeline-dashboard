package web

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"encoding/csv"
	"encoding/json"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/constants"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/jenkins_types"
	"github.com/puppetlabs/pipeline-dashboard/lib/report/utils"
)

type Page struct {
    Title       string
    LastUpdated string
    Jobs        []jenkins_types.Pipeline
    Trains      []jenkins_types.TrainStrings
    Links       []jenkins_types.Link
    Products    []jenkins_types.Product
}

func (h *Handlers) GeneratePageData() *Page {
    title := "CI Dashboard"

    fmt.Println("Generating Page Data.")

    // Lock for accessing files, but defer to ensure we release it on return.
    utils.Lock()
    defer utils.Unlock()

    // Does file exist - otherwise return with no page data
    statinfo, err := os.Stat(constants.Results_filename)
    if os.IsNotExist(err) {
        h.Page = nil
        return h.Page
    }
    lastupdated := statinfo.ModTime().Format(http.TimeFormat)

    csvFile, err := os.Open(constants.Results_filename)
    if (err != nil) {
        fmt.Println(err.Error())
        h.Page = nil
        return h.Page
    }
    defer csvFile.Close()
    reader := csv.NewReader(bufio.NewReader(csvFile))

    var jobs []jenkins_types.Pipeline

    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }

        buildNumber, _ := strconv.Atoi(line[5])

        errors, _ := strconv.Atoi(line[12])
        transients, _ := strconv.Atoi(line[13])

        jobs = append(jobs, jenkins_types.Pipeline{
            URL:         line[0],
            Server:      line[1],
            Pipeline:    line[3],
            PipelineJob: line[2],
            Version:     line[4],
            BuildNumber: buildNumber,
            JobDataStrings: &jenkins_types.JobDataStrings{
                StartTime:            line[6],
                EndTime:              line[7],
                WallClockTimeHours:   line[8],
                WallClockTimeMinutes: line[9],
                TotalHours:           line[10],
                TotalMinutes:         line[11],
                QueueTimeHours:       line[12],
                QueueTimeMinutes:     line[13],
            },
            Errors:     errors,
            Transients: transients,
        })
    }

    trainCSVFile, _ := os.Open(constants.Trains_filename)
    if (err != nil) {
        fmt.Println(err.Error())
        return h.Page
    }
    defer trainCSVFile.Close()
    trainReader := csv.NewReader(bufio.NewReader(trainCSVFile))

    var trains []jenkins_types.TrainStrings

    for {
        line, error := trainReader.Read()

        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }

        minutes, _ := strconv.ParseFloat(line[4], 64)
        hours := minutes / 60
        minutesLeft := int(minutes) % 60

        queueMinutes, _ := strconv.ParseFloat(line[5], 64)
        queueHours := queueMinutes / 60
        queueMinutesLeft := int(queueMinutes) % 60

        errors, _ := strconv.Atoi(line[9])
        transients, _ := strconv.Atoi(line[10])

        trains = append(trains, jenkins_types.TrainStrings{
            Pipeline:             line[0],
            Version:              line[1],
            URL:                  line[2],
            Name:                 line[3],
            DurationSortMinutes:  int(minutes),
            DurationMinutes:      fmt.Sprintf("%d", int(minutesLeft)),
            DurationHours:        fmt.Sprintf("%d", int(hours)),
            QueueTimeSortMinutes: int(queueMinutes),
            QueueTimeMinutes:     fmt.Sprintf("%d", int(queueMinutesLeft)),
            QueueTimeHours:       fmt.Sprintf("%d", int(queueHours)),
            StartTime:            line[6],
            EndTime:              line[7],
            Timestamp:            line[8],
            Errors:               errors,
            Transients:           transients,
            Platform:             line[11],
            PlatformVersion:      line[12],
        })
    }

    sort.Slice(trains, func(i, j int) bool {
        return trains[i].DurationSortMinutes > trains[j].DurationSortMinutes
    })

    for i, product := range h.Products {
        product.SetVals(jobs)
        h.Products[i] = product
    }

    h.Page = &Page{Title: title, LastUpdated: lastupdated, Jobs: jobs, Trains: trains, Products: h.Products, Links: h.Links}

    return h.Page
}

func (handlers *Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, handlers.Page)
}

type Handlers struct {
    Products []jenkins_types.Product
    Links    []jenkins_types.Link
    Page     *Page
}

func (handlers *Handlers) ProductsHandler(w http.ResponseWriter, r *http.Request) {
    handlers.GeneratePageData()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(handlers.Page)
}

func (handlers *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)

    w.Write([]byte("OK"))
}

func Serve() {
    fs := http.FileServer(http.Dir("./public/"))

    handlers := &Handlers{
        Products: jenkins_types.GetProducts(),
        Links:    jenkins_types.GetLinks(),
    }

    http.Handle("/", http.FileServer(http.Dir("./public/")))

    http.Handle("/css/", http.FileServer(http.Dir("./public/")))

    http.Handle("/static/css/", fs)
    http.Handle("/static/js/", fs)

    http.HandleFunc("/api/1/products", handlers.ProductsHandler)
    http.Handle("/metrics", handlers.GenerateMetrics(promhttp.Handler()))
    http.HandleFunc("/api/1/healthcheck", handlers.HealthCheck)

    fmt.Println("Listening on port :8080")

    log.Fatal(http.ListenAndServe(":8080", nil))
}
