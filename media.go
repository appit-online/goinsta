package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	neturl "net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StoryReelMention represent story reel mention
type StoryReelMention struct {
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           int     `json:"z"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Rotation    float64 `json:"rotation"`
	IsPinned    int     `json:"is_pinned"`
	IsHidden    int     `json:"is_hidden"`
	IsSticker   int     `json:"is_sticker"`
	IsFBSticker int     `json:"is_fb_sticker"`
	User        User
	DisplayType string `json:"display_type"`
}

// StoryCTA represent story cta
type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

// Item represents media items
//
// All Item has Images or Videos objects which contains the url(s).
// You can use Download function to get the best quality Image or Video from Item.
type Item struct {
	insta    *Instagram
	media    Media
	Comments *Comments `json:"-"`

	TakenAt           int64  `json:"taken_at"`
	Pk                int64  `json:"pk"`
	ID                string `json:"id"`
	CommentsDisabled  bool   `json:"comments_disabled"`
	DeviceTimestamp   int64  `json:"device_timestamp"`
	FacepileTopLikers []struct {
		FollowFrictionType float64 `json:"follow_friction_type"`
		FullNeme           string  `json:"ful_name"`
		IsPrivate          bool    `json:"is_private"`
		IsVerified         bool    `json:"is_verified"`
		Pk                 float64 `json:"pk"`
		ProfilePicID       string  `json:"profile_pic_id"`
		ProfilePicURL      string  `json:"profile_pic_url"`
		Username           string  `json:"username"`
	} `json:"facepile_top_likers"`
	MediaType             int     `json:"media_type"`
	Code                  string  `json:"code"`
	ClientCacheKey        string  `json:"client_cache_key"`
	FilterType            int     `json:"filter_type"`
	CarouselParentID      string  `json:"carousel_parent_id"`
	CarouselMedia         []Item  `json:"carousel_media,omitempty"`
	User                  User    `json:"user"`
	CanViewerReshare      bool    `json:"can_viewer_reshare"`
	Caption               Caption `json:"caption"`
	CaptionIsEdited       bool    `json:"caption_is_edited"`
	LikeViewCountDisabled bool    `json:"like_and_view_counts_disabled"`
	IsCommercial          bool    `json:"is_commercial"`
	CommercialityStatus   string  `json:"commerciality_status"`
	FundraiserTag         struct {
		HasStandaloneFundraiser bool `json:"has_standalone_fundraiser"`
	} `json:"fundraiser_tag"`
	IsPaidPartnership bool   `json:"is_paid_partnership"`
	ProductType       string `json:"product_type"`
	Likes             int    `json:"like_count"`
	HasLiked          bool   `json:"has_liked"`
	// Toplikers can be `string` or `[]string`.
	// Use TopLikers function instead of getting it directly.
	Toplikers                    interface{} `json:"top_likers"`
	Likers                       []User      `json:"likers"`
	CommentLikesEnabled          bool        `json:"comment_likes_enabled"`
	CommentThreadingEnabled      bool        `json:"comment_threading_enabled"`
	HasMoreComments              bool        `json:"has_more_comments"`
	MaxNumVisiblePreviewComments int         `json:"max_num_visible_preview_comments"`
	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	// Use PreviewComments function instead of getting it directly.
	Previewcomments interface{} `json:"preview_comments,omitempty"`
	CommentCount    int         `json:"comment_count"`
	PhotoOfYou      bool        `json:"photo_of_you"`
	// Tags are tagged people in photo
	Tags struct {
		In []Tag `json:"in"`
	} `json:"usertags,omitempty"`
	FbUserTags           Tag    `json:"fb_user_tags"`
	CanViewerSave        bool   `json:"can_viewer_save"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	// Images contains URL images in different versions.
	// Version = quality.
	Images          Images   `json:"image_versions2,omitempty"`
	OriginalWidth   int      `json:"original_width,omitempty"`
	OriginalHeight  int      `json:"original_height,omitempty"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty"`
	Location        Location `json:"location,omitempty"`
	Lat             float64  `json:"lat,omitempty"`
	Lng             float64  `json:"lng,omitempty"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty"`
	VideoCodec        string  `json:"video_codec"`
	HasAudio          bool    `json:"has_audio,omitempty"`
	VideoDuration     float64 `json:"video_duration,omitempty"`
	ViewCount         float64 `json:"view_count,omitempty"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty"`
	IsUnifiedVideo    bool    `json:"is_unified_video"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty"`

	// Only for stories
	StoryEvents              []interface{}      `json:"story_events"`
	StoryHashtags            []interface{}      `json:"story_hashtags"`
	StoryPolls               []interface{}      `json:"story_polls"`
	StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	StorySoundOn             []interface{}      `json:"story_sound_on"`
	CreativeConfig           interface{}        `json:"creative_config"`
	StoryLocations           []interface{}      `json:"story_locations"`
	StorySliders             []interface{}      `json:"story_sliders"`
	StoryQuestions           []interface{}      `json:"story_questions"`
	StoryProductItems        []interface{}      `json:"story_product_items"`
	StoryCTA                 []StoryCTA         `json:"story_cta"`
	IntegrityReviewDecision  string             `json:"integrity_review_decision"`
	IsReelMedia              bool               `json:"is_reel_media"`
	ProfileGridControl       bool               `json:"profile_grid_control_enabled"`
	ReelMentions             []StoryReelMention `json:"reel_mentions"`
	ExpiringAt               int64              `json:"expiring_at"`
	CanSendCustomEmojis      bool               `json:"can_send_custom_emojis"`
	SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int64              `json:"has_shared_to_fb"`
	Mentions                 []Mentions
	Audience                 string `json:"audience,omitempty"`
	StoryMusicStickers       []struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Z              int     `json:"z"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Rotation       float64 `json:"rotation"`
		IsPinned       int     `json:"is_pinned"`
		IsHidden       int     `json:"is_hidden"`
		IsSticker      int     `json:"is_sticker"`
		MusicAssetInfo struct {
			ID                       string `json:"id"`
			Title                    string `json:"title"`
			Subtitle                 string `json:"subtitle"`
			DisplayArtist            string `json:"display_artist"`
			CoverArtworkURI          string `json:"cover_artwork_uri"`
			CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadURL   string `json:"progressive_download_url"`
			HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
			IsExplicit               bool   `json:"is_explicit"`
			DashManifest             string `json:"dash_manifest"`
			HasLyrics                bool   `json:"has_lyrics"`
			AudioAssetID             string `json:"audio_asset_id"`
			IgArtist                 struct {
				Pk            int    `json:"pk"`
				Username      string `json:"username"`
				FullName      string `json:"full_name"`
				IsPrivate     bool   `json:"is_private"`
				ProfilePicURL string `json:"profile_pic_url"`
				ProfilePicID  string `json:"profile_pic_id"`
				IsVerified    bool   `json:"is_verified"`
			} `json:"ig_artist"`
			PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio          bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
			OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
		} `json:"music_asset_info"`
	} `json:"story_music_stickers,omitempty"`
}

// Comment pushes a text comment to media item.
//
// If parent media is a Story this function will send a private message
// replying the Instagram story.
func (item *Item) Comment(text string) error {
	var opt *reqOptions
	var err error
	insta := item.media.instagram()

	switch item.media.(type) {
	case *StoryMedia:
		to, err := prepareRecipients(item)
		if err != nil {
			return err
		}

		query := insta.prepareDataQuery(
			map[string]interface{}{
				"recipient_users": to,
				"action":          "send_item",
				"media_id":        item.ID,
				"client_context":  generateUUID(),
				"text":            text,
				"entry":           "reel",
				"reel_id":         item.User.ID,
			},
		)
		opt = &reqOptions{
			Connection: "keep-alive",
			Endpoint:   fmt.Sprintf("%s?media_type=%s", urlReplyStory, item.MediaToString()),
			Query:      query,
			IsPost:     true,
		}
	case *FeedMedia: // normal media
		var data string
		data, err = insta.prepareData(
			map[string]interface{}{
				"comment_text": text,
			},
		)
		opt = &reqOptions{
			Endpoint: fmt.Sprintf(urlCommentAdd, item.Pk),
			Query:    generateSignature(data),
			IsPost:   true,
		}
	}
	if err != nil {
		return err
	}

	// ignoring response
	_, _, err = insta.sendRequest(opt)
	return err
}

// MediaToString returns Item.MediaType as string.
func (item *Item) MediaToString() string {
	switch item.MediaType {
	case 1:
		return "photo"
	case 2:
		return "video"
	case 8:
		return "carousel"
	}
	return ""
}

func setToItem(item *Item, media Media) {
	item.media = media
	item.User.insta = media.instagram()
	item.Comments = newComments(item)
	for i := range item.CarouselMedia {
		item.CarouselMedia[i].User = item.User
		setToItem(&item.CarouselMedia[i], media)
	}
}

// setToMediaItem is a utility function that
// mimics the setToItem but for the SavedMedia items
func setToMediaItem(item *MediaItem, media Media) {
	item.Media.media = media
	item.Media.User.insta = media.instagram()

	item.Media.Comments = newComments(&item.Media)

	for i := range item.Media.CarouselMedia {
		item.Media.CarouselMedia[i].User = item.Media.User
		setToItem(&item.Media.CarouselMedia[i], media)
	}
}

func getname(name string) string {
	nname := name
	i := 1
	for {
		ext := path.Ext(name)

		_, err := os.Stat(name)
		if err != nil {
			break
		}
		if ext != "" {
			nname = strings.Replace(nname, ext, "", -1)
		}
		name = fmt.Sprintf("%s.%d%s", nname, i, ext)
		i++
	}
	return name
}

func download(insta *Instagram, url, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := insta.c.Get(url)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	return dst, err
}

type bestMedia struct {
	w, h int
	url  string
}

// GetBest returns best quality image or video.
//
// Arguments can be []Video or []Candidate
func GetBest(obj interface{}) string {
	m := bestMedia{}

	switch t := obj.(type) {
	// getting best video
	case []Video:
		for _, video := range t {
			if m.w < video.Width && video.Height > m.h && video.URL != "" {
				m.w = video.Width
				m.h = video.Height
				m.url = video.URL
			}
		}
		// getting best image
	case []Candidate:
		for _, image := range t {
			if m.w < image.Width && image.Height > m.h && image.URL != "" {
				m.w = image.Width
				m.h = image.Height
				m.url = image.URL
			}
		}
	}
	return m.url
}

var rxpTags = regexp.MustCompile(`#\w+`)

// Hashtags returns caption hashtags.
//
// Item media parent must be FeedMedia.
//
// See example: examples/media/hashtags.go
func (item *Item) Hashtags() []Hashtag {
	tags := rxpTags.FindAllString(item.Caption.Text, -1)

	hsh := make([]Hashtag, len(tags))

	i := 0
	for _, tag := range tags {
		hsh[i].Name = tag[1:]
		i++
	}

	for _, comment := range item.PreviewComments() {
		tags := rxpTags.FindAllString(comment.Text, -1)

		for _, tag := range tags {
			hsh = append(hsh, Hashtag{Name: tag[1:]})
		}
	}

	return hsh
}

// Delete deletes your media item. StoryMedia or FeedMedia
//
// See example: examples/media/mediaDelete.go
func (item *Item) Delete() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaDelete, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// SyncLikers fetch new likers of a media
//
// This function updates Item.Likers value
func (item *Item) SyncLikers() error {
	resp := respLikers{}
	insta := item.media.instagram()
	body, err := insta.sendSimpleRequest(urlMediaLikers, item.ID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err == nil {
		item.Likers = resp.Users
	}
	return err
}

// Unlike mark media item as unliked.
//
// See example: examples/media/unlike.go
func (item *Item) Unlike() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaUnlike, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Like mark media item as liked.
//
// See example: examples/media/like.go
func (item *Item) Like() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaLike, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Save saves media item.
//
// You can get saved media using Account.Saved()
func (item *Item) Save() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaSave, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Unsave unsaves media item.
func (item *Item) Unsave() error {
	insta := item.media.instagram()
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": item.ID,
		},
	)
	if err != nil {
		return err
	}

	_, _, err = insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaUnsave, item.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	return err
}

