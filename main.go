package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/deploymenttheory/go-bindings-win32/bindings/runtime/win32"
	"github.com/deploymenttheory/go-bindings-winrt/bindings/winrt/foundation"
	"github.com/deploymenttheory/go-bindings-winrt/bindings/winrt/media/control"
	streams "github.com/deploymenttheory/go-bindings-winrt/bindings/winrt/storage/streams"
)

func toSeconds(ts foundation.TimeSpan) float64 {
	return float64(ts.Duration) / 10_000_000
}

func formatDuration(seconds float64) string {
	total := int(seconds)

	minutes := total / 60
	secs := total % 60

	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		`\`, "_",
		`/`, "_",
		`:`, "_",
		`*`, "_",
		`?`, "_",
		`"`, "_",
		`<`, "_",
		`>`, "_",
		`|`, "_",
	)

	return strings.TrimSpace(replacer.Replace(name))
}

func saveArtwork(title string, data []byte) (string, error) {
	title = sanitizeFilename("artworks/" + title)

	filename := title + ".jpg"
	if _, err := os.Stat(filename); err == nil {
		return filename, err
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.WriteFile("artworks/"+filename+".jpg", data, 0644); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	fmt.Println("Artwork Base64:", encoded[:50])

	return "", nil
}

func GetNowPlaying() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	statics, err := control.GlobalSystemMediaTransportControlsSessionManagerStatics()

	if err != nil {
		log.Fatal("Get statics:", err)
	}
	defer statics.Release()

	mediaSession, err := statics.RequestAsync()
	if err != nil {
		log.Fatal("Get Meida Session:", err)
	}
	defer mediaSession.Release()
	manager, err := mediaSession.Await()
	if err != nil {
		log.Fatal("Await Manager:", err)
	}
	defer manager.Release()
	currentSession, err := manager.GetCurrentSession()
	if err != nil {
		log.Fatal("Get Current Session:", err)
	}
	defer currentSession.Release()
	mediaOp, err := currentSession.TryGetMediaPropertiesAsync()
	if err != nil {
		log.Fatal("Get Media Properties:", err)
	}
	defer mediaOp.Release()
	media, err := mediaOp.Await()
	if err != nil {
		log.Fatal("Await media:", err)
	}
	defer media.Release()
	artist, _ := media.Artist()
	title, _ := media.Title()
	album, _ := media.AlbumTitle()
	thumbnail, err := media.Thumbnail()
	if err != nil {
		log.Fatal(err)
	}
	defer thumbnail.Release()

	streamOp, err := thumbnail.OpenReadAsync()
	if err != nil {
		log.Fatal(err)
	}
	defer streamOp.Release()

	stream, err := streamOp.Await()
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Release()

	randomStream, err := win32.QueryInterface[streams.IRandomAccessStream](
		&stream.IUnknown,
		&streams.IID_IRandomAccessStream,
	)
	if err != nil {
		log.Fatal("QueryInterface:", err)
	}
	defer randomStream.Release()

	size, err := randomStream.Size()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Thumbnail size:", size)

	inputStream, err := randomStream.GetInputStreamAt(0)
	if err != nil {
		log.Fatal(err)
	}
	defer inputStream.Release()

	buffer, err := streams.Create(uint32(size))
	if err != nil {
		log.Fatal(err)
	}
	defer buffer.Release()

	readOp, err := inputStream.ReadAsync(
		&buffer.IBuffer,
		uint32(size),
		streams.InputStreamOptionsNone,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer readOp.Release()

	readBuffer, err := readOp.Await()
	if err != nil {
		log.Fatal(err)
	}
	defer readBuffer.Release()
	if err != nil {
		log.Fatal(err)
	}
	defer readBuffer.Release()

	length, err := readBuffer.Length()
	if err != nil {
		log.Fatal(err)
	}

	readerStatics, err := streams.DataReaderStatics()
	if err != nil {
		log.Fatal(err)
	}
	defer readerStatics.Release()

	reader, err := readerStatics.FromBuffer(readBuffer)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Release()

	data := make([]byte, length)

	for i := uint32(0); i < length; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		data[i] = b
	}

	fmt.Println("Artwork bytes:", len(data))
	saveArtwork(title, data)
	// playbacktype, _ := media.PlaybackType()
	// println(artist, title, subtitle, album, thumbnail, playbacktype)
	playback, err := currentSession.GetPlaybackInfo()
	if err != nil {
		log.Fatal("GetPlaybackInfo:", err)
	}

	defer playback.Release()
	status, err := playback.PlaybackStatus()
	if err != nil {
		log.Fatal("PlaybackStatus:", err)
	}
	timeline, err := currentSession.GetTimelineProperties()
	if err != nil {
		log.Fatal("Timeline:", err)
	}
	position, err := timeline.Position()
	if err != nil {
		log.Fatal(err)
	}

	endTime, err := timeline.EndTime()
	if err != nil {
		log.Fatal(err)
	}

	posSeconds := toSeconds(position)
	endSeconds := toSeconds(endTime)

	fmt.Println("Position:", formatDuration(posSeconds))
	fmt.Println("Duration:", formatDuration(endSeconds))

	// fmt.Println(posMinutes, endMinites)
	fmt.Println("========== NOW PLAYING ==========")
	fmt.Println("Title:  ", title)
	fmt.Println("Artist: ", artist)
	fmt.Println("Album:  ", album)
	fmt.Println("Status: ", status.String())
	// fmt.Println("Thumnail:  ", stream.As)
	fmt.Println("=================================")
}

func main() {
	// for {
	GetNowPlaying()
	// time.Sleep(1)
	// }
}
