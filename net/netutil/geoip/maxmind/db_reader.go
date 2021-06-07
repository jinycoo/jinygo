/**------------------------------------------------------------**
 * @filename maxmind/db_writer.go
 * @author   jiny - caojingyin@baimaohui.net
 * @version  1.0.0
 * @date     2020/11/6 12:22
 * @desc     maxmind-db_writer - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"jinycoo.com/jinygo/ctime"
	"jinycoo.com/jinygo/errors"
)

type Reader interface {
	Get(destination Writer, editionID string) error
}

type HTTPDatabaseReader struct {
	client            *http.Client
	retryFor          ctime.Duration
	url               string
	licenseKey        string
	accountID         int
	preserveFileTimes bool
	verbose           bool
}

func (reader *HTTPDatabaseReader) Get(destination Writer, editionID string) error {
	defer func() {
		if err := destination.Close(); err != nil {
			log.Println(err)
		}
	}()

	maxMindURL := fmt.Sprintf(
		"%s/geoip/databases/%s/update?db_md5=%s",
		reader.url,
		url.PathEscape(editionID),
		url.QueryEscape(destination.GetHash()),
	)

	req, err := http.NewRequest(http.MethodGet, maxMindURL, nil) // nolint: noctx
	if err != nil {
		return errors.Wrap(err, "error creating request")
	}
	req.SetBasicAuth(fmt.Sprintf("%d", reader.accountID), reader.licenseKey)

	if reader.verbose {
		log.Printf("Performing update request to %s", maxMindURL)
	}
	response, err := MaybeRetryRequest(reader.client, time.Duration(reader.retryFor), req)
	if err != nil {
		return errors.Wrap(err, "error performing HTTP request")
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Fatalf("Error closing response body: %+v", errors.Wrap(err, "closing body"))
		}
	}()

	if response.StatusCode == http.StatusNotModified {
		if reader.verbose {
			log.Printf("No new updates available for %s", editionID)
		}
		return nil
	}

	if response.StatusCode != http.StatusOK {
		buf, err := io.ReadAll(io.LimitReader(response.Body, 256))
		if err == nil {
			return errors.Errorf("unexpected HTTP status code: %s: %s", response.Status, buf)
		}
		return errors.Errorf("unexpected HTTP status code: %s", response.Status)
	}

	gzReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return errors.Wrap(err, "encountered an error creating GZIP reader")
	}
	defer func() {
		if err := gzReader.Close(); err != nil {
			log.Printf("error closing gzip reader: %s", err)
		}
	}()

	if _, err = io.Copy(destination, gzReader); err != nil { //nolint:gosec
		return errors.Wrap(err, "error writing response")
	}

	newMD5 := response.Header.Get("X-Database-MD5")
	if newMD5 == "" {
		return errors.New("no X-Database-MD5 header found")
	}
	if err := destination.ValidHash(newMD5); err != nil {
		return err
	}

	if err := destination.Commit(); err != nil {
		return errors.Wrap(err, "encountered an issue committing database update")
	}

	if reader.preserveFileTimes {
		modificationTime, err := lastModified(response.Header.Get("Last-Modified"))
		if err != nil {
			return errors.Wrap(err, "unable to get last modified time")
		}
		err = destination.SetFileModificationTime(modificationTime)
		if err != nil {
			return errors.Wrap(err, "unable to set modification time")
		}
	}

	return nil
}

func lastModified(lastModified string) (time.Time, error) {
	if lastModified == "" {
		return time.Time{}, errors.New("no Last-Modified header found")
	}

	t, err := time.ParseInLocation(time.RFC1123, lastModified, time.UTC)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "error parsing time")
	}

	return t, nil
}