// Download downloads media item (video or image) with the best quality.
//
// Input parameters are folder and filename. If filename is "" will be saved with
// the default value name.
//
// If file exists it will be saved
// This function makes folder automatically
//
// This function returns an slice of location of downloaded items
// The returned values are the output path of images and videos.
//
// This function does not download CarouselMedia.
//
// See example: examples/media/itemDownload.go
func (item *Item) Download(folder, name string) (imgs, vds string, err error) {
	var u *neturl.URL
	var nname string
	imgFolder := path.Join(folder, "images")
	vidFolder := path.Join(folder, "videos")
	insta := item.media.instagram()

	os.MkdirAll(folder, 0o777)
	os.MkdirAll(imgFolder, 0o777)
	os.MkdirAll(vidFolder, 0o777)

	vds = GetBest(item.Videos)
	if vds != "" {
		if name == "" {
			u, err = neturl.Parse(vds)
			if err != nil {
				return
			}

			nname = path.Join(vidFolder, path.Base(u.Path))
		} else {
			nname = path.Join(vidFolder, name)
		}
		nname = getname(nname)

		vds, err = download(insta, vds, nname)
		return "", vds, err
	}

	imgs = GetBest(item.Images.Versions)
	if imgs != "" {
		if name == "" {
			u, err = neturl.Parse(imgs)
			if err != nil {
				return
			}

			nname = path.Join(imgFolder, path.Base(u.Path))
		} else {
			nname = path.Join(imgFolder, name)
		}
		nname = getname(nname)

		imgs, err = download(insta, imgs, nname)
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

// TopLikers returns string slice or single string (inside string slice)
// Depending on TopLikers parameter.
func (item *Item) TopLikers() []string {
	switch s := item.Toplikers.(type) {
	case string:
		return []string{s}
	case []string:
		return s
	}
	return nil
}

// PreviewComments returns string slice or single string (inside Comment slice)
// Depending on PreviewComments parameter.
// If PreviewComments are string or []string only the Text field will be filled.
func (item *Item) PreviewComments() []Comment {
	switch s := item.Previewcomments.(type) {
	case []interface{}:
		if len(s) == 0 {
			return nil
		}

		switch s[0].(type) {
		case interface{}:
			comments := make([]Comment, 0)
			for i := range s {
				if buf, err := json.Marshal(s[i]); err != nil {
					return nil
				} else {
					comment := &Comment{}

					if err = json.Unmarshal(buf, comment); err != nil {
						return nil
					} else {
						comments = append(comments, *comment)
					}
				}
			}
			return comments
		case string:
			comments := make([]Comment, 0)
			for i := range s {
				comments = append(comments, Comment{
					Text: s[i].(string),
				})
			}
			return comments
		}
	case string:
		comments := []Comment{
			{
				Text: s,
			},
		}
		return comments
	}
	return nil
}

// StoryIsCloseFriends returns a bool
// If the returned value is true the story was published only for close friends
func (item *Item) StoryIsCloseFriends() bool {
	return item.Audience == "besties"
}

// Media interface defines methods for both StoryMedia and FeedMedia.
type Media interface {
	// Next allows pagination
	Next(...interface{}) bool
	// Error returns error (in case it have been occurred)
	Error() error
	// ID returns media id
	ID() string
	// Delete removes media
	Delete() error

	instagram() *Instagram
}

// StoryMedia is the struct that handles the information from the methods to get info about Stories.
type StoryMedia struct {
	insta    *Instagram
	endpoint string
	uid      int64

	err error

	Pk                     interface{} `json:"id"`
	MediaCount             int64       `json:"media_count"`
	MediaIDs               []int64     `json:"media_ids"`
	Muted                  bool        `json:"muted"`
	LatestReelMedia        int64       `json:"latest_reel_media"`
	LatestBestiesReelMedia float64     `json:"latest_besties_reel_media"`
	ExpiringAt             float64     `json:"expiring_at"`
	Seen                   float64     `json:"seen"`
	SeenRankedPosition     int         `json:"seen_ranked_position"`
	CanReply               bool        `json:"can_reply"`
	CanGifQuickReply       bool        `json:"can_gif_quick_reply"`
	ClientPrefetchScore    float64     `json:"client_prefetch_score"`
	Title                  string      `json:"title"`
	CanReshare             bool        `json:"can_reshare"`
	ReelType               string      `json:"reel_type"`
	User                   User        `json:"user"`
	Items                  []Item      `json:"items"`
	ReelMentions           []string    `json:"reel_mentions"`
	PrefetchCount          int         `json:"prefetch_count"`
	// this field can be int or bool
	HasBestiesMedia       interface{} `json:"has_besties_media"`
	HasPrideMedia         bool        `json:"has_pride_media"`
	HasVideo              bool        `json:"has_video"`
	IsCacheable           bool        `json:"is_cacheable"`
	IsSensitiveVerticalAd bool        `json:"is_sensitive_vertical_ad"`
	RankedPosition        int         `json:"ranked_position"`
	RankerScores          struct {
		Fp   float64 `json:"fp"`
		Ptap float64 `json:"ptap"`
		Vm   float64 `json:"vm"`
	} `json:"ranker_scores"`
	StoryRankingToken    string      `json:"story_ranking_token"`
	Broadcasts           []Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
	HasNewNuxStory       bool        `json:"has_new_nux_story"`
	Status               string      `json:"status"`
}

// Delete removes instragram story.
//
// See example: examples/media/deleteStories.go
func (media *StoryMedia) Delete() error {
	insta := media.insta
	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": media.ID(),
		},
	)
	if err == nil {
		_, _, err = insta.sendRequest(
			&reqOptions{
				Endpoint: fmt.Sprintf(urlMediaDelete, media.ID()),
				Query:    generateSignature(data),
				IsPost:   true,
			},
		)
	}
	return err
}

