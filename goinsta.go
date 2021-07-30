package goinsta

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Instagram represent the main API handler
//
// Profiles: Represents instragram's user profile.
// Account:  Represents instagram's personal account.
// Search:   Represents instagram's search.
// Timeline: Represents instagram's timeline.
// Activity: Represents instagram's user activity.
// Inbox:    Represents instagram's messages.
// Location: Represents instagram's locations.
//
// See Scheme section in README.md for more information.
//
// We recommend to use Export and Import functions after first Login.
//
// Also you can use SetProxy and UnsetProxy to set and unset proxy.
// Golang also provides the option to set a proxy using HTTP_PROXY env var.
type Instagram struct {
	user string
	pass string
	// device id: android-1923fjnma8123
	dID string
	// family device id: 71cd1aec-e146-4380-8d60-d216127c7b4e
	fID string
	// uuid: 8493-1233-4312312-5123
	uuid string
	// rankToken
	rankToken string
	// token -- I think this is depricated, as I don't see any csrf tokens being used anymore, but not 100% sure
	token string
	// phone id
	pid string
	// ads id
	adid string
	// challenge URL
	challengeURL string
	// pigeonSessionId
	psID string
	// contains header options set by Instagram
	headerOptions map[string]string
	// expiry of X-Mid cookie
	xmidExpiry int64

	// Instagram objects

	// Timeline allows to receive timeline media.
	Timeline *Timeline
	// Discover provides access to the discover/explore page
	Discover *Discover
	// Profiles is the user interaction
	Profiles *Profiles
	// Account stores all personal data of the user and his/her options.
	Account *Account
	// Searchbar performs searching of multiple things (users, locations...)
	Searchbar *Search
	// Activity are instagram notifications.
	Activity *Activity
	// Inbox are instagram message/chat system.
	Inbox *Inbox
	// Feed for search over feeds
	Feed *Feed
	// User contacts from mobile address book
	Contacts *Contacts
	// Location instance
	Locations *LocationInstance
	// Challenge controls security side of account (Like sms verify / It was me)
	Challenge *Challenge

	c *http.Client

	// Non fatal err handler, which don't get returned
	// By default they will be printed out, you can e.g. pass them to a logger
	ErrHandler func(...interface{})
}

func DefaultErrHandler(args ...interface{}) {
	fmt.Println(args...)
}

func (insta *Instagram) SetErrorHandler(f func(...interface{})) {
	insta.ErrHandler = f
}

// SetHTTPClient sets http client.  This further allows users to use this functionality
// for HTTP testing using a mocking HTTP client Transport, which avoids direct calls to
// the Instagram, instead of returning mocked responses.
func (insta *Instagram) SetHTTPClient(client *http.Client) {
	insta.c = client
}

// SetHTTPTransport sets http transport. This further allows users to tweak the underlying
// low level transport for adding additional fucntionalities.
func (insta *Instagram) SetHTTPTransport(transport http.RoundTripper) {
	insta.c.Transport = transport
}

// SetDeviceID sets device id
func (insta *Instagram) SetDeviceID(id string) {
	insta.dID = id
}

// SetUUID sets uuid
func (insta *Instagram) SetUUID(uuid string) {
	insta.uuid = uuid
}

// SetPhoneID sets phone id
func (insta *Instagram) SetPhoneID(id string) {
	insta.pid = id
}

// SetCookieJar sets the Cookie Jar. This further allows to use a custom implementation
// of a cookie jar which may be backed by a different data store such as redis.
func (insta *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(instaAPIUrl)
	if err != nil {
		return err
	}
	// First grab the cookies from the existing jar and we'll put it in the new jar.
	cookies := insta.c.Jar.Cookies(url)
	insta.c.Jar = jar
	insta.c.Jar.SetCookies(url, cookies)
	return nil
}

// New creates Instagram structure
func New(username, password string) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	insta := &Instagram{
		user: username,
		pass: password,
		dID: generateDeviceID(
			generateMD5Hash(username + password),
		),
		uuid:          generateUUID(), // both uuid must be differents
		pid:           generateUUID(),
		fID:           generateUUID(),
		psID:          "UFS-" + generateUUID() + "-0",
		headerOptions: map[string]string{"X-Ig-Www-Claim": "0"},
		xmidExpiry:    -1,
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Jar: jar,
		},
		ErrHandler: DefaultErrHandler,
	}
	insta.init()

	return insta
}

func (insta *Instagram) init() {
	insta.Challenge = newChallenge(insta)
	insta.Profiles = newProfiles(insta)
	insta.Activity = newActivity(insta)
	insta.Timeline = newTimeline(insta)
	insta.Searchbar = newSearch(insta)
	insta.Inbox = newInbox(insta)
	insta.Feed = newFeed(insta)
	insta.Contacts = newContacts(insta)
	insta.Locations = newLocation(insta)
	insta.Discover = newDiscover(insta)
}

