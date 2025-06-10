package models

type FileFeedback struct {
	Score               int
	Path                string
	DeadCodeResults     DeadCodeReport
	BarrelExportResults BarrelExportReport
}
type DeadCodeReport struct {
	Score int
}
type BarrelExportReport struct {
	Score int
}
