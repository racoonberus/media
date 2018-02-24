package preset

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetVideoCodec(filename string) (string, error) {
	cc, err := Exec(fmt.Sprintf(
		"ffprobe -v error -select_streams v:0 -show_entries stream=codec_name -of default=nokey=1:noprint_wrappers=1 /bucket/input/%s | tail -1",
		filename,
	))
	if err != nil {
		log.Fatal(err)
	}
	codec := string(cc)
	codec = strings.Trim(codec, "\n")
	return codec, err
}

func Exec(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).Output()
}

func SmtInSlice(needle interface{}, haystack []interface{}) bool {
	for _, b := range haystack {
		if b == needle {
			return true
		}
	}
	return false
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type Preset interface {
	GetName() string
	Execute(filePaths []string) ([]string, error)
}

type WebVideo struct{}

func (wv *WebVideo) GetName() string {
	return "Web Video"
}

func (wv *WebVideo) Execute(filePaths []string) ([]string, error) {
	var out []string = []string{}

	for _, filename := range filePaths {

		if ok, _ := FileExists(filename); ok {
			continue
		}

		parts := strings.Split(filename, ".")
		name, extension := parts[0], parts[1]
		codec, err := GetVideoCodec(filename)
		outFile := fmt.Sprintf("/bucket/output/%s.mp4", name)

		var convertCmd string

		if extension == "avi" {
			convertCmd = fmt.Sprintf("ffmpeg -y -i /bucket/input/%s -vcodec libx264 -acodec libfaac %s",
				filename, outFile)
		}

		if codec == "h264" {
			convertCmd = fmt.Sprintf("ffmpeg -y -i /bucket/input/%s %s",
				filename, outFile)
		}

		if SmtInSlice(extension, []interface{}{"3gp", "3g2"}) && codec == "mpeg4" {
			convertCmd = fmt.Sprintf("ffmpeg -i /bucket/input/%s -vcodec copy -acodec copy %s",
				filename, outFile)
		}

		if extension == "3gp" && codec == "h263" {
			convertCmd = fmt.Sprintf("ffmpeg -i /bucket/input/%s -c:v libx264 -c:a aac -strict experimental %s",
				filename, outFile)
		}

		//log.Println(convertCmd)

		_, err = Exec(convertCmd)
		if err != nil {
			return nil, err
		}

		//time.Sleep(3 * time.Second)
		//Exec("ls -lah")
		//
		//if ok, _ := FileExists(outFile); ok {
		//	return nil, errors.New("Not found result file " + outFile + ". Maybe you have not decision for case CODEC=" + codec + ";EXTENSION=" + extension)
		//}
		//
		//_, err = Exec(fmt.Sprintf("rm -f /bucket/input/%s", filename))
		//if err != nil {
		//	return nil, err
		//}

		out = append(out, outFile)
	}

	return out, nil
}