/**------------------------------------------------------------**
 * @filename maxmind/update.go
 * @author   jiny - caojingyin@baimaohui.net
 * @version  1.0.0
 * @date     2020/11/6 14:01
 * @desc     maxmind-update - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"jinycoo.com/jinygo/errors"
	"jinycoo.com/jinygo/log"
	"jinycoo.com/jinygo/utils"
)

const (
	ZeroMD5 = "00000000000000000000000000000000"
)

type MMDbUWork struct {
	Config *Config
	Client *http.Client
}

func New(config *Config) *MMDbUWork {
	config.Init()
	transport := http.DefaultTransport
	if len(config.Proxy) > 0 {
		px, _ := parseProxy(config.Proxy, config.ProxyUserPassword)
		proxy := http.ProxyURL(px)
		transport.(*http.Transport).Proxy = proxy
	}
	return &MMDbUWork{
		Config: config,
		Client: &http.Client{Transport: transport},
	}
}

func (uw *MMDbUWork) Run() error {
	var dbReader = &HTTPDatabaseReader{
		client:            uw.Client,
		retryFor:          uw.Config.RetryFor,
		url:               uw.Config.URL,
		licenseKey:        uw.Config.LicenseKey,
		accountID:         uw.Config.AccountID,
		preserveFileTimes: uw.Config.PreserveFileTimes,
		verbose:           uw.Config.Verbose,
	}

	for _, editionID := range uw.Config.EditionIDs {
		filename, err := uw.GetFilename(editionID)
		fmt.Println(filename)
		if err != nil {
			return errors.Wrapf(err, "error retrieving filename for %s", editionID)
		}
		filePath := filepath.Join(utils.RootDir(), uw.Config.DatabaseDirectory, "current", filename)
		fmt.Println(filePath)
		dbWriter, err := NewLocalFileDatabaseWriter(filePath, uw.Config.LockFile, uw.Config.Verbose)
		if err != nil {
			return errors.Wrapf(err, "error creating database writer for %s", editionID)
		}
		if err := dbReader.Get(dbWriter, editionID); err != nil {
			return errors.WithMessagef(err, "error while getting database for %s", editionID)
		}
	}
	return nil
}

func (uw *MMDbUWork) GetFilename(editionID string) (string, error) {
	maxMindURL := fmt.Sprintf("%s/app/update_getfilename?product_id=%s", uw.Config.URL, url.QueryEscape(editionID))

	if uw.Config.Verbose {
		log.Infof("Performing get filename request to %s", maxMindURL)
	}
	req, err := http.NewRequest(http.MethodGet, maxMindURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "error creating HTTP request")
	}
	res, err := MaybeRetryRequest(uw.Client, time.Duration(uw.Config.RetryFor), req)
	if err != nil {
		return "", errors.Wrap(err, "error performing HTTP request")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatalf("error closing response body: %+v", errors.Wrap(err, "closing body"))
		}
	}()

	buf, err := io.ReadAll(io.LimitReader(res.Body, 256))
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("unexpected HTTP status code: %s: %s", res.Status, buf)
	}

	if len(buf) == 0 {
		return "", errors.New("response body is empty")
	}

	if bytes.Count(buf, []byte("\n")) > 0 ||
		bytes.Count(buf, []byte("\x00")) > 0 {
		return "", errors.New("invalid characters in filename")
	}

	return string(buf), nil
}

func MaybeRetryRequest(c *http.Client, retryFor time.Duration, req *http.Request) (*http.Response, error) {
	if retryFor < 0 {
		return nil, errors.New("negative retry duration")
	}
	if req.Body != nil {
		return nil, errors.New("can't retry requests with bodies")
	}
	var resp *http.Response
	var err error

	start := time.Now()
	for i := uint(0); ; i++ {
		resp, err = c.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		currentDuration := time.Since(start)

		waitDuration := 200 * time.Millisecond * (1 << i)
		if currentDuration+waitDuration > retryFor/time.Millisecond {
			break
		}
		if err == nil {
			_ = resp.Body.Close()
		}
		time.Sleep(waitDuration)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error performing http request")
	}
	return resp, nil
}