// ID returns Story id
func (media *StoryMedia) ID() string {
	switch id := media.Pk.(type) {
	case int64:
		return strconv.FormatInt(id, 10)
	case string:
		return id
	}
	return ""
}

func (media *StoryMedia) instagram() *Instagram {
	return media.insta
}

func (media *StoryMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

// Error returns error happened any error
func (media StoryMedia) Error() error {
	return media.err
}

// Seen marks story as seen.
/*
func (media *StoryMedia) Seen() error {
	insta := media.inst
	data, err := insta.prepareData(
		map[string]interface{}{
			"container_module":   "feed_timeline",
			"live_vods_skipped":  "",
			"nuxes_skipped":      "",
			"nuxes":              "",
			"reels":              "", // TODO xd
			"live_vods":          "",
			"reel_media_skipped": "",
		},
	)
	if err == nil {
		_, _, err = insta.sendRequest(
			&reqOptions{
				Endpoint: urlMediaSeen, // reel=1&live_vod=0
				Query:    generateSignature(data),
				IsPost:   true,
				UseV2:    true,
			},
		)
	}
	return err
}
*/

type trayRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Sync function is used when Highlight must be sync.
// Highlight must be sync when User.Highlights does not return any object inside StoryMedia slice.
//
// This function does NOT update Stories items.
//
// This function updates (fetches) StoryMedia.Items
func (media *StoryMedia) Sync() error {
	insta := media.insta
	query := []trayRequest{
		{"SUPPORTED_SDK_VERSIONS", supportedSdkVersions},
		{"FACE_TRACKER_VERSION", facetrackerVersion},
		{"segmentation", segmentation},
		{"COMPRESSION", compression},
		{"world_tracker", worldTracker},
		{"gyroscope", gyroscope},
	}
	qjson, err := json.Marshal(query)
	if err != nil {
		return err
	}

	id := media.Pk.(string)
	data, err := insta.prepareData(
		map[string]interface{}{
			"exclude_media_ids":          "[]",
			"supported_capabilities_new": string(qjson),
			"source":                     "feed_timeline",
			"_uid":                       insta.Account.ID,
			"_uuid":                      insta.uuid,
			"user_ids":                   []string{id},
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: urlReelMedia,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err == nil {
		resp := trayResp{}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			m, ok := resp.Reels[id]
			if ok {
				media.Items = m.Items
				media.setValues()
				return nil
			}
			err = fmt.Errorf("cannot find %s structure in response", id)
		}
	}
	return err
}

// Next allows pagination after calling:
// User.Stories
//
//
// returns false when list reach the end
// if StoryMedia.Error() is ErrNoMore no problem have been occurred.
func (media *StoryMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.insta
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	body, err := insta.sendSimpleRequest(endpoint)
	if err == nil {
		m := StoryMedia{}
		err = json.Unmarshal(body, &m)
		if err == nil {
			// TODO check NextID media
			*media = m
			media.insta = insta
			media.endpoint = endpoint
			media.err = ErrNoMore // TODO: See if stories has pagination
			media.setValues()
			return true
		}
	}
	media.err = err
	return false
}

// FeedMedia represent a set of media items
type FeedMedia struct {
	insta *Instagram

	err error

	uid       int64
	endpoint  string
	timestamp string

	Items               []Item `json:"items"`
	NumResults          int    `json:"num_results"`
	MoreAvailable       bool   `json:"more_available"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Status              string `json:"status"`
	// Can be int64 and string
	// this is why we recommend Next() usage :')
	NextID interface{} `json:"next_max_id"`
}

// Delete deletes all items in media. Take care...
//
// See example: examples/media/mediaDelete.go
func (media *FeedMedia) Delete() error {
	for i := range media.Items {
		media.Items[i].Delete()
	}
	return nil
}

func (media *FeedMedia) instagram() *Instagram {
	return media.insta
}

// SetInstagram set instagram
func (media *FeedMedia) SetInstagram(insta *Instagram) {
	media.insta = insta
}

// SetID sets media ID
// this value can be int64 or string
func (media *FeedMedia) SetID(id interface{}) {
	media.NextID = id
}

// Sync updates media values.
func (media *FeedMedia) Sync() error {
	id := media.ID()
	insta := media.insta

	data, err := insta.prepareData(
		map[string]interface{}{
			"media_id": id,
		},
	)
	if err != nil {
		return err
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlMediaInfo, id),
			Query:    generateSignature(data),
			IsPost:   false,
		},
	)
	if err != nil {
		return err
	}

	m := FeedMedia{}
	err = json.Unmarshal(body, &m)
	*media = m
	media.endpoint = urlMediaInfo
	media.insta = insta
	media.NextID = id
	media.setValues()
	return err
}

func (media *FeedMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

func (media FeedMedia) Error() error {
	return media.err
}

// ID returns media id.
func (media *FeedMedia) ID() string {
	switch s := media.NextID.(type) {
	case string:
		return s
	case int64:
		return strconv.FormatInt(s, 10)
	case json.Number:
		return string(s)
	}
	return ""
}

// Next allows pagination after calling:
// User.Feed
// extra query arguments can be passes one after another as func(key, value).
// Only if an even number of string arguements will be passed, they will be
//   used in the query.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problems have occurred.
func (media *FeedMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.insta
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	query := map[string]string{
		"exclude_comment":                 "true",
		"only_fetch_first_carousel_media": "false",
	}
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i = i + 2 {
			query[params[i].(string)] = params[i+1].(string)
		}
	}

	if next := media.ID(); next != "" {
		query["max_id"] = next
	}

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    query,
		},
	)
	if err == nil {
		m := FeedMedia{}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&m)
		if err == nil {
			media.NextID = m.NextID
			media.MoreAvailable = m.MoreAvailable
			media.NumResults = m.NumResults
			media.AutoLoadMoreEnabled = m.AutoLoadMoreEnabled
			media.Status = m.Status
			if m.NextID == 0 || !m.MoreAvailable {
				media.err = ErrNoMore
			}
			m.setValues()
			media.Items = append(media.Items, m.Items...)
			return true
		}
	}
	return false
}

// Latest returns a slice of the latest fetched items of the list of all items.
// The Next method keeps adding to the list, with Latest you can retrieve only
// the newest items.
func (media *FeedMedia) Latest() []Item {
	return media.Items[len(media.Items)-media.NumResults:]
}

// MediaItem defines a item media for the
// SavedMedia struct
type MediaItem struct {
	Media Item `json:"media"`
}

// SavedMedia stores the information about media being saved before in my account.
type SavedMedia struct {
	insta    *Instagram
	endpoint string

	err error

	Items []MediaItem `json:"items"`

	NumResults          int    `json:"num_results"`
	MoreAvailable       bool   `json:"more_available"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Status              string `json:"status"`

	NextID interface{} `json:"next_max_id"`
}

// Next allows pagination
func (media *SavedMedia) Next(params ...interface{}) bool {
	// Inital error check
	// if last pagination had errors
	if media.err != nil {
		return false
	}

	insta := media.insta
	endpoint := media.endpoint
	next := media.ID()

	opts := &reqOptions{
		Endpoint: endpoint,
		Query: map[string]string{
			"max_id": next,
		},
	}

	body, _, err := insta.sendRequest(opts)
	if err != nil {
		media.err = err
		return false
	}

	m := SavedMedia{}

	if err := json.Unmarshal(body, &m); err != nil {
		media.err = err
		return false
	}

	*media = m

	media.insta = insta
	media.endpoint = endpoint
	media.err = nil

	if m.NextID == 0 || !m.MoreAvailable {
		media.err = ErrNoMore
	}

	media.setValues()

	return true
}

// Error returns the SavedMedia error
func (media *SavedMedia) Error() error {
	return media.err
}

// ID returns the SavedMedia next id
func (media *SavedMedia) ID() string {
	switch id := media.NextID.(type) {
	case int64:
		return strconv.FormatInt(id, 10)
	case string:
		return id
	}
	return ""
}

// Delete method TODO
//
// I think this method should use the
// Unsave method, instead of the Delete.
func (media *SavedMedia) Delete() error {
	return nil
}

// instagram returns the media instagram
func (media *SavedMedia) instagram() *Instagram {
	return media.insta
}

// setValues set the SavedMedia items values
func (media *SavedMedia) setValues() {
	for i := range media.Items {
		setToMediaItem(&media.Items[i], media)
	}
}

// UploadPhoto post image from io.Reader to instagram.
func (insta *Instagram) UploadPhoto(photo io.Reader, photoCaption string, quality int, filterType int) (Item, error) {
	out := Item{}

	config, err := insta.postPhoto(photo, photoCaption, quality, filterType, false)
	if err != nil {
		return out, err
	}
	data, err := insta.prepareData(config)
	if err != nil {
		return out, err
	}

	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: "media/configure/?",
		Query:    generateSignature(data),
		IsPost:   true,
	})
	if err != nil {
		return out, err
	}
	var uploadResult struct {
		Media    Item   `json:"media"`
		UploadID string `json:"upload_id"`
		Status   string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}

// UploadVideo post video and thumbnail from io.Reader to instagram.
func (insta *Instagram) UploadVideo(video io.Reader, title string, caption string, thumbnail io.Reader) (Item, error) {
	out := Item{}
	config, err := insta.postVideo(video, title, caption, thumbnail)
	if err != nil {
		return out, err
	}

	data, err := insta.prepareData(config)
	if err != nil {
		return out, err
	}

	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: "media/configure/?",
		Query:    generateSignature(data),
		IsPost:   true,
	})
	if err != nil {
		return out, err
	}

	var result struct {
		Media    Item   `json:"media"`
		UploadID string `json:"upload_id"`
		Status   string `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return out, err
	}
	if result.Status != "ok" {
		return out, fmt.Errorf("unknown error, status: %s", result.Status)
	}

	return result.Media, nil
}

func (insta *Instagram) postThumbnail(uploadID int64, name string, thumbnail io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(thumbnail)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", baseUrl+"rupload_igphoto/"+name, buf)
	if err != nil {
		return err
	}
	req.Header.Set("X-IG-Capabilities", "3Q4=")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-type", "image/jpeg")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", instaUserAgent)
	req.Header.Set("X-Entity-Name", name)
	req.Header.Set("X-Entity-Length", strconv.FormatInt(req.ContentLength, 10))
	req.Header.Set("Offset", "0")
	ruploadParams := map[string]string{
		"media_type":          "2",
		"upload_id":           strconv.FormatInt(uploadID, 10),
		"upload_media_height": "240",
		"upload_media_width":  "320",
	}
	params, err := json.Marshal(ruploadParams)
	if err != nil {
		return err
	}
	req.Header.Set("X-Instagram-Rupload-Params", string(params))

	resp, err := insta.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("postThumbnail invalid status code, result: %s with body %s", resp.Status, string(body))
	}
	var result struct {
		UploadID       string      `json:"upload_id"`
		XsharingNonces interface{} `json:"xsharing_nonces"`
		Status         string      `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.Status != "ok" {
		return fmt.Errorf("unknown error, status: %s", result.Status)
	}

	return nil
}

