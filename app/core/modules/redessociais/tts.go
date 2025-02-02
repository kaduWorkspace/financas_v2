package redessociais

import (
	"goravel/app/core"
	"os/exec"
)

const CURL_TTS_EXP = "curl -X POST 'https://ttsmp3.com/makemp3_new.php' -H 'Content-Type: application/x-www-form-urlencoded' -d 'msg=Olá, como posso ajudar?&lang=Joanna'"

func DownloadTTS(texto string) (error, string) {
    res, err := core.HttpRequest(
        "https://ttsmp3.com/makemp3_new.php",
        "POST",
        map[string]string{
            "Content-Type": "application/x-www-form-urlencoded",
        },
        "msg=Eu te amo meu amor?&lang=Joanna",
    )
    if err != nil {
        return err, ""
    }
    mapa := core.JsonToMap(res)
    if mapa["URL"] == nil {
        return nil, ""
    }
    exec.Command("wget", mapa["URL"].(string), "-O", "/tmp/audio.mp3")
    return nil, mapa["URL"].(string)
}

