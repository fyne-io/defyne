package guidefs

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// IconNames is an array with the list of names of all the Icons
	IconNames []string

	// IconReverse Contains the key value pair where the key is the address of the icon and the value is the Name
	IconReverse map[string]string

	// Icons Has the hashmap of Icons from the standard theme.
	// ToDo: Will have to look for a way to sync the list from `fyne_demo`
	Icons map[string]fyne.Resource
)

func initIcons() {
	Icons = map[string]fyne.Resource{
		"CancelIcon":        theme.CancelIcon(),
		"ConfirmIcon":       theme.ConfirmIcon(),
		"DeleteIcon":        theme.DeleteIcon(),
		"SearchIcon":        theme.SearchIcon(),
		"SearchReplaceIcon": theme.SearchReplaceIcon(),

		"CheckButtonIcon":        theme.CheckButtonIcon(),
		"CheckButtonCheckedIcon": theme.CheckButtonCheckedIcon(),
		"CheckButtonFillIcon":    theme.CheckButtonFillIcon(),
		"RadioButtonIcon":        theme.RadioButtonIcon(),
		"RadioButtonCheckedIcon": theme.RadioButtonCheckedIcon(),
		"RadioButtonFillIcon":    theme.RadioButtonFillIcon(),

		"ColorAchromaticIcon": theme.ColorAchromaticIcon(),
		"ColorChromaticIcon":  theme.ColorChromaticIcon(),
		"ColorPaletteIcon":    theme.ColorPaletteIcon(),

		"ContentAddIcon":    theme.ContentAddIcon(),
		"ContentRemoveIcon": theme.ContentRemoveIcon(),
		"ContentClearIcon":  theme.ContentClearIcon(),
		"ContentCutIcon":    theme.ContentCutIcon(),
		"ContentCopyIcon":   theme.ContentCopyIcon(),
		"ContentPasteIcon":  theme.ContentPasteIcon(),
		"ContentRedoIcon":   theme.ContentRedoIcon(),
		"ContentUndoIcon":   theme.ContentUndoIcon(),

		"InfoIcon":     theme.InfoIcon(),
		"ErrorIcon":    theme.ErrorIcon(),
		"QuestionIcon": theme.QuestionIcon(),
		"WarningIcon":  theme.WarningIcon(),

		"BrokenImageIcon": theme.BrokenImageIcon(),

		"DocumentIcon":       theme.DocumentIcon(),
		"DocumentCreateIcon": theme.DocumentCreateIcon(),
		"DocumentPrintIcon":  theme.DocumentPrintIcon(),
		"DocumentSaveIcon":   theme.DocumentSaveIcon(),

		"FileIcon":            theme.FileIcon(),
		"FileApplicationIcon": theme.FileApplicationIcon(),
		"FileAudioIcon":       theme.FileAudioIcon(),
		"FileImageIcon":       theme.FileImageIcon(),
		"FileTextIcon":        theme.FileTextIcon(),
		"FileVideoIcon":       theme.FileVideoIcon(),
		"FolderIcon":          theme.FolderIcon(),
		"FolderNewIcon":       theme.FolderNewIcon(),
		"FolderOpenIcon":      theme.FolderOpenIcon(),
		"ComputerIcon":        theme.ComputerIcon(),
		"HomeIcon":            theme.HomeIcon(),
		"HelpIcon":            theme.HelpIcon(),
		"HistoryIcon":         theme.HistoryIcon(),
		"SettingsIcon":        theme.SettingsIcon(),
		"StorageIcon":         theme.StorageIcon(),
		"DownloadIcon":        theme.DownloadIcon(),
		"UploadIcon":          theme.UploadIcon(),

		"ViewFullScreenIcon": theme.ViewFullScreenIcon(),
		"ViewRestoreIcon":    theme.ViewRestoreIcon(),
		"ViewRefreshIcon":    theme.ViewRefreshIcon(),
		"VisibilityIcon":     theme.VisibilityIcon(),
		"VisibilityOffIcon":  theme.VisibilityOffIcon(),
		"ZoomFitIcon":        theme.ZoomFitIcon(),
		"ZoomInIcon":         theme.ZoomInIcon(),
		"ZoomOutIcon":        theme.ZoomOutIcon(),

		"MoreHorizontalIcon": theme.MoreHorizontalIcon(),
		"MoreVerticalIcon":   theme.MoreVerticalIcon(),

		"MoveDownIcon": theme.MoveDownIcon(),
		"MoveUpIcon":   theme.MoveUpIcon(),

		"NavigateBackIcon": theme.NavigateBackIcon(),
		"NavigateNextIcon": theme.NavigateNextIcon(),

		"MenuIcon":         theme.MenuIcon(),
		"MenuExpandIcon":   theme.MenuExpandIcon(),
		"MenuDropDownIcon": theme.MenuDropDownIcon(),
		"MenuDropUpIcon":   theme.MenuDropUpIcon(),

		"MailAttachmentIcon": theme.MailAttachmentIcon(),
		"MailComposeIcon":    theme.MailComposeIcon(),
		"MailForwardIcon":    theme.MailForwardIcon(),
		"MailReplyIcon":      theme.MailReplyIcon(),
		"MailReplyAllIcon":   theme.MailReplyAllIcon(),
		"MailSendIcon":       theme.MailSendIcon(),

		"MediaFastForwardIcon":  theme.MediaFastForwardIcon(),
		"MediaFastRewindIcon":   theme.MediaFastRewindIcon(),
		"MediaPauseIcon":        theme.MediaPauseIcon(),
		"MediaPlayIcon":         theme.MediaPlayIcon(),
		"MediaStopIcon":         theme.MediaStopIcon(),
		"MediaRecordIcon":       theme.MediaRecordIcon(),
		"MediaReplayIcon":       theme.MediaReplayIcon(),
		"MediaSkipNextIcon":     theme.MediaSkipNextIcon(),
		"MediaSkipPreviousIcon": theme.MediaSkipPreviousIcon(),

		"VolumeDownIcon": theme.VolumeDownIcon(),
		"VolumeMuteIcon": theme.VolumeMuteIcon(),
		"VolumeUpIcon":   theme.VolumeUpIcon(),

		"AccountIcon": theme.AccountIcon(),
		"LoginIcon":   theme.LoginIcon(),
		"LogoutIcon":  theme.LogoutIcon(),

		"ListIcon": theme.ListIcon(),
		"GridIcon": theme.GridIcon(),
	}
	IconNames = extractIconNames()
	IconReverse = reverseIconMap()
}

// extractIconNames returns all the list of names of all the Icons from the hashmap `Icons`
func extractIconNames() []string {
	var iconNamesFromData = make([]string, len(Icons))
	i := 0
	for k := range Icons {
		iconNamesFromData[i] = k
		i++
	}

	sort.Strings(iconNamesFromData)
	return iconNamesFromData
}

// reverseIconMap returns all the list of Icons and their addresses
func reverseIconMap() map[string]string {
	var iconReverseFromData = make(map[string]string, len(Icons))
	for k, v := range Icons {
		iconReverseFromData[v.Name()] = k
	}

	return iconReverseFromData
}
