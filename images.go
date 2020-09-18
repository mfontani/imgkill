package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
)

// "docker images ..." returns JSON that unmarshals to this:
type imageJSONLine struct {
	Tag        string
	Repository string
	ID         string
	Size       string
	CreatedAt  string
}

func grabImages() []imageJSONLine {
	cmd := exec.Command(cmdType, "images", "--format", `{"Tag":"{{.Tag}}","Repository":"{{.Repository}}","ID":"{{.ID}}","Size":"{{.Size}}","CreatedAt":"{{.CreatedAt}}"}`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("Executing '%s': %v", cmd, err))
	}
	images := make([]imageJSONLine, 0)
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if len(line) > 0 {
			var v imageJSONLine
			if err := json.Unmarshal(line, &v); err != nil {
				panic(fmt.Sprintf("Unmarshalling '%v': %v", line, err))
			}
			if v.Repository != "<none>" && v.Tag != "<none>" {
				images = append(images, v)
			}
		}
	}
	sort.SliceStable(images, func(i, j int) bool {
		if images[i].Repository == images[j].Repository {
			return images[i].Tag > images[j].Tag
		}
		return images[i].Repository < images[j].Repository
	})
	return images
}
