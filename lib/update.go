package lib

import (
	"path/filepath"

	"github.com/coreservice-io/package_client"
)

func StartAutoUpdate(sync_remote_update_secs bool, ini_update bool, download_folder string, update_success_callback func(), logger func(string)) {

	pc, _ := package_client.NewPackageClient(AUTO_UPDATE_CONFIG_TOKEN, AUTO_UPDATE_CONFIG_PACKAGEID,
		AUTO_UPDATE_CONFIG_CURRENT_VERSION, sync_remote_update_secs, func(pc *package_client.PackageClient, m *package_client.Msg_resp_app_version, err error) bool {

			if err == nil {
				app_detail_s := &package_client.AppDetail_Standard{}
				decode_err := pc.DecodeAppDetail(m, app_detail_s)
				if decode_err == nil {
					download_err := package_client.DownloadFile(filepath.Join(download_folder, "temp"), app_detail_s.Download_url, app_detail_s.File_hash)
					if download_err == nil {
						unziperr := package_client.UnZipTo(filepath.Join(download_folder, "temp"), download_folder, true)
						if unziperr == nil {
							update_success_callback()
							return true
						}
					}
				}
			}

			return false

		}, func(logstr string) {
			logger(logstr)
		})

	pc.SetAutoUpdateInterval(AUTO_UPDATE_CONFIG_UPDATE_INTERVAL_SECS) //.Update().StartAutoUpdate()

	if ini_update {
		pc.Update()
	}

	pc.StartAutoUpdate()

}
