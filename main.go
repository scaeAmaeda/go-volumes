package main

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pgaskin/kepubify/v4/kepub"
)

func main() {
	fgr := ""
	for i := 1; i < len(os.Args); i++ {
		fgr = fgr + os.Args[i] + "_"
	}

	os.Mkdir(".\\workingdir", os.ModePerm)
	temp := ".\\workingdir"
	defer os.RemoveAll(temp)
	var nameBook string = temp + "\\" + fgr + ".pdf"
	var nameBookEpub string = temp + "\\" + fgr + ".epub"
	var nameBookKepub string = ".\\" + fgr + ".kepub.epub"
	conf := model.NewDefaultConfiguration()
	chapters, err := os.ReadDir(".\\chapters")

	if err != nil {
		log.Fatal(err)
	}
	for _, f := range chapters {
		if f.IsDir() {
			var path string = ".\\chapters\\" + f.Name()
			fi, _ := os.ReadDir(path)
			allChap := []string{}
			for _, g := range fi {
				localFile := ".\\chapters\\" + f.Name() + "\\" + g.Name()
				allChap = append(allChap, localFile)
			}
			outputPDF := temp + "\\" + f.Name() + ".pdf"
			api.ImportImagesFile(allChap, outputPDF, nil, conf)
		}
	}
	// tout est dans le temp il faut concat
	allChapPDF := []string{}
	tempDir, _ := os.ReadDir(temp)
	for _, c := range tempDir {
		locFile := temp + "\\" + c.Name()
		allChapPDF = append(allChapPDF, locFile)
	}
	api.MergeCreateFile(allChapPDF, nameBook, false, nil)
	// ici virer le sommaire qui se fout parce qu'on est pas en CLI
	ctxPDF, _ := api.ReadContextFile(nameBook)
	totalPages := ctxPDF.PageCount + 1
	pagesToRemove := fmt.Sprintf("%d", totalPages)
	api.RemovePagesFile(nameBook, nameBook, []string{pagesToRemove}, conf)

	convertiToEpub(nameBook, nameBookEpub)

	ctx := context.Background()
	// 0 securité, fuck la sécurité on est pas en entreprise
	epubFile, _ := os.Open(nameBookEpub)
	defer epubFile.Close()
	stat, _ := epubFile.Stat()
	size := stat.Size()
	zipReader, _ := zip.NewReader(epubFile, size)
	outFile, _ := os.Create(nameBookKepub)
	defer outFile.Close()
	conv := kepub.NewConverter()
	err = conv.Convert(ctx, outFile, zipReader)
	if err != nil {
		log.Fatal(err)
	}
	os.RemoveAll(".\\chapters")
	os.Mkdir(".\\chapters", os.ModePerm)

}

func convertiToEpub(input string, output string) error {
	cmd := exec.Command("ebook-convert", input, output)
	err := cmd.Run()
	return err
}
