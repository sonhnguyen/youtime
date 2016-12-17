package youtime

import (
	"fmt"
	"io"
	"math"
	"time"
)

// Subtitle contains the subtitle data for a .srt file
type Subtitle struct {
	Number int           // A sequential number for this subtitle
	Start  time.Duration // The duration since the start of the file when the subtitle should be shown
	End    time.Duration // The duration since the start of the file when the subtitle should be hidden
	Text   string        // The contents of the subtitle
	Bounds Rectangle     // Optional bounding box of the subtitles. Use .IsEmpty() to determine if it's set
}

// Rectangle represents a on-screen axis-aligned bounding box
type Rectangle struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

// Width calculates and returns the width of a Rectangle
func (r *Rectangle) Width() int {
	return r.Right - r.Left
}

// Height calculates and returns the height of a Rectangle
func (r *Rectangle) Height() int {
	return r.Bottom - r.Top
}

// FromSizes creates a new rectangle with its top-left corner at (x, y)
// and with the specified with and height
func FromSizes(x, y, wid, hgt int) Rectangle {
	return Rectangle{x, y, x + wid, y + hgt}
}

// IsEmpty checks if the Rectangle has the nil value
func (r *Rectangle) IsEmpty() bool {
	return r.Left == r.Right && r.Top == r.Bottom
}

// Writes a duration formatted as hours:minues:seconds,milliseconds
func writeTime(w io.Writer, dur time.Duration) (nbytes int, err error) {
	hoursToPrint := int(math.Floor(dur.Hours()))
	minutesToPrint := int(math.Floor(dur.Minutes() - (time.Duration(hoursToPrint) * time.Hour).Minutes()))
	secondsToPrint := int(math.Floor(dur.Seconds() - (time.Duration(hoursToPrint)*time.Hour + time.Duration(minutesToPrint)*time.Minute).Seconds()))
	millisecondsToPrint := int(math.Floor(float64(dur/time.Millisecond - (time.Duration(hoursToPrint)*time.Hour+time.Duration(minutesToPrint)*time.Minute+time.Duration(secondsToPrint)*time.Second)/time.Millisecond)))

	nbytes, err = fmt.Fprintf(w, "%02d:%02d:%02d,%03d", hoursToPrint, minutesToPrint, secondsToPrint, millisecondsToPrint)
	return
}

// Writes a bounding rectangle X1:left X2:right Y1:top Y2:bottom
func writeRect(w io.Writer, r Rectangle) (nbytes int, err error) {
	nbytes, err = fmt.Fprintf(w, "X1:%d X2:%d Y1:%d Y2:%d", r.Left, r.Right, r.Top, r.Bottom)
	return
}

// WriteTo writes a Subtitle-object to the given writer in srt-format.
// No validation of the Subtitle object is performed
func (s *Subtitle) WriteTo(writer io.Writer) (nbytes int, err error) {
	var wlen int

	wlen, err = fmt.Fprintf(writer, "%v\n", s.Number)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = writeTime(writer, s.Start)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, " --> ")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = writeTime(writer, s.End)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	if !s.Bounds.IsEmpty() {
		wlen, err = io.WriteString(writer, " ")
		nbytes += wlen
		if err != nil {
			return nbytes, err
		}

		wlen, err = writeRect(writer, s.Bounds)
		nbytes += wlen
		if err != nil {
			return nbytes, err
		}
	}

	wlen, err = io.WriteString(writer, "\n")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, s.Text)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, "\n\n")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	return nbytes, nil
}

// Write file
func WriteSubs(w io.Writer, subs []Subtitle) error {
	for _, sub := range subs {
		_, err := sub.WriteTo(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSubtitleByID(id string, mongodb Mongodb) (Video, error) {

	result, err := GetVideoByIdMongo(id, mongodb)
	if err != nil {
		return Video{}, err
	}
	return result, nil
}