func (insta *Instagram) postVideo(
	video io.Reader,
	title string,
	caption string,
	thumbnail io.Reader,
) (map[string]interface{}, error) {
	uploadID := time.Now().Unix()
	rndNumber := rand.Intn(9999999999-1000000000) + 1000000000
	name := "igtv_" + strconv.Itoa(rndNumber)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(video)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", baseUrl+"rupload_igvideo/"+name, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-IG-Capabilities", "3Q4=")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-type", "video/mp4")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", instaUserAgent)
	req.Header.Set("X-Entity-Name", name)
	req.Header.Set("X-Entity-Length", strconv.FormatInt(req.ContentLength, 10))
	req.Header.Set("Offset", "0")
	ruploadParams := map[string]string{
		"media_type":   "2",
		"video_format": "video/mp4",
		"upload_id":    strconv.FormatInt(uploadID, 10),
	}
	params, err := json.Marshal(ruploadParams)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Instagram-Rupload-Params", string(params))

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("postVideo invalid status code, result: %s with body %s", resp.Status, string(body))
	}
	var result struct {
		UploadID       string      `json:"upload_id"`
		XsharingNonces interface{} `json:"xsharing_nonces"`
		Status         string      `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("unknown error, status: %s", result.Status)
	}

	err = insta.postThumbnail(uploadID, name, thumbnail)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	config := map[string]interface{}{
		"caption":            caption,
		"upload_id":          strconv.FormatInt(uploadID, 10),
		"device_id":          insta.dID,
		"source_type":        4,
		"date_time_original": now.Format("2020:51:21 22:51:37"),
	}

	return config, nil
}

func (insta *Instagram) postPhoto(
	photo io.Reader,
	photoCaption string,
	quality int,
	filterType int,
	isSidecar bool,
) (map[string]interface{}, error) {
	uploadID := time.Now().Unix()
	rndNumber := rand.Intn(9999999999-1000000000) + 1000000000
	name := strconv.FormatInt(uploadID, 10) + "_0_" + strconv.Itoa(rndNumber)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(photo)
	if err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	req, err := http.NewRequest("POST", baseUrl+"rupload_igphoto/"+name, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-IG-Capabilities", "3Q4=")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-type", "application/octet-stream")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", instaUserAgent)
	req.Header.Set("X-Entity-Name", name)
	ruploadParams := map[string]string{
		"retry_context":     `{"num_step_auto_retry": 0, "num_reupload": 0, "num_step_manual_retry": 0}`,
		"media_type":        "1",
		"upload_id":         strconv.FormatInt(uploadID, 10),
		"xsharing_user_ids": "[]",
		"image_compression": `{"lib_name": "moz", "lib_version": "3.1.m", "quality": "80"}`,
	}
	params, err := json.Marshal(ruploadParams)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Instagram-Rupload-Params", string(params))
	req.Header.Set("Offset", "0")
	req.Header.Set("X-Entity-Length", strconv.FormatInt(req.ContentLength, 10))

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code, result: %s with body %s", resp.Status, string(body))
	}
	var result struct {
		UploadID       string      `json:"upload_id"`
		XsharingNonces interface{} `json:"xsharing_nonces"`
		Status         string      `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("unknown error, status: %s", result.Status)
	}
	width, height, err := getImageDimensionFromReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	now := time.Now()

	config := map[string]interface{}{
		"media_folder": "Instagram",
		"source_type":  4,
		"caption":      photoCaption,
		"upload_id":    strconv.FormatInt(uploadID, 10),
		"device_id":    insta.dID,
		"device":       goInstaDeviceSettings,
		"edits": map[string]interface{}{
			"crop_original_size": []int{width * 1.0, height * 1.0},
			"crop_center":        []float32{0.0, 0.0},
			"crop_zoom":          1.0,
			"filter_type":        filterType,
		},
		"extra": map[string]interface{}{
			"source_width":  width,
			"source_height": height,
		},
		"height":                height,
		"width":                 width,
		"camera_model":          goInstaDeviceSettings["model"],
		"scene_capture_type":    "standard",
		"timezone_offset":       "3600",
		"date_time_original":    now.Format("2020:51:21 22:51:37"),
		"date_time_digitalized": now.Format("2020:51:21 22:51:37"),
		"software":              "1",
	}
	return config, nil
}

// UploadAlbum post image from io.Reader to instagram.
func (insta *Instagram) UploadAlbum(
	photos []io.Reader,
	photoCaption string,
	quality int,
	filterType int,
) (Item, error) {
	out := Item{}

	var childrenMetadata []map[string]interface{}
	for _, photo := range photos {
		config, err := insta.postPhoto(photo, photoCaption, quality, filterType, true)
		if err != nil {
			return out, err
		}

		childrenMetadata = append(childrenMetadata, config)
	}
	albumUploadID := time.Now().Unix()

	config := map[string]interface{}{
		"caption":           photoCaption,
		"client_sidecar_id": albumUploadID,
		"children_metadata": childrenMetadata,
	}
	data, err := insta.prepareData(config)
	if err != nil {
		return out, err
	}

	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: "media/configure_sidecar/?",
		Query:    generateSignature(data),
		IsPost:   true,
	})
	if err != nil {
		return out, err
	}

	var uploadResult struct {
		Media           Item   `json:"media"`
		ClientSideCarID int64  `json:"client_sidecar_id"`
		Status          string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}
