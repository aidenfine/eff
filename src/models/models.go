package models

type FileFeedback struct {
	Score               int
	Path                string
	DeadCodeResults     DeadCodeReport
	BarrelExportResults BarrelExportReport
}
type DeadCodeReport struct {
	Score     int
	Locations []string
}
type BarrelExportReport struct {
	Score     int
	Locations []string
}
