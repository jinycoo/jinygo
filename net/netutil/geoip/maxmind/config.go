/**------------------------------------------------------------**
 * @filename maxmind/config.go
 * @author   jiny - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2020/11/6 10:01
 * @desc     maxmind-config - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jinycoo/jinygo/ctime"
	"github.com/jinycoo/jinygo/errors"
)

const (
	defaultDatabaseDirectory = "/tmp"
)

type Config struct {
	AccountID         int
	DatabaseDirectory string
	LicenseKey        string
	LockFile          string
	URL               string
	EditionIDs        []string
	Proxy             string
	ProxyUserPassword string
	PreserveFileTimes bool
	Verbose           bool
	RetryFor          ctime.Duration
}

func (cfg *Config) Init() (err error) {

	if cfg.DatabaseDirectory == "" {
		cfg.DatabaseDirectory = filepath.Clean(defaultDatabaseDirectory)
	} else {
		cfg.DatabaseDirectory = filepath.Clean(cfg.DatabaseDirectory)
	}

	if cfg.URL == "" {
		cfg.URL = "https://updates.maxmind.com"
	}

	if cfg.LockFile == "" {
		cfg.LockFile = filepath.Join(cfg.DatabaseDirectory, ".geoipupdate.lock")
	} else {
		cfg.LockFile = filepath.Clean(cfg.LockFile)
	}
	return
}

var schemeRE = regexp.MustCompile(`(?i)\A([a-z][a-z0-9+\-.]*)://`)

func parseProxy(proxy, proxyUserPassword string) (*url.URL, error) {
	if proxy == "" {
		return nil, nil
	}

	matches := schemeRE.FindStringSubmatch(proxy)
	if matches == nil {
		proxy = "http://" + proxy
	} else {
		scheme := strings.ToLower(matches[1])
		if scheme != "http" && scheme != "socks5" {
			return nil, errors.Errorf("unsupported proxy type: %s", scheme)
		}
	}

	u, err := url.Parse(proxy)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing proxy URL")
	}

	if !strings.Contains(u.Host, ":") {
		u.Host += ":1080" // The 1080 default historically came from cURL.
	}

	// Historically if the Proxy option had a username and password they would
	// override any specified in the ProxyUserPassword option. Continue that.
	if u.User != nil {
		return u, nil
	}

	if proxyUserPassword == "" {
		return u, nil
	}

	userPassword := strings.SplitN(proxyUserPassword, ":", 2)
	if len(userPassword) != 2 {
		return nil, errors.New("proxy user/password is malformed")
	}
	u.User = url.UserPassword(userPassword[0], userPassword[1])

	return u, nil
}
