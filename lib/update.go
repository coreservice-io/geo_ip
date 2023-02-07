package lib

import (
	"path/filepath"

	"github.com/coreservice-io/package_client"
)

// func Update(update_key string, download_folder string, logger func(string)) error {

// 	pc, pc_err := package_client.NewPackageClient(update_key, AUTO_UPDATE_CONFIG_PACKAGEID,
// 		"0.0.0", false, func(pc *package_client.PackageClient, m *package_client.Msg_resp_app_version) error {

// 			app_detail_s := &package_client.AppDetail_Standard{}
// 			decode_err := pc.DecodeAppDetail(m, app_detail_s)
// 			if decode_err != nil {
// 				return decode_err
// 			}

// 			logger("start downloading geo files")

// 			download_err := package_client.DownloadFile(filepath.Join(download_folder, "temp"), app_detail_s.Download_url, app_detail_s.File_hash)
// 			if download_err != nil {
// 				return download_err
// 			}

// 			logger("start unzipping files geo files")
// 			unziperr := package_client.UnZipTo(filepath.Join(download_folder, "temp"), download_folder, true)
// 			if unziperr != nil {
// 				return unziperr
// 			}

// 			return nil

// 		}, func(logstr string) {
// 			logger(logstr)
// 		}, func(logstr string) {
// 			logger(logstr)
// 		})

// 	if pc_err != nil {
// 		return pc_err
// 	}

// 	update_err := pc.Update()
// 	if update_err != nil {
// 		return update_err
// 	}

// 	return nil
// }

func StartAutoUpdate(update_key string, current_version string, sync_remote_update_secs bool, download_folder string, update_success_callback func(), logger func(string), err_logger func(string)) (*package_client.PackageClient, error) {

	pc, pc_err := package_client.NewPackageClient(update_key, AUTO_UPDATE_CONFIG_PACKAGEID,
		current_version, sync_remote_update_secs, func(pc *package_client.PackageClient, m *package_client.Msg_resp_app_version) error {

			app_detail_s := &package_client.AppDetail_Standard{}
			decode_err := pc.DecodeAppDetail(m, app_detail_s)
			if decode_err != nil {
				return decode_err
			}

			logger("starting download geoip data")
			download_err := package_client.DownloadFile(filepath.Join(download_folder, "temp"), app_detail_s.Download_url, app_detail_s.File_hash)
			if download_err != nil {
				return download_err
			}

			logger("starting unzip geoip data")
			unziperr := package_client.UnZipTo(filepath.Join(download_folder, "temp"), download_folder, true)
			if unziperr != nil {
				return unziperr
			}

			logger("unzip geoip data success")
			update_success_callback()
			return nil

		}, func(logstr string) {
			logger(logstr)
		}, func(logstr string) {
			err_logger(logstr)
		})

	if pc_err != nil {
		return nil, pc_err
	}

	start_err := pc.SetAutoUpdateInterval(AUTO_UPDATE_CONFIG_UPDATE_INTERVAL_SECS).StartAutoUpdate()
	if start_err != nil {
		return nil, start_err
	}

	return pc, nil
}