// SetProxy sets proxy for connection.
func (insta *Instagram) SetProxy(url string, insecure bool, forceHTTP2 bool) error {
	uri, err := neturl.Parse(url)
	if err == nil {
		insta.c.Transport = &http.Transport{
			Proxy:             http.ProxyURL(uri),
			ForceAttemptHTTP2: forceHTTP2,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}
	}
	return err
}

// UnsetProxy unsets proxy for connection.
func (insta *Instagram) UnsetProxy() {
	insta.c.Transport = nil
}

// Save exports config to ~/.goinsta
func (insta *Instagram) Save() error {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("home") // for plan9
	}
	return insta.Export(filepath.Join(home, ".goinsta"))
}

// Export exports *Instagram object options
func (insta *Instagram) Export(path string) error {
	url, err := neturl.Parse(instaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:            insta.Account.ID,
		User:          insta.user,
		DeviceID:      insta.dID,
		FamilyID:      insta.fID,
		UUID:          insta.uuid,
		RankToken:     insta.rankToken,
		Token:         insta.token,
		PhoneID:       insta.pid,
		XmidExpiry:    insta.xmidExpiry,
		HeaderOptions: insta.headerOptions,
		Cookies:       insta.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0o644)
}

// Export exports selected *Instagram object options to an io.Writer
func Export(insta *Instagram, writer io.Writer) error {
	url, err := neturl.Parse(instaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:            insta.Account.ID,
		User:          insta.user,
		DeviceID:      insta.dID,
		FamilyID:      insta.fID,
		UUID:          insta.uuid,
		RankToken:     insta.rankToken,
		Token:         insta.token,
		PhoneID:       insta.pid,
		XmidExpiry:    insta.xmidExpiry,
		HeaderOptions: insta.headerOptions,
		Cookies:       insta.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

// ImportReader imports instagram configuration from io.Reader
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportReader(r io.Reader) (*Instagram, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := ConfigFile{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return ImportConfig(config)
}

// ImportConfig imports instagram configuration from a configuration object.
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportConfig(config ConfigFile) (*Instagram, error) {
	url, err := neturl.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	insta := &Instagram{
		user:          config.User,
		dID:           config.DeviceID,
		fID:           config.FamilyID,
		uuid:          config.UUID,
		rankToken:     config.RankToken,
		token:         config.Token,
		pid:           config.PhoneID,
		xmidExpiry:    config.XmidExpiry,
		headerOptions: config.HeaderOptions,
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
		ErrHandler: DefaultErrHandler,
	}
	insta.c.Jar, err = cookiejar.New(nil)
	if err != nil {
		return insta, err
	}
	insta.c.Jar.SetCookies(url, config.Cookies)

	insta.init()
	insta.Account = &Account{insta: insta, ID: config.ID}
	insta.Account.Sync()

	return insta, nil
}

// Import imports instagram configuration
//
// This function does not set proxy automatically. Use SetProxy after this call.
func Import(path string) (*Instagram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportReader(f)
}

// Login performs instagram login sequence in close resemblance to the android apk.
//
// Password will be deleted after login
func (insta *Instagram) Login() (err error) {
	// pre-login sequence
	err = insta.zrToken()
	if err != nil {
		return
	}
	_, _, err = insta.sync()
	if err != nil {
		return
	}

	err = insta.getPrefill()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching prefill:", err)
	}

	err = insta.contactPrefill()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching contact prefill:", err)
	}

	pkey, pkeyID, err := insta.sync()
	if err != nil {
		return
	}
	if pkey == "" || pkeyID == 0 {
		return errors.New("Sync returned empty public key and/or public key id")
	}

	err = insta.login(pkey, pkeyID)
	if err != nil {
		return err
	}

	// post-login sequence
	err = insta.OpenApp()
	if err != nil {
		return err
	}

	return
}

// Logout closes current session
func (insta *Instagram) Logout() error {
	_, err := insta.sendSimpleRequest(urlLogout)
	insta.c.Jar = nil
	insta.c = nil
	return err
}

func (insta *Instagram) OpenApp() (err error) {
	err = insta.zrToken()
	if err != nil {
		return
	}

	err = insta.getAccountFamily()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching account family:", err)
	}
	_, _, err = insta.syncb()
	if err != nil {
		return
	}

	err = insta.getNdxSteps()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching ndx steps:", err)
	}

	if !insta.Timeline.Next() {
		return errors.New("Failed to fetch timeline during login procedure: " +
			insta.Timeline.err.Error())
	}

	err = insta.callNotifBadge()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching notify badge", err)
	}

	err = insta.banyan()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching banyan", err)
	}

	err = insta.callMediaBlocked()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching blocked media", err)
	}

	// no clue what theses values could be used for
	_, err = insta.getCooldowns()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching cool downs", err)
	}

	if !insta.Discover.Next() {
		insta.ErrHandler("Non fatal error while fetching explore page",
			insta.Discover.Error())
	}

	err = insta.getConfig()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching config", err)
	}

	// no clue what theses values could be used for
	_, err = insta.getScoresBootstrapUsers()
	if err != nil {
		insta.ErrHandler("Non fatal error while fetching bootstrap user scores", err)
	}

	if !insta.Activity.Next() {
		return errors.New("Failed to fetch recent activity: " +
			insta.Activity.err.Error())
	}

	err = insta.sendAdID()
	if err != nil {
		insta.ErrHandler("Non fatal error while sending ad id", err)
	}

	err = insta.callStClPushPerm()
	if err != nil {
		insta.ErrHandler("Non fatal error while calling store client push permissions", err)
	}

	if !insta.Inbox.initialSnapshot() {
		return errors.New("Failed to fetch initial messages inbox snapshot: " +
			insta.Inbox.err.Error())
	}

	err = insta.callContPointSig()
	if err != nil {
		insta.ErrHandler("Non fatal error while calling contact point signal:", err)
	}

	return nil
}

