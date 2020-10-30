package factorio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ModPortalStruct struct {
	DownloadsCount int    `json:"downloads_count"`
	Name           string `json:"name"`
	Owner          string `json:"owner"`
	Releases       []struct {
		DownloadURL string `json:"download_url"`
		FileName    string `json:"file_name"`
		InfoJSON    struct {
			FactorioVersion Version `json:"factorio_version"`
		} `json:"info_json"`
		ReleasedAt    time.Time `json:"released_at"`
		Sha1          string    `json:"sha1"`
		Version       Version   `json:"version"`
		Compatibility bool      `json:"compatibility"`
	} `json:"releases"`
	Summary string `json:"summary"`
	Title   string `json:"title"`
}

// get all mods uploaded to the factorio modPortal
func ModPortalList() (interface{}, error, int) {
	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods?page_size=max", nil)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		err = ModPortalError(string(text))
		log.Println(err)
		return nil, err, resp.StatusCode
	}

	var jsonVal interface{}
	err = json.Unmarshal(text, &jsonVal)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return jsonVal, nil, resp.StatusCode
}

// get the details (mod-info, releases, etc.) from a specific mod from the modPortal
func ModPortalModDetails(modId string) (ModPortalStruct, error, int) {
	var mod ModPortalStruct

	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods/"+modId, nil)
	if err != nil {
		return mod, err, http.StatusInternalServerError
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mod, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return mod, err, http.StatusInternalServerError
	}

	err = json.Unmarshal(text, &mod)
	if err != nil {
		return mod, err, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		err = ModPortalError(string(text))
		log.Println(err)
		return ModPortalStruct{}, err, resp.StatusCode
	}

	server := GetFactorioServer()

	installedBaseVersion := Version{}
	_ = installedBaseVersion.UnmarshalText([]byte(server.BaseModVersion))
	requiredVersion := NilVersion

	for key, release := range mod.Releases {
		requiredVersion = release.InfoJSON.FactorioVersion
		areVersionIdentical := requiredVersion.Equals(installedBaseVersion)
		isException := installedBaseVersion.Equals(Version{1, 0, 0, 0}) && requiredVersion.Equals(Version{0, 18, 0, 0})
		release.Compatibility = areVersionIdentical || isException
		mod.Releases[key] = release
	}

	return mod, nil, resp.StatusCode
}

//Log the user into factorio, so mods can be downloaded
func FactorioLogin(username string, password string) (error, int) {
	var err error

	resp, err := http.PostForm("https://auth.factorio.com/api-login",
		url.Values{"require_game_ownership": {"true"}, "username": {username}, "password": {password}})

	if err != nil {
		return err, http.StatusInternalServerError
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		err = ModPortalError(bodyString)
		log.Println(err)
		return err, resp.StatusCode
	}

	var successResponse []string
	err = json.Unmarshal(bodyBytes, &successResponse)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	credentials := Credentials{
		Username: username,
		Userkey:  successResponse[0],
	}

	err = credentials.Save()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
