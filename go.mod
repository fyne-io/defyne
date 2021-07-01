module github.com/fyne-io/defyne

go 1.15

require (
	fyne.io/fyne/v2 v2.0.3-rc2.0.20210626134702-d695b2f5a69e
	fyne.io/x/fyne v0.0.0-20210407180700-b277e3225fbe
	github.com/fyne-io/terminal v0.0.0-20210410133030-a03d1d963afd
)

// TODO figure out why using the version inline above does not work!?
replace fyne.io/fyne/v2 v2.0.2 => fyne.io/fyne/v2 v2.0.2-0.20210409183941-f4c9800228a3
