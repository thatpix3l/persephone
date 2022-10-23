// Utilties for working with the query response you receive after sending query actions
package query

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/thatpix3l/persephone/pkg/zeropad"
)

type Response struct {
	HasInternalBattery               bool              // `queryID:"1"`
	BatteryLevelBars                 uint              // `queryID:"2"`
	HasExternalBattery               bool              // `queryID:"3"`
	ExternalBatteryPercent           uint              // `queryID:"4"`
	IsOverHeating                    bool              // `queryID:"6"`
	IsBusy                           bool              // `queryID:"8"`
	IsQuickCaptureEnabled            bool              // `queryID:"9"`
	IsEncoding                       bool              // `queryID:"10"`
	IsLcdLockActive                  bool              // `queryID:"11"`
	VideoProgressCounter             uint              // `queryID:"13"`
	IsWirelessConnectionsEnabled     bool              // `queryID:"17"`
	PairingStatus                    uint              // `queryID:"19"`
	PairingType                      uint              // `queryID:"20"`
	TimeSinceSuccessfulPairing       time.Duration     // `queryID:"21"`
	WifiScanStatus                   uint              // `queryID:"22"`
	TimeSinceCompletedWifiScan       time.Duration     // `queryID:"23"`
	WifiProvisionStatus              uint              // `queryID:"24"`
	RemoteControlVersion             uint              // `queryID:"26"`
	IsRemoteControlConnected         bool              // `queryID:"27"`
	WirelessPairingStatus            uint              // `queryID:"28"`
	WlanApSsid                       string            // `queryID:"29"`
	CameraApSsid                     string            // `queryID:"30"`
	WirelessDeviceCount              uint              // `queryID:"31"`
	IsPreviewStreamEnabled           bool              // `queryID:"32"`
	StorageStatus                    int               // `queryID:"33"`
	PhotosBeforeFull                 uint              // `queryID:"34"`
	VideoTimeBeforeFull              time.Duration     // `queryID:"35"`
	GroupPhotosBeforeFull            uint              // `queryID:"36"`
	TotalGroupVideos                 uint              // `queryID:"37"`
	TotalPhotos                      uint              // `queryID:"38"`
	TotalVideos                      uint              // `queryID:"39"`
	UpdateStatus                     uint              // `queryID:"41"`
	IsCancellingUpdate               bool              // `queryID:"42"`
	IsLocateCameraActive             bool              // `queryID:"45"`
	MultishotCountdown               uint              // `queryID:"49"`
	RemainingSpace                   datasize.ByteSize // `queryID:"54"`
	IsPreviewStreamSupported         bool              // `queryID:"55"`
	WifiBarStrentgh                  uint              // `queryID:"56"`
	TagHilightsCount                 uint              // `queryID:"58"`
	TimeSinceBootTagHilight          time.Duration     // `queryID:"59"`
	StatusUpdateMinIntervalMS        uint              // `queryID:"60"`
	TimelapseTimeBeforeFull          time.Duration     // `queryID:"64"`
	ExposureMode                     uint              // `queryID:"65"`
	ExposureX                        uint              // `queryID:"66"`
	ExposureY                        uint              // `queryID:"67"`
	IsGpsLocked                      bool              // `queryID:"68"`
	IsWifiRadioEnabled               bool              // `queryID:"69"`
	InternalBatteryPercent           uint              // `queryID:"70"`
	MicAccessoryStatus               uint              // `queryID:"74"`
	DigitalZoomPercent               uint              // `queryID:"75"`
	WifiBandMode                     uint              // `queryID:"76"`
	IsDigitalZoomActive              bool              // `queryID:"77"`
	IsVideoSettingsMobileFriendly    bool              // `queryID:"78"`
	IsFirstTimeMode                  bool              // `queryID:"79"`
	IsWifi5GHzBandAvailable          bool              // `queryID:"81"`
	IsReadyForCommands               bool              // `queryID:"82"`
	IsBatteryGoodForUpdates          bool              // `queryID:"83"`
	IsTooCold                        bool              // `queryID:"85"`
	Orientation                      uint              // `queryID:"86"`
	IsZoomableWhileEncoding          bool              // `queryID:"88"`
	FlatMode                         uint              // `queryID:"89"`
	VideoPresetID                    uint              // `queryID:"93"`
	PhotoPresetID                    uint              // `queryID:"94"`
	TimelapsePresetID                uint              // `queryID:"95"`
	PresetGroupID                    uint              // `queryID:"96"`
	PresetID                         uint              // `queryID:"97"`
	PresetModified                   uint              // `queryID:"98"`
	LiveBurstsBeforeFull             uint              // `queryID:"99"`
	LiveBursts                       uint              // `queryID:"100"`
	IsCaptureDelayCountingDown       bool              // `queryID:"101"`
	MediaModeStatus                  uint              // `queryID:"102"`
	TimeWarpSpeed                    uint              // `queryID:"103"`
	IsLinuxCoreActive                bool              // `queryID:"104"`
	CameraLensType                   uint              // `queryID:"105"`
	IsVideoHindsightCaptureActive    bool              // `queryID:"106"`
	ScheduledCapturePresetID         uint              // `queryID:"107"`
	IsScheduledCaptureSet            bool              // `queryID:"108"`
	MediaModeStatusBitmasked         uint              // `queryID:"110"`
	HasStorageMinimumWriteSpeed      bool              // `queryID:"111"`
	StorageWriteSpeedErrorsSinceBoot uint              // `queryID:"112"`
	IsTurboTransferActive            bool              // `queryID:"113"`
	CameraControlStatus              uint              // `queryID:"114"`
	IsConnectedViaUSB                bool              // `queryID:"115"`
	UsbControlStaus                  uint              // `queryID:"116"`
	TotalStorageSpace                datasize.ByteSize // `queryID:"117"`
}