func (insta *Instagram) login(pkey string, pkeyID int) error {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	encrypted, err := EncryptPassword(insta.pass, pkey, pkeyID, timestamp)
	if err != nil {
		return err
	}

	result, err := json.Marshal(
		map[string]interface{}{
			"jazoest":      Jazoest(insta.dID),
			"country_code": "[{\"country_code\":\"44\",\"source\":[\"default\"]}]",
			"phone_id":     insta.fID,
			"enc_password": fmt.Sprintf(
				"#PWD_INSTAGRAM:4:%s:%s",
				timestamp,
				encrypted,
			),
			"username":            insta.user,
			"adid":                insta.adid,
			"guid":                insta.uuid,
			"device_id":           insta.dID,
			"google_tokens":       "[]",
			"login_attempt_count": 0,
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlLogin,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(result)},
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}

	insta.pass = ""
	res := accountResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	insta.Account = &res.Account
	insta.Account.insta = insta
	insta.rankToken = strconv.FormatInt(insta.Account.ID, 10) + "_" + insta.uuid

	return nil
}

func (insta *Instagram) getPrefill() error {
	data, err := json.Marshal(
		map[string]string{
			"android_device_id": insta.dID,
			"phone_id":          insta.fID,
			"usages":            "[\"account_recovery_omnibox\"]",
			"device_id":         insta.uuid,
		},
	)
	if err != nil {
		return err
	}

	// ignore the error returned by the request, because 429 if often returned
	insta.sendRequest(
		&reqOptions{
			Endpoint: urlGetPrefill,
			IsPost:   true,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(data)},
		},
	)
	return nil
}

func (insta *Instagram) contactPrefill() error {
	data, err := json.Marshal(
		map[string]string{
			"phone_id": insta.fID,
			"usage":    "prefill",
		},
	)
	if err != nil {
		return err
	}

	// ignore the error returned by the request, because 429 if often returned
	insta.sendRequest(
		&reqOptions{
			Endpoint: urlContactPrefill,
			IsPost:   true,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(data)},
		},
	)
	return nil
}

func (insta *Instagram) zrToken() error {
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlZrToken,
			IsPost:   false,
			Query: map[string]string{
				"device_id":        insta.dID,
				"token_hash":       "",
				"custom_device_id": insta.uuid,
				"fetch_reason":     "token_expired",
			},
			IgnoreHeaders: []string{
				"X-Pigeon-Session-Id",
				"X-Pigeon-Rawclienttime",
				"X-Ig-App-Locale",
				"X-Ig-Device-Locale",
				"X-Ig-Mapped-Locale",
				"X-Ig-App-Startup-Country",
			},
		},
	)

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	// Get the expiry time of the token
	token := res["token"].(map[string]interface{})
	ttl := token["ttl"].(float64)
	t := token["request_time"].(float64)
	insta.xmidExpiry = int64(t + ttl)

	return err
}

func (insta *Instagram) sendAdID() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlLogAttribution,
			IsPost:   true,
			Query:    map[string]string{"signed_body": "SIGNATURE.{}"},
		},
	)
	return err
}

func (insta *Instagram) callStClPushPerm() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlStoreClientPushPermissions,
			IsPost:   true,
			Query: map[string]string{
				"enabled":   "true",
				"device_id": insta.uuid,
				"_uuid":     insta.uuid,
			},
		},
	)
	return err
}

