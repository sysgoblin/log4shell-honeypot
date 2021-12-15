package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type ElasticBody struct {
	Status      int                `json:"status"`
	Name        string             `json:"name"`
	ClusterName string             `json:"cluster_name"`
	ClusterUUID string             `json:"cluster_uuid"`
	Version     ElasticBodyVersion `json:"version"`
	Tagline     string             `json:"tagline"`
}

type ElasticBodyVersion struct {
	Number        string `json:"number"`
	BuildHash     string `json:"build_hash"`
	BuildDate     string `json:"build_date"`
	BuildSnapshot bool   `json:"build_snapshot"`
	LuceneVersion string `json:"lucene_version"`
}

var (
	elasticHeaders = map[string]string{
		"Content-Type":              "application/json; charset=UTF-8",
		"X-Cloud-Request-Id":        "adJdTYZdQ6GLK6UGVpB9Tx",
		"X-Elastic-Product":         "Elasticsearch",
		"X-Found-Handling-Cluster":  "607a036b350db1d65291d2520ec0a0d22630eb5c",
		"X-Found-Handling-Instance": "instance-0000000420",
	}
	elasticBody = ElasticBody{
		Status:      200,
		Name:        "instance-0000000420",
		ClusterName: "607a036b350db1d65291d2520ec0a0d22630eb5c",
		ClusterUUID: "3wJaQxyzRDWhF50MFpwLXQ",
		Version: ElasticBodyVersion{
			Number:        "7.14.0",
			BuildHash:     "41010b6aef45716303ffbe6f9591503281d79f62",
			BuildDate:     "2021-07-29T20:49:32.864135063Z",
			BuildSnapshot: false,
			LuceneVersion: "8.0.0",
		},
		Tagline: "You Know, for Search",
	}
)

func CreateElasticResponse(w http.ResponseWriter) {
	// set header for content type json
	for key, value := range elasticHeaders {
		w.Header().Set(key, value)
	}

	// pretty print elasticBody json
	prettyElasticBody, err := json.MarshalIndent(elasticBody, "", "  ")
	if err != nil {
		log.Fatal("Unable to marshal elasticBody json: ", err)
	}

	// add newline to end of prettyElasticBody
	prettyElasticBody = append(prettyElasticBody, '\n')
	w.Write(prettyElasticBody)
	return
}
