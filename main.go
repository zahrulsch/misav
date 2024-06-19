package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func loadTarget(target string) []string {
	if target == "" {
		target = "target.txt"
	}

	f, err := os.ReadFile(target)

	if err != nil {
		panic(err)
	}

	targets := []string{}
	lines := strings.Split(string(f), "\n")

	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\t", "")

		if line != "" {
			targets = append(targets, line)
		}
	}

	return targets
}

func main() {
	app := NewApp()
	tars := loadTarget(app.Config.TargetFile)

	log.Printf("%s", app.Config.JSONString())

	for _, target := range tars {
		id, err := app.GetInitialPage(target)

		if err != nil {
			panic(err)
		}

		uriMedia, res, err := app.GetPlaylist(id)
		if err != nil {
			panic(err)
		}

		medias, err := app.GetMedia(id, res, uriMedia)
		if err != nil {
			panic(err)
		}

		var wg sync.WaitGroup
		var maxCon = app.Config.MaxConcurrentMedia

		mediaC := make(chan *Media, 100)

		go func() {
			for _, med := range medias {
				mediaC <- med
			}

			close(mediaC)
		}()

		successMap := []map[string]string{}
		bar := progressbar.DefaultBytes(
			int64(len(medias)),
			"Downloading "+id,
		)

		var mtx sync.Mutex

		for i := 0; i < maxCon; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for med := range mediaC {
					dir := ".tmp/"
					name := dir + id + "." + strconv.Itoa(med.ID)
					tempName := name + ".temp"
					file, err := os.OpenFile(tempName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

					if err != nil {
						panic(err)
					}

					err = app.GetMediaData(file, med.URI)

					if err != nil {
						panic(err)
					}

					file.Close()

					mtx.Lock()

					bar.Add(1)
					successMap = append(successMap, map[string]string{
						"id":   strconv.Itoa(med.ID),
						"path": tempName,
					})

					mtx.Unlock()
				}
			}()
		}

		wg.Wait()

		sort.Slice(successMap, func(i, j int) bool {
			ix, _ := strconv.Atoi(successMap[i]["id"])
			jx, _ := strconv.Atoi(successMap[j]["id"])

			return ix < jx
		})

		rootDir := "video"

		if app.Config.OutDir != "" {
			rootDir = app.Config.OutDir
		}

		out, _ := filepath.Abs(filepath.Join(rootDir, id+".mp4"))

		file, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		file.Truncate(0)

		for _, c := range successMap {
			p, _ := filepath.Abs(c["path"])
			fileIn, err := os.ReadFile(p)

			if err != nil {
				panic(err)
			}

			if _, err := file.Write(fileIn); err != nil {
				panic(err)
			}

			if err := os.Remove(p); err != nil {
				panic(err)
			}
		}
	}

}