func (insta *Instagram) sync() (string, int, error) {
	data, err := json.Marshal(
		map[string]interface{}{
			"id":                      insta.uuid,
			"server_config_retrieval": 1,
		},
	)
	if err != nil {
		return "", -1, err
	}

	_, h, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSync,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(data)},
			IsPost:   true,
			IgnoreHeaders: []string{
				"Authorization",
			},
		},
	)

	hkey := h["Ig-Set-Password-Encryption-Pub-Key"]
	hkeyID := h["Ig-Set-Password-Encryption-Key-Id"]
	var key string
	var keyID string
	if len(hkey) > 0 && len(hkeyID) > 0 && hkey[0] != "" && hkeyID[0] != "" {
		key = hkey[0]
		keyID = hkeyID[0]
	}

	id, err := strconv.Atoi(keyID)
	if err != nil {
		return "", -1, err
	}

	return key, id, err
}

// syncb calls b.i.instagram.com, with auth data after login
func (insta *Instagram) syncb() (string, int, error) {
	data, err := json.Marshal(
		map[string]interface{}{
			"id":                      insta.Account.ID,
			"_id":                     insta.Account.ID,
			"_uuid":                   insta.uuid,
			"server_config_retrieval": 1,
		},
	)
	if err != nil {
		return "", -1, err
	}

	_, h, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlSync,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(data)},
			IsPost:   true,
			Useb:     true,
			IgnoreHeaders: []string{
				"X-Pigeon-Session-Id",
				"X-Pigeon-Rawclienttime",
			},
		},
	)

	hkey := h["Ig-Set-Password-Encryption-Pub-Key"]
	hkeyID := h["Ig-Set-Password-Encryption-Key-Id"]
	var key string
	var keyID string
	if len(hkey) > 0 && len(hkeyID) > 0 && hkey[0] != "" && hkeyID[0] != "" {
		key = hkey[0]
		keyID = hkeyID[0]
	}

	id, err := strconv.Atoi(keyID)
	if err != nil {
		return "", -1, err
	}

	return key, id, err
}

func (insta *Instagram) getAccountFamily() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlGetAccFamily,
		},
	)
	return err
}

func (insta *Instagram) getNdxSteps() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlGetNdxSteps,
		},
	)
	return err
}

func (insta *Instagram) banyan() error {
	// TODO: process body, and put the data in a struct
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlBanyan,
			Query: map[string]string{
				"views": `["story_share_sheet","direct_user_search_nullstate","forwarding_recipient_sheet","threads_people_picker","direct_inbox_active_now","group_stories_share_sheet","call_recipients","reshare_share_sheet","direct_user_search_keypressed"]`,
			},
		},
	)
	return err
}

func (insta *Instagram) callNotifBadge() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlNotifBadge,
			IsPost:   true,
			Query: map[string]string{
				"phone_id":  insta.fID,
				"user_ids":  strconv.Itoa(int(insta.Account.ID)),
				"device_id": insta.uuid,
				"_uuid":     insta.uuid,
			},
		},
	)
	return err
}

func (insta *Instagram) callContPointSig() error {
	query := map[string]string{
		"phone_id":      insta.fID,
		"_uid":          strconv.Itoa(int(insta.Account.ID)),
		"device_id":     insta.uuid,
		"_uuid":         insta.uuid,
		"google_tokens": "[]",
	}
	b, err := json.Marshal(query)
	if err != nil {
		return err
	}
	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: urlProcessContactPointSignals,
			IsPost:   true,
			Query:    map[string]string{"signed_body": "SIGNATURE." + string(b)},
		},
	)
	return err
}

func (insta *Instagram) callMediaBlocked() error {
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlMediaBlocked,
		},
	)
	return err
}

func (insta *Instagram) getCooldowns() (*Cooldowns, error) {
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCooldowns,
			Query: map[string]string{
				"signed_body": "SIGNATURE.{}",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// No clue what to use these values for
	temp := Cooldowns{}
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

func (insta *Instagram) getScoresBootstrapUsers() (*ScoresBootstrapUsers, error) {
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlCooldowns,
			Query: map[string]string{
				"surfaces": `["autocomplete_user_list","coefficient_besties_list_ranking","coefficient_rank_recipient_user_suggestion","coefficient_ios_section_test_bootstrap_ranking","coefficient_direct_recipients_ranking_variant_2"]`,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// No clue what to use these values for
	s := ScoresBootstrapUsers{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	for _, u := range s.Users {
		u.insta = insta
	}
	return &s, nil
}

func (insta *Instagram) getConfig() error {
	// returns a bunch of values with single letter labels
	// see unparsedResp/loom_fetch_config/*.json for examples
	_, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlFetchConfig,
		},
	)
	return err
}

func (insta *Instagram) expose() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"id":         insta.Account.ID,
			"experiment": "ig_android_profile_contextual_feed",
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: urlExpose,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)

	return err
}

// GetMedia returns media specified by id.
//
// The argument can be int64 or string
//
// See example: examples/media/like.go
func (insta *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
	media := &FeedMedia{
		insta:  insta,
		NextID: o,
	}
	return media, media.Sync()
}
