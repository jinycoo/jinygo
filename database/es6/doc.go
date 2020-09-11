/*
package es6 provides a Go client for Elasticsearch.

Create the client with the NewDefaultClient function:

		es6.NewDefaultClient()

The ELASTICSEARCH_URL environment variable is used instead of the default URL, when set.
Use a comma to separate multiple URLs.

To configure the client, pass a Config object to the NewClient function:

		cfg := es6.Config{
		  Addresses: []string{
		    "http://localhost:9200",
		    "http://localhost:9201",
		  },
		  Username: "foo",
		  Password: "bar",
		  Transport: &http.Transport{
		    MaxIdleConnsPerHost:   10,
		    ResponseHeaderTimeout: time.Second,
		    DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		    TLSClientConfig: &tls.Config{
		      MinVersion:         tls.VersionTLS11,
		    },
		  },
		}

		es6.NewClient(cfg)

When using the Elastic Service (https://elastic.co/cloud), you can use CloudID instead of Addresses.
When either Addresses or CloudID is set, the ELASTICSEARCH_URL environment variable is ignored.

See the elasticsearch_integration_test.go file and the _examples folder for more information.

Call the Elasticsearch APIs by invoking the corresponding methods on the client:

		res, err := es.Info()
		if err != nil {
		  log.Fatalf("Error getting response: %s", err)
		}

		log.Println(res)

See the github.com/elastic/go-elasticsearch/esapi package for more information about using the API.

See the github.com/elastic/go-elasticsearch/estransport package for more information about configuring the transport.
*/
package es6
