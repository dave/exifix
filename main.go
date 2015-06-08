package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var date1 = regexp.MustCompile(`([12][90]\d\d)-([01]\d)-([0123]\d)`)
var date2 = regexp.MustCompile(`([12][90]\d\d)([01]\d)([0123]\d)_(\d\d)(\d\d)(\d\d)`)

func matchNowDate1(str string) (date time.Time, ok bool) {
	if date1.MatchString(str) {
		matches := date1.FindStringSubmatch(str)
		year, _ := strconv.ParseInt(matches[1], 10, 64)
		month, _ := strconv.ParseInt(matches[2], 10, 64)
		day, _ := strconv.ParseInt(matches[3], 10, 64)
		date := time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.UTC)
		return date, true
	}
	return time.Time{}, false
}
func matchNowDate2(str string) (date time.Time, ok bool) {
	if date2.MatchString(str) {
		matches := date2.FindStringSubmatch(str)
		year, _ := strconv.ParseInt(matches[1], 10, 64)
		month, _ := strconv.ParseInt(matches[2], 10, 64)
		day, _ := strconv.ParseInt(matches[3], 10, 64)
		h, _ := strconv.ParseInt(matches[4], 10, 64)
		m, _ := strconv.ParseInt(matches[5], 10, 64)
		s, _ := strconv.ParseInt(matches[6], 10, 64)
		date := time.Date(int(year), time.Month(int(month)), int(day), int(h), int(m), int(s), 0, time.UTC)
		return date, true
	}
	return time.Time{}, false
}

func betweenTimes(time1 time.Time, time2 time.Time) time.Duration {
	if time1.After(time2) {
		return time1.Sub(time2)
	}
	return time2.Sub(time1)
}

func areDatesClose(time1 time.Time, time2 time.Time) bool {
	between := betweenTimes(time1, time2)
	return between.Hours() < 24*7
}

func getPathDate(path string) time.Time {

	parts := strings.Split(path, "/")

	date := time.Time{}

	for _, part := range parts {
		// Date1 is of the format 2001-02-03, which I've usually added by hand, and never has the time.
		date1out, ok := matchNowDate1(part)
		if ok {
			date = date1out
		}
	}

	for _, part := range parts {
		// Date2 is of the format 20010203_010203, which is usually on the filename, and always contains
		// the time. So, we should use this in priority over date1.
		date2out, ok := matchNowDate2(part)
		if ok {
			date = date2out
		}
	}

	return date

}
func getFileDate(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
func getExifDate(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}
	}

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, err := x.DateTime()
	if err != nil {
		return time.Time{}
	}
	return tm
}

func setFileDate(path string, date time.Time) error {
	fmt.Printf("Setting file date to %v... ", date)
	//return nil
	touchDate := fmt.Sprintf("%04d%02d%02d%02d%02d.%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
	cmd := exec.Command("touch", "-t", touchDate, path)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done.\n")
	}
	return nil
}

func setExifDate(path string, date time.Time) error {
	fmt.Printf("Setting exif date to %v... ", date)
	//return nil
	exifDate := fmt.Sprintf("-AllDates='%04d:%02d:%02d %02d:%02d:%02d'", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
	cmd := exec.Command("exiftool", exifDate, "-overwrite_original", path)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done.\n")
	}
	return nil
}

func updateDateIfWrong(path string) error {

	pathDate := getPathDate(path)

	zero := time.Time{}
	if pathDate == zero {
		// If we don't have a date in the path, then just give up.
		fmt.Printf("Can't find path date\n")
		return nil
	}

	fileDate := getFileDate(path)
	exifDate := getExifDate(path)

	actualDate := pathDate

	if exifDate != zero && areDatesClose(pathDate, exifDate) {
		actualDate = exifDate
	} else if fileDate != zero && areDatesClose(pathDate, fileDate) {
		actualDate = fileDate
	}

	if exifDate == zero || !areDatesClose(actualDate, exifDate) {
		err := setExifDate(path, actualDate)
		if err != nil {
			return err
		}
	}

	fileDate = getFileDate(path)

	if fileDate == zero || !areDatesClose(actualDate, fileDate) {
		err := setFileDate(path, actualDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func process(path string) {
	err := updateDateIfWrong(path)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func main() {

	filepath.Walk("/Users/dave/Dropbox (Personal)/Pictures", func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".DS_Store") {
			return nil
		}

		fmt.Printf("%s\n", path)
		process(path)
		return nil
	})

}
