package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/hegedustibor/htgo-tts/handlers"
)

/**
 * Required:
 * - mplayer
 *
 * Use:
 *
 * speech := htgotts.Speech{Folder: "audio", Language: "en", Handler: MPlayer}
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Proxy    string
	Handler  handlers.PlayerInterface
}

func GenerateName(text string) string {
	textAsFilename := []byte(text)
	hash := sha256.Sum256(textAsFilename)

	// Convert hash to string
	hashString := hex.EncodeToString(hash[:])

	return "irp-" + hashString + ".mp3"
}

// Creates a speech file with a given name
func (speech *Speech) CreateSpeechFile(text string) (string, error) {
	var err error

	generateName := GenerateName(text)

	f := speech.Folder + "/" + generateName
	if err = speech.createFolderIfNotExists(speech.Folder); err != nil {
		return "", err
	}

	if err = speech.downloadIfNotExists(f, text); err != nil {
		return "", err
	}

	return f, nil
}

// Plays an existent .mp3 file
func (speech *Speech) PlaySpeechFile(fileName string) error {
	if speech.Handler == nil {
		mplayer := handlers.MPlayer{}
		return mplayer.Play(fileName)
	}

	return speech.Handler.Play(fileName)
}

// Speak downloads speech and plays it using mplayer
func (speech *Speech) Speak(text string) error {

	var err error
	// generatedHashName := speech.generateHashName(text)

	fileName, err := speech.CreateSpeechFile(text)
	if err != nil {
		return err
	}

	return speech.PlaySpeechFile(fileName)
}

/**
 * Create the folder if does not exists.
 */
func (speech *Speech) createFolderIfNotExists(folder string) error {
	dir, err := os.Open(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0700)
	}

	dir.Close()
	return nil
}

/**
 * Download the voice file if does not exists.
 */
func (speech *Speech) downloadIfNotExists(fileName string, text string) error {
	f, err := os.Open(fileName)
	if err != nil {
		dlURL := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)

		response, err := speech.urlResponse(dlURL, f)

		if err != nil {
			return err
		}
		defer response.Body.Close()

		output, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(output, response.Body)
		return err
	}

	f.Close()
	return nil
}

func (speech *Speech) urlResponse(dlUrl string, f *os.File) (resp *http.Response, err error) {
	var (
		response *http.Response
	)

	if speech.Proxy != "" {
		var proxyURL *url.URL
		proxyURL, err = url.Parse(speech.Proxy)

		if err != nil {
			return response, err
		}

		httpCli := &http.Client{Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

		return httpCli.Get(dlUrl)
	}

	return http.Get(dlUrl)
}

func main() {

	handlerIndex := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}

	// downloadAudio := func(w http.ResponseWriter, r *http.Request) {

	// }

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/index", handlerIndex)
	http.HandleFunc("/handle-text", handleText)
	http.HandleFunc("/download/audio", downloadAudio)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)

	// buildSpeech := Speech{Folder: "audio", Language: "in"}

	// text := "saya sangat suka sekali"

	// fileName, _ := buildSpeech.CreateSpeechFile(text)

	// buildSpeech.PlaySpeechFile(fileName)

}

func handleText(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	speech := r.FormValue("text")

	fmt.Printf(speech)
}

func downloadAudio(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("path")

	f, err := os.Open(path)

	if f != nil {
		defer f.Close()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())

	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
