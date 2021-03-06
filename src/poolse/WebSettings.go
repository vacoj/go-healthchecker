package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func settingsWeb(rw http.ResponseWriter, req *http.Request) {
	SettingsMu.Lock()
	defer SettingsMu.Unlock()

	blob, err := json.Marshal(&SETTINGS)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}

func settingsReloadWeb(rw http.ResponseWriter, req *http.Request) {
	StatusMu.Lock()
	SettingsMu.Lock()
	defer SettingsMu.Unlock()
	defer StatusMu.Unlock()

	ts := STATUS.Targets

	longest := findLongestPollingInterval(ts)

	checkTime := time.Now().Unix()-SETTINGS.LastReload.Unix() < int64(longest+5)
	if checkTime {
		io.WriteString(rw, fmt.Sprintf(`
<html>
	<head>
		<meta http-equiv="refresh" content="%d;URL=/status">
	</head>
	<body>
		Settings are still being reloaded.
		<br/>
		<br/>
		New settings will be applied once the longest-running application monitor 
		checks in.  This could take up to %d seconds, and the page will refresh 
		automatically.
	</body>
</html>`,
			longest+6, longest+5))
	} else {
		SETTINGS.LastReload = time.Now()
		GlobalWaitGroupHelper(true)
		go SETTINGS.reloadSettings()
		// show caller new settings
		io.WriteString(rw, fmt.Sprintf(`
<html>
	<head>
		<meta http-equiv="refresh" content="%d;URL=/status">
	</head>
	<body>
		Settings are being reloaded. This could take up to %d seconds and the page will refresh 
		automatically.
	</body>
</html>`,
			longest+6, longest+5))
	}
}