// Convert a 64-bit byte slice to an unsigned int
func bytesToUint(b []byte) uint {
	return uint(binary.BigEndian.Uint64(zeropad.BigEndian64(b)))
}

// Unmarshal from "data", a Big-Endian encoded byte array consisting of the [status_ID, count_of_values, val_1, val2, ...] extracted from a full GoPro Query Response, into the struct.
func UnmarshalPartial(data []byte, r *Response) (int, error) {

	if data == nil {
		return 0, errors.New("byte array is nil")
	}

	if len(data) < 3 {
		return 0, fmt.Errorf("byte array length %d is less than minimum of 3: %v", len(data), data)
	}

	// Status ID
	id := int(data[0])

	// Suggested count of values, according to byte array
	suggestedValLength := int(data[1])

	// Slice of actual values
	valBytes := data[2:]

	// Actual count of values
	actualValLength := len(valBytes)
	if suggestedValLength != actualValLength {
		return 0, fmt.Errorf("byte array suggests value count of %d, does not match actual count %d: %v", suggestedValLength, actualValLength, data)

	}

	// Value of payload, as an unsigned int
	valUint := bytesToUint(valBytes)
	var updateBoolErr error = nil

	updateBool := func(boolVar *bool) {

		if valUint == 0 {
			*boolVar = false

		} else if valUint == 1 {
			*boolVar = true

		} else {
			updateBoolErr = fmt.Errorf("number is not 0 or 1: \"%v\"", valUint)

		}

	}

	updateTime := func(durationVar *time.Duration, durationType time.Duration) {
		*durationVar = durationType * time.Duration(valUint)
	}

	updateSpace := func(spaceVar *datasize.ByteSize, spaceType datasize.ByteSize) {
		*spaceVar = spaceType * datasize.ByteSize(valUint)
	}

	switch id {

	case 1:
		updateBool(&r.HasInternalBattery)

	case 2:
		r.BatteryLevelBars = valUint

	case 3:
		updateBool(&r.HasExternalBattery)

	case 4:
		r.ExternalBatteryPercent = valUint

	case 6:
		updateBool(&r.IsOverHeating)

	case 8:
		updateBool(&r.IsBusy)

	case 9:
		updateBool(&r.IsQuickCaptureEnabled)

	case 10:
		updateBool(&r.IsEncoding)

	case 11:
		updateBool(&r.IsLcdLockActive)

	case 13:
		r.VideoProgressCounter = valUint

	case 17:
		updateBool(&r.IsWirelessConnectionsEnabled)

	case 19:
		r.PairingStatus = valUint

	case 20:
		r.PairingType = valUint

	case 21:
		updateTime(&r.TimeSinceSuccessfulPairing, time.Millisecond)

	case 22:
		r.WifiScanStatus = valUint

	case 23:
		updateTime(&r.TimeSinceCompletedWifiScan, time.Millisecond)

	case 24:
		r.WifiProvisionStatus = valUint

	case 26:
		r.RemoteControlVersion = valUint

	case 27:
		updateBool(&r.IsRemoteControlConnected)

	case 28:
		r.WirelessPairingStatus = valUint

	case 29:
		r.WlanApSsid = string(valBytes[:])

	case 30:
		r.CameraApSsid = string(valBytes[:])

	case 31:
		r.WirelessDeviceCount = valUint

	case 32:
		updateBool(&r.IsPreviewStreamEnabled)

	case 33:
		_ = r.StorageStatus // TODO: Add case for converting the value bytes into a signed int for the camera storage status

	case 34:
		r.PhotosBeforeFull = valUint

	case 35:
		updateTime(&r.VideoTimeBeforeFull, time.Minute)

	case 36:
		r.GroupPhotosBeforeFull = valUint

	case 37:
		r.TotalGroupVideos = valUint

	case 38:
		r.TotalPhotos = valUint

	case 39:
		r.TotalVideos = valUint

	case 41:
		r.UpdateStatus = valUint

	case 42:
		updateBool(&r.IsCancellingUpdate)

	case 45:
		updateBool(&r.IsLocateCameraActive)

	case 49:
		r.MultishotCountdown = valUint

	case 54:
		updateSpace(&r.RemainingSpace, datasize.KB)

	case 55:
		updateBool(&r.IsPreviewStreamSupported)

	case 56:
		r.WifiBarStrentgh = valUint

	case 58:
		r.TagHilightsCount = valUint

	case 59:
		updateTime(&r.TimeSinceBootTagHilight, time.Millisecond)

	case 60:
		r.StatusUpdateMinIntervalMS = valUint

	case 64:
		updateTime(&r.TimelapseTimeBeforeFull, time.Minute)

	case 65:
		r.ExposureMode = valUint

	case 66:
		r.ExposureX = valUint

	case 67:
		r.ExposureY = valUint

	case 68:
		updateBool(&r.IsGpsLocked)

	case 69:
		updateBool(&r.IsWifiRadioEnabled)

	case 70:
		r.InternalBatteryPercent = valUint

	case 74:
		r.MicAccessoryStatus = valUint

	case 75:
		r.DigitalZoomPercent = valUint

	case 76:
		r.WifiBandMode = valUint

	case 77:
		updateBool(&r.IsDigitalZoomActive)

	case 78:
		updateBool(&r.IsVideoSettingsMobileFriendly)

	case 79:
		updateBool(&r.IsFirstTimeMode)

	case 81:
		updateBool(&r.IsWifi5GHzBandAvailable)

	case 82:
		updateBool(&r.IsReadyForCommands)

	case 83:
		updateBool(&r.IsBatteryGoodForUpdates)

	case 85:
		updateBool(&r.IsTooCold)

	case 86:
		r.Orientation = valUint

	case 88:
		updateBool(&r.IsZoomableWhileEncoding)

	case 89:
		r.FlatMode = valUint

	case 93:
		r.VideoPresetID = valUint

	case 94:
		r.PhotoPresetID = valUint

	case 95:
		r.TimelapsePresetID = valUint

	case 96:
		r.PresetGroupID = valUint

	case 97:
		r.PresetID = valUint

	case 98:
		r.PresetModified = valUint

	case 99:
		r.LiveBurstsBeforeFull = valUint

	case 100:
		r.LiveBursts = valUint

	case 101:
		updateBool(&r.IsCaptureDelayCountingDown)

	case 102:
		r.MediaModeStatus = valUint

	case 103:
		r.TimeWarpSpeed = valUint

	case 104:
		updateBool(&r.IsLinuxCoreActive)

	case 105:
		r.CameraLensType = valUint

	case 106:
		updateBool(&r.IsVideoHindsightCaptureActive)

	case 107:
		r.ScheduledCapturePresetID = valUint

	case 108:
		updateBool(&r.IsScheduledCaptureSet)

	case 110:
		r.MediaModeStatusBitmasked = valUint

	case 111:
		updateBool(&r.HasStorageMinimumWriteSpeed)

	case 112:
		r.StorageWriteSpeedErrorsSinceBoot = valUint

	case 113:
		updateBool(&r.IsTurboTransferActive)

	case 114:
		r.CameraControlStatus = valUint

	case 115:
		updateBool(&r.IsConnectedViaUSB)

	case 116:
		r.CameraControlStatus = valUint

	case 117:
		updateSpace(&r.TotalStorageSpace, datasize.KB)

	default:
		updateBoolErr = fmt.Errorf("status ID %d does not exist: \"%v\"", id, data)

	}

	return suggestedValLength, updateBoolErr

}
